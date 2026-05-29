package response

import (
	"encoding/json"
	"launchpad/model"
)

type ProductContractResponse struct {
	ID                        uint    `json:"id"`
	Name                      string  `json:"name"`
	Description               string  `json:"description"`
	Img                       string  `json:"img"`
	Status                    uint    `json:"status"`
	Amount                    string  `json:"amount"`
	SaleContractAddress       string  `json:"saleContractAddress"`
	TokenAddress              string  `json:"tokenAddress"`
	PaymentToken              string  `json:"paymentToken"`
	Follower                  uint    `json:"follower"`
	Tge                       int64   `json:"tge"`
	ProjectWebsite            string  `json:"projectWebsite"`
	AboutHTML                 string  `json:"aboutHtml"`
	RegistrationTimeStarts    int64   `json:"registrationTimeStarts"`
	RegistrationTimeEnds      int64   `json:"registrationTimeEnds"`
	SaleStart                 int64   `json:"saleStart"`
	SaleEnd                   int64   `json:"saleEnd"`
	MaxParticipation          string  `json:"maxParticipation"`
	TokenPriceInPT            string  `json:"tokenPriceInPT"`
	TotalTokensSold           string  `json:"totalTokensSold"`
	AmountOfTokensToSell      string  `json:"amountOfTokensToSell"`
	TotalRaised               string  `json:"totalRaised"`
	Symbol                    string  `json:"symbol"`
	Decimals                  uint    `json:"decimals"`
	UnlockTime                int64   `json:"unlockTime"`
	Medias                    string  `json:"medias"`
	NumberOfRegistrations     uint    `json:"numberOfRegistrations"`
	Vesting                   string  `json:"vesting"`
	Tricker                   string  `json:"tricker"`
	TokenName                 string  `json:"tokenName"`
	Roi                       string  `json:"roi"`
	VestingPortionsUnlockTime []int64 `json:"vestingPortionsUnlockTime"`
	VestingPercentPerPortion  []int64 `json:"vestingPercentPerPortion"`
	CreateTime                int64   `json:"createTime"`
	UpdateTime                int64   `json:"updateTime"`
	Type                      uint    `json:"type"`
	CardLink                  string  `json:"cardLink"`
	TweetID                   string  `json:"tweetId"`
	ChainID                   uint    `json:"chainId"`
	PaymentTokenDecimals      uint    `json:"paymentTokenDecimals"`
	CurrentPrice              uint    `json:"currentPrice"`
}

// ProductContractToProductContractResponse 转换 ProductContract 为 ProductContractResponse
func ProductContractToProductContractResponse(productContract model.ProductContract) ProductContractResponse {
	vestingPortionsUnlockTime := []int64{}
	if productContract.VestingPortionsUnlockTime != "" {
		_ = json.Unmarshal([]byte(productContract.VestingPortionsUnlockTime), &vestingPortionsUnlockTime)
	}
	vestingPercentPerPortion := []int64{}
	if productContract.VestingPercentPerPortion != "" {
		_ = json.Unmarshal([]byte(productContract.VestingPercentPerPortion), &vestingPercentPerPortion)
	}
	return ProductContractResponse{
		ID:                        productContract.ID,
		Name:                      productContract.Name,
		Description:               productContract.Description,
		Img:                       productContract.Img,
		Status:                    productContract.Status,
		Amount:                    productContract.Amount,
		SaleContractAddress:       productContract.SaleContractAddress,
		TokenAddress:              productContract.TokenAddress,
		PaymentToken:              productContract.PaymentToken,
		Follower:                  productContract.Follower,
		Tge:                       productContract.Tge.UnixMilli(),
		ProjectWebsite:            productContract.ProjectWebsite,
		AboutHTML:                 productContract.AboutHTML,
		RegistrationTimeStarts:    productContract.RegistrationTimeStarts.UnixMilli(),
		RegistrationTimeEnds:      productContract.RegistrationTimeEnds.UnixMilli(),
		SaleStart:                 productContract.SaleStart.UnixMilli(),
		SaleEnd:                   productContract.SaleEnd.UnixMilli(),
		MaxParticipation:          productContract.MaxParticipation,
		TokenPriceInPT:            productContract.TokenPriceInPT,
		TotalTokensSold:           productContract.TotalTokensSold,
		AmountOfTokensToSell:      productContract.AmountOfTokensToSell,
		TotalRaised:               productContract.TotalRaised,
		Symbol:                    productContract.Symbol,
		Decimals:                  productContract.Decimals,
		UnlockTime:                productContract.UnlockTime.UnixMilli(),
		Medias:                    productContract.Medias,
		NumberOfRegistrations:     productContract.NumberOfRegistrations,
		Vesting:                   productContract.Vesting,
		Tricker:                   productContract.Tricker,
		TokenName:                 productContract.TokenName,
		Roi:                       productContract.Roi,
		VestingPortionsUnlockTime: vestingPortionsUnlockTime,
		VestingPercentPerPortion:  vestingPercentPerPortion,
		CreateTime:                productContract.CreateTime.UnixMilli(),
		UpdateTime:                productContract.UpdateTime.UnixMilli(),
		Type:                      productContract.Type,
		CardLink:                  productContract.CardLink,
		TweetID:                   productContract.TweetID,
		ChainID:                   productContract.ChainID,
		PaymentTokenDecimals:      productContract.PaymentTokenDecimals,
		CurrentPrice:              productContract.CurrentPrice,
	}
}
