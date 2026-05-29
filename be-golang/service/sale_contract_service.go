package service

import (
	"context"
	"fmt"
	"launchpad/config"
	"launchpad/model"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

type SaleContractService struct {
	productContractService *ProductContractService
}

// NewSaleContractServiceImpl 创建销售合约服务实例
func NewSaleContractService(productContractService *ProductContractService) *SaleContractService {
	return &SaleContractService{
		productContractService: productContractService,
	}
}

// 事件常量
const (
	SaleEventCreated          = "SaleCreated(address,uint256,uint256,uint256)"
	SaleEventStartTime        = "StartTimeSet(uint256)"
	SaleEventRegistrationTime = "RegistrationTimeSet(uint256,uint256)"
	SALE_DEPLOY_TOPIC         = "SaleDeployed(address)"
)

// QuerySaleInfo 查询销售信息
func (s *SaleContractService) QuerySaleInfo(saleAddress string) (*model.ProductContract, error) {
	// 连接到以太坊节点
	client, err := ethclient.Dial(config.AppConfig.Owner.NetworkUrl)

	if err != nil {
		return nil, err
	}
	defer client.Close()

	// 创建合约实例
	contractAddress := common.HexToAddress(saleAddress)

	contract, err := model.NewLaunchPadSale(contractAddress, client)

	if err != nil {
		return nil, err
	}

	// 调用合约方法
	opts := &bind.CallOpts{
		Context: context.Background(),
	}

	// 获取sale信息
	saleInfo, err := contract.Sale(opts)
	if err != nil {
		return nil, err
	}

	// 获取registration信息
	registrationInfo, err := contract.Registration(opts)
	if err != nil {
		return nil, err
	}

	// 解析数据
	productPO := &model.ProductContract{
		SaleContractAddress:    contractAddress.Hex(),
		TokenAddress:           saleInfo.Token.Hex(),
		TokenPriceInPT:         saleInfo.TokenPriceInETH.String(),
		TotalTokensSold:        saleInfo.AmountOfTokensToSell.String(),
		SaleEnd:                time.Unix(saleInfo.SaleEnd.Int64(), 0),
		UnlockTime:             time.Unix(saleInfo.TokensUnlockTime.Int64(), 0),
		SaleStart:              time.Unix(saleInfo.SaleStart.Int64(), 0),
		RegistrationTimeEnds:   time.Unix(registrationInfo.RegistrationTimeEnds.Int64(), 0),
		RegistrationTimeStarts: time.Unix(registrationInfo.RegistrationTimeStarts.Int64(), 0),
	}

	return productPO, nil
}

func (s *SaleContractService) StartSaleFactoryListen() error {
	ctx := context.Background()
	client, err := ethclient.Dial(config.AppConfig.Owner.NetworkUrl)
	if err != nil {
		return fmt.Errorf("连接以太坊节点失败: %v", err)
	}
	defer client.Close()

	contractAddr := common.HexToAddress(config.AppConfig.Sales.SalesFactoryAddress)
	saleDeploy := s.CalculateEventSignature(SALE_DEPLOY_TOPIC)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddr},
		Topics: [][]common.Hash{
			{saleDeploy},
		},
	}

	log.Infof("初始化服务开启销售工厂扫描,address=%s", config.AppConfig.Sales.SalesFactoryAddress)
	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		return fmt.Errorf("扫描销售工厂历史日志失败: %w", err)
	}
	if err := s.handleSaleFactoryLogs(logs); err != nil {
		return err
	}

	latestBlock, err := client.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("获取最新区块失败: %w", err)
	}

	liveLogs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(ctx, query, liveLogs)
	if err == nil {
		for {
			select {
			case err := <-sub.Err():
				return fmt.Errorf("销售工厂订阅中断: %w", err)
			case vLog := <-liveLogs:
				if err := s.handleSaleFactoryLogs([]types.Log{vLog}); err != nil {
					log.Errorf("处理销售工厂事件失败: %v", err)
				}
			}
		}
	}

	log.Warnf("销售工厂订阅不可用，切换到轮询模式: %v", err)
	return s.pollSaleFactoryLogs(ctx, client, query, latestBlock+1)
}

func (s *SaleContractService) ListenSaleChange(saleAddress string) error {
	ctx := context.Background()
	log.Infof("开始扫描销售合约变化,address=%s", saleAddress)
	client, err := ethclient.Dial(config.AppConfig.Owner.NetworkUrl)
	if err != nil {
		return fmt.Errorf("连接以太坊节点失败: %v", err)
	}
	defer client.Close()

	contractAddr := common.HexToAddress(saleAddress)
	topicCreated := s.CalculateEventSignature(SaleEventCreated)
	topicStartTime := s.CalculateEventSignature(SaleEventStartTime)
	topicRegistration := s.CalculateEventSignature(SaleEventRegistrationTime)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddr},
		Topics: [][]common.Hash{
			{topicCreated, topicStartTime, topicRegistration},
		},
	}

	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		return fmt.Errorf("扫描销售合约日志失败: %w", err)
	}

	return s.handleSaleChangeLogs(saleAddress, logs)
}

func (s *SaleContractService) handleSaleFactoryLogs(logs []types.Log) error {
	for _, vLog := range logs {
		saleAddress, err := bytesToAddressFrom20Bytes(vLog.Data)
		if err != nil {
			return err
		}
		log.Infof("监听到销售部署事件,factory=%s,saleAddress=%s", vLog.Address.Hex(), saleAddress.Hex())
		if err := s.ListenSaleChange(saleAddress.Hex()); err != nil {
			return err
		}
	}
	return nil
}

func (s *SaleContractService) handleSaleChangeLogs(saleAddress string, logs []types.Log) error {
	for _, vLog := range logs {
		log.Infof("监听到销售变更消息,address=%s,topic0=%s", saleAddress, vLog.Topics[0].Hex())
		sale, err := s.QuerySaleInfo(saleAddress)
		if err != nil {
			return err
		}
		if err := s.productContractService.UpsertSaleInfo(sale); err != nil {
			return err
		}
		log.Infof("销售合约数据已同步,address=%s", saleAddress)
	}

	return nil
}

func (s *SaleContractService) pollSaleFactoryLogs(ctx context.Context, client *ethclient.Client, baseQuery ethereum.FilterQuery, fromBlock uint64) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			latestBlock, err := client.BlockNumber(ctx)
			if err != nil {
				log.Errorf("获取最新区块失败: %v", err)
				continue
			}
			if latestBlock < fromBlock {
				continue
			}
			query := baseQuery
			query.FromBlock = new(big.Int).SetUint64(fromBlock)
			query.ToBlock = new(big.Int).SetUint64(latestBlock)
			logs, err := client.FilterLogs(ctx, query)
			if err != nil {
				log.Errorf("轮询销售工厂日志失败: %v", err)
				continue
			}
			if err := s.handleSaleFactoryLogs(logs); err != nil {
				log.Errorf("处理销售工厂日志失败: %v", err)
			}
			fromBlock = latestBlock + 1
		}
	}
}

func (s *SaleContractService) CalculateEventSignature(eventSignature string) common.Hash {
	hash := crypto.Keccak256Hash([]byte(eventSignature))

	return hash
}

func bytesToAddressFrom20Bytes(input []byte) (common.Address, error) {

	// 直接转换
	address := common.BytesToAddress(input)
	return address, nil
}
