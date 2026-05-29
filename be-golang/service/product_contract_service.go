package service

import (
	"launchpad/api/request"
	"launchpad/model"
	"launchpad/repository"
)

type ProductContractService struct {
	ProductContractRepository *repository.ProductContractRepository // Ensure this type is defined in the repository package
}

func NewProductContractService(productContractRepository *repository.ProductContractRepository) *ProductContractService {
	return &ProductContractService{
		ProductContractRepository: productContractRepository,
	}
}

func (s *ProductContractService) GetById(productId string) (*model.ProductContract, error) {
	return s.ProductContractRepository.GetById(productId)
}

func (s *ProductContractService) List() ([]model.ProductContract, error) {
	return s.ProductContractRepository.List()
}

func (s *ProductContractService) Update(productContractUpdateRequest *request.ProductContractUpdateRequest) error {
	return s.ProductContractRepository.Update(productContractUpdateRequest)
}

func (s *ProductContractService) UpsertSaleInfo(productContract *model.ProductContract) error {
	return s.ProductContractRepository.UpsertSaleInfo(productContract)
}
