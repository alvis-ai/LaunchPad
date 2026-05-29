package service

import (
	utils "launchpad/utils"
	"log"
)

type EncodeService struct {
}

func NewEncodeService() *EncodeService {
	return &EncodeService{}
}

func (s *EncodeService) Sign(hexString string) string {
	log.Printf("signing hex string: %s", hexString)
	return utils.GetSign(hexString)
}
