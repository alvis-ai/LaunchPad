package repository

import (
	"launchpad/api/request"
	"launchpad/model"

	"gorm.io/gorm"
)

type ProductContractRepository struct {
	DB *gorm.DB
}

func NewProductContractRepository(db *gorm.DB) *ProductContractRepository {
	return &ProductContractRepository{DB: db}
}

func (r *ProductContractRepository) GetById(productId string) (*model.ProductContract, error) {
	var productContract model.ProductContract

	if err := r.DB.Where("id = ?", productId).Find(&productContract).Error; err != nil {
		return nil, err
	}
	return &productContract, nil
}

func (r *ProductContractRepository) List() ([]model.ProductContract, error) {
	var productContracts []model.ProductContract

	result := r.DB.Model(&model.ProductContract{}).Find(&productContracts)
	if result.Error != nil {
		return nil, result.Error
	}

	return productContracts, nil
}

func (r *ProductContractRepository) Update(productContractUpdateRequest *request.ProductContractUpdateRequest) error {
	query := r.DB.Model(&model.ProductContract{})
	if productContractUpdateRequest.ID > 0 {
		query = query.Where("id = ?", productContractUpdateRequest.ID)
	} else {
		query = query.Where("sale_contract_address = ?", productContractUpdateRequest.SaleAddress)
	}

	return query.Select(
		"SaleContractAddress",
		"TokenAddress",
		"PaymentToken",
		"TokenPriceInPT",
		"TotalTokensSold",
		"SaleEnd",
		"UnlockTime",
		"RegistrationTimeStarts",
		"RegistrationTimeEnds",
		"SaleStart",
	).Updates(model.ProductContract{
		SaleContractAddress:    productContractUpdateRequest.SaleAddress,
		TokenAddress:           productContractUpdateRequest.SaleToken,
		TokenPriceInPT:         productContractUpdateRequest.TokenPriceInPT,
		TotalTokensSold:        productContractUpdateRequest.TotalTokens,
		SaleEnd:                productContractUpdateRequest.SaleEndTime.Time(),
		UnlockTime:             productContractUpdateRequest.TokensUnlockTime.Time(),
		RegistrationTimeStarts: productContractUpdateRequest.RegistrationStart.Time(),
		RegistrationTimeEnds:   productContractUpdateRequest.RegistrationEnd.Time(),
		SaleStart:              productContractUpdateRequest.SaleStartTime.Time(),
	}).Error
}

func (r *ProductContractRepository) UpsertSaleInfo(productContract *model.ProductContract) error {
	var existing model.ProductContract
	err := r.DB.Where("sale_contract_address = ?", productContract.SaleContractAddress).First(&existing).Error
	if err == nil {
		productContract.ID = existing.ID
		return r.DB.Model(&existing).Select(
			"SaleContractAddress",
			"TokenAddress",
			"TokenPriceInPT",
			"TotalTokensSold",
			"SaleEnd",
			"UnlockTime",
			"RegistrationTimeStarts",
			"RegistrationTimeEnds",
			"SaleStart",
		).Updates(productContract).Error
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	if productContract.Name == "" {
		productContract.Name = "IDO " + productContract.SaleContractAddress
	}
	return r.DB.Create(productContract).Error
}
