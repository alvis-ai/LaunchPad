package handler

import (
	"launchpad/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AprHandler struct{}

func NewAprHandler() *AprHandler {
	return &AprHandler{}
}

// @Summary BRE Price
// @Description Get BRE/LaunchPad reference price used by the frontend APR widgets
// @Tags apr
// @Produce json
// @Success 200 {object} model.Result[string]
// @Router /boba/apr/bre_price [get]
func (h *AprHandler) BrePrice(c *gin.Context) {
	c.JSON(http.StatusOK, model.OkWithData[string]("0"))
}
