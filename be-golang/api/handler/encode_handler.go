package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"

	"launchpad/enums"
	"launchpad/middleware"
	"launchpad/model"
	"launchpad/service"
	"launchpad/utils"
)

const HEX_PREFIX = "0x"

type EncodeHandler struct {
	EncodeService *service.EncodeService
}

func NewEncodeHandler(encodeService *service.EncodeService) *EncodeHandler {
	return &EncodeHandler{
		EncodeService: service.NewEncodeService(),
	}
}

// @Summary Sign Registration
// @Description Sign Registration
// @Tags encode
// @Accept json
// @Produce json
// @Param userAddress formData string true "User Address"
// @Param contractAddress formData string true "Contract Address"
// @Success 200 {string} string "Sign"
// @Router /boba/encode/sign_registration [post]
func (h *EncodeHandler) SignRegistration(c *gin.Context) {
	userAddress := c.PostForm("userAddress")
	contractAddress := c.PostForm("contractAddress")
	log.Printf("userAddress: %s, contractAddress: %s", userAddress, contractAddress)
	if userAddress == "" || contractAddress == "" {
		panic("userAddress or contractAddress is empty")
	}
	contractAddressCleaned := utils.CleanHexPrefix(contractAddress)
	userAddressCleaned := utils.CleanHexPrefix(userAddress)
	concat := prependHexPrefix(strings.ToLower(userAddressCleaned + contractAddressCleaned))
	sign := h.EncodeService.Sign(concat)
	c.JSON(http.StatusOK, model.OkWithData[string](sign))
}

func containsHexPrefix(input string) bool {
	return len(input) > 1 && input[0] == '0' && input[1] == 'x'
}

func prependHexPrefix(input string) string {
	if !containsHexPrefix(input) {
		return HEX_PREFIX + input
	}
	return input
}

// @Summary Sign Participation
// @Description Sign Participation
// @Tags encode
// @Accept json
// @Produce json
// @Param userAddress formData string true "User Address"
// @Param contractAddress formData string true "Contract Address"
// @Param amount formData string true "Amount"
// @Success 200 {string} string "Sign"
// @Router /boba/encode/sign_participation [post]
func (h *EncodeHandler) SignParticipation(c *gin.Context) {
	userAddress := c.PostForm("userAddress")
	contractAddress := c.PostForm("contractAddress")
	amount := c.PostForm("amount")
	log.Printf("userAddress: %s, contractAddress: %s amount: %s", userAddress, contractAddress, amount)
	if userAddress == "" || contractAddress == "" || amount == "" {
		panic(&middleware.BusinessError{ReCode: enums.INVALID_PARAMETERS})
	}
	contractAddressCleaned := utils.CleanHexPrefix(contractAddress)
	userAddressCleaned := utils.CleanHexPrefix(userAddress)
	amountHex, err := utils.ToHexStringNoPrefixZeroPaddedFromString(amount, 64)
	if err != nil {
		log.Printf("Error converting amount to hex: %v", err)
		panic(&middleware.BusinessError{ReCode: enums.INVALID_PARAMETERS})
	}
	concat := strings.ToLower(userAddressCleaned + amountHex + contractAddressCleaned)
	sign := h.EncodeService.Sign(concat)
	c.JSON(http.StatusOK, model.OkWithData[string](sign))
}
