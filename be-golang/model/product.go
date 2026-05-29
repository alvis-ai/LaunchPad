package model

import "time"

// ProductPO 产品模型（使用指针类型处理可选字段）
type ProductPO struct {
	ID                uint    
	SaleAddress       string   
	SaleToken         string   
	SaleOwner         string   
	TokenPriceInPT    string    
	PaymentToken      string   
	TotalTokens       string    
	SaleEndTime       time.Time
	TokensUnlockTime  time.Time
	RegistrationStart time.Time
	RegistrationEnd   time.Time
	SaleStartTime     time.Time
}