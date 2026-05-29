package handler

import (
	"fmt"
	"launchpad/api/request"
	"launchpad/api/response"
	"launchpad/enums"
	"launchpad/middleware"
	"launchpad/model"
	"launchpad/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductContractHandler struct {
	ProductContractService *service.ProductContractService
}

func NewProductContractHandler(productContractService *service.ProductContractService) *ProductContractHandler {
	return &ProductContractHandler{
		ProductContractService: productContractService,
	}
}

// @Summary Product Base Info
// @Description Get product base info
// @Tags product
// @Accept json
// @Produce json
// @Param productId query string true "Product Id"
// @Success 200 {object} model.Result[response.ProductContractResponse]
// @Router /boba/product/base_info [get]
func (h *ProductContractHandler) BaseInfo(c *gin.Context) {
	productId := c.Query("productId")
	if productId == "" {
		panic(&middleware.BusinessError{ReCode: enums.INVALID_PARAMETERS})
	}

	productContract, err := h.ProductContractService.GetById(productId)
	if err != nil {
		panic(&middleware.BusinessError{ReCode: enums.INVALID_PARAMETERS})
	}
	productContractResponse := response.ProductContractToProductContractResponse(*productContract)
	c.JSON(http.StatusOK, model.OkWithData(productContractResponse))
}

// @Summary Product List
// @Description Get product list
// @Tags product
// @Accept json
// @Produce json
// @Success 200 {object} model.Result[response.ProductContractResponse]
// @Router /boba/product/list [get]
func (h *ProductContractHandler) List(c *gin.Context) {
	productContracts, err := h.ProductContractService.List()
	if err != nil {
		fmt.Printf("List error: %v\n", err)
		panic(&middleware.BusinessError{ReCode: enums.INVALID_PARAMETERS})
	}
	fmt.Printf("Found %d product contracts\n", len(productContracts))
	// 初始化为空数组，确保返回 [] 而不是 null
	productContractResponses := make([]response.ProductContractResponse, 0)
	if productContracts != nil {
		for _, productContract := range productContracts {
			productContractResponses = append(productContractResponses, response.ProductContractToProductContractResponse(productContract))
		}
	}

	c.JSON(http.StatusOK, model.OkWithData(productContractResponses))
}

// @Summary Update Product Contract
// @Description Update product contract
// @Tags product
// @Accept json
// @Produce json
// @Param productId formData string true "Product Id"
// @Param name formData string true "Name"
// @Param description formData string true "Description"
// @Param price formData string true "Price"
// @Param status formData string true "Status"
// @Success 200 {object} model.Result[string]
// @Router /boba/update [post]
func (h *ProductContractHandler) Update(c *gin.Context) {
	var productContractUpdateRequest request.ProductContractUpdateRequest
	productId := c.PostForm("productId")
	fmt.Printf("productId: %s\n", productId)

	if err := c.ShouldBindBodyWithJSON(&productContractUpdateRequest); err != nil {
		fmt.Printf("Bind error: %v\n", err)
		panic(&middleware.BusinessError{ReCode: enums.INVALID_PARAMETERS})
	}

	fmt.Printf("productContractUpdateRequest: %+v\n", productContractUpdateRequest)

	if err := h.ProductContractService.Update(&productContractUpdateRequest); err != nil {
		fmt.Printf("Update error: %v\n", err)
		panic(&middleware.BusinessError{ReCode: enums.INVALID_PARAMETERS})
	}

	c.JSON(http.StatusOK, model.OkWithData[string]("success"))
}
