package model

import (
	"time"
)

type ProductContract struct {
	ID                        uint      `json:"id" gorm:"primarykey"`
	Name                      string    `json:"name" gorm:"type:varchar(80);not null"`
	Description               string    `json:"description" gorm:"type:longtext"`
	Img                       string    `json:"img" gorm:"type:varchar(500)"`
	TwitterName               string    `json:"twitterName" gorm:"type:varchar(40)"`
	Status                    uint      `json:"status" gorm:"type:int(4);default:0;not null"`
	Amount                    string    `json:"amount" gorm:"type:varchar(40)"`
	SaleContractAddress       string    `json:"saleContractAddress" gorm:"type:varchar(42)"`
	TokenAddress              string    `json:"tokenAddress" gorm:"type:varchar(42)"`
	PaymentToken              string    `json:"paymentToken" gorm:"type:varchar(42)"`
	Follower                  uint      `json:"follower" gorm:"type:int(8);default:0;not null"`
	Tge                       time.Time `json:"tge" gorm:"type:datetime"`
	ProjectWebsite            string    `json:"projectWebsite" gorm:"type:varchar(500)"`
	AboutHTML                 string    `json:"aboutHtml" gorm:"type:varchar(500)"`
	RegistrationTimeStarts    time.Time `json:"registrationTimeStarts" gorm:"type:datetime"`
	RegistrationTimeEnds      time.Time `json:"registrationTimeEnds" gorm:"type:datetime"`
	SaleStart                 time.Time `json:"saleStart" gorm:"type:datetime"`
	SaleEnd                   time.Time `json:"saleEnd" gorm:"type:datetime"`
	MaxParticipation          string    `json:"maxParticipation" gorm:"type:varchar(40)"`
	TokenPriceInPT            string    `json:"tokenPriceInPT" gorm:"type:varchar(40);column:token_price_in_PT"`
	TotalTokensSold           string    `json:"totalTokensSold" gorm:"type:varchar(40)"`
	AmountOfTokensToSell      string    `json:"amountOfTokensToSell" gorm:"type:varchar(60)"`
	TotalRaised               string    `json:"totalRaised" gorm:"type:varchar(60)"`
	Symbol                    string    `json:"symbol" gorm:"type:varchar(60)"`
	Decimals                  uint      `json:"decimals" gorm:"type:int(8);default:8"`
	UnlockTime                time.Time `json:"unlockTime" gorm:"type:datetime"`
	Medias                    string    `json:"medias" gorm:"type:varchar(200)"`
	NumberOfRegistrations     uint      `json:"numberOfRegistrations" gorm:"type:int(8);column:number_of_registrants"`
	Vesting                   string    `json:"vesting" gorm:"type:varchar(40)"`
	Tricker                   string    `json:"tricker" gorm:"type:varchar(40)"`
	TokenName                 string    `json:"tokenName" gorm:"type:varchar(20)"`
	Roi                       string    `json:"roi" gorm:"type:varchar(20)"`
	VestingPortionsUnlockTime string    `json:"vestingPortionsUnlockTime" gorm:"type:varchar(60)"`
	VestingPercentPerPortion  string    `json:"vestingPercentPerPortion" gorm:"type:varchar(60)"`
	CreateTime                time.Time `json:"createTime" gorm:"type:datetime; not null"`
	UpdateTime                time.Time `json:"updateTime" gorm:"type:datetime; not null"`
	Type                      uint      `json:"type" gorm:"type:int(8)"`
	CardLink                  string    `json:"cardLink" gorm:"type:varchar(200)"`
	TweetID                   string    `json:"tweetId" gorm:"type:varchar(40)"`
	ChainID                   uint      `json:"chainId" gorm:"type:int(8)"`
	PaymentTokenDecimals      uint      `json:"paymentTokenDecimals" gorm:"type:int(8)"`
	CurrentPrice              uint      `json:"currentPrice" gorm:"type:bigint(12)"`
}

// TableName 指定表名
func (ProductContract) TableName() string {
	return "product_contract"
}
