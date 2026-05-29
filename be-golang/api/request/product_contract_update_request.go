package request

import (
	"encoding/json"
	"strconv"
	"time"
)

// TimestampMs 自定义时间类型，用于处理毫秒时间戳字符串
type TimestampMs time.Time

func (t *TimestampMs) UnmarshalJSON(data []byte) error {
	var timestampStr string
	if err := json.Unmarshal(data, &timestampStr); err != nil {
		return err
	}
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return err
	}
	*t = TimestampMs(time.Unix(0, timestamp*int64(time.Millisecond)))
	return nil
}

func (t TimestampMs) Time() time.Time {
	return time.Time(t)
}

type ProductContractUpdateRequest struct {
	ID                uint        `json:"id" binding:"required"`
	SaleAddress       string      `json:"saleAddress"`
	SaleToken         string      `json:"saleToken"`
	SaleOwner         string      `json:"saleOwner"`
	TokenPriceInPT    string      `json:"tokenPriceInPT"`
	PaymentToken      string      `json:"paymentToken"`
	TotalTokens       string      `json:"totalTokens"`
	SaleEndTime       TimestampMs `json:"saleEndTime"`
	TokensUnlockTime  TimestampMs `json:"tokensUnlockTime"`
	RegistrationStart TimestampMs `json:"registrationStart"`
	RegistrationEnd   TimestampMs `json:"registrationEnd"`
	SaleStartTime     TimestampMs `json:"saleStartTime"`
}
