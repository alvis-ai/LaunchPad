package route

import (
	"launchpad/api/handler"
	"launchpad/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "launchpad/docs"
)

func SetupRouter(
	encodeHandler *handler.EncodeHandler,
	helloHandler *handler.HelloHandler,
	productContractHandler *handler.ProductContractHandler,
	aprHandler *handler.AprHandler,
) *gin.Engine {
	r := gin.New()
	// 使用 CORS 中间件处理跨域
	r.Use(middleware.CORS())
	// 使用自定义的 Recovery 中间件
	r.Use(middleware.CustomRecovery())

	// 注册路由
	r.GET("/boba/hello", helloHandler.GetHelloWord)
	r.POST("/boba/encode/sign_registration", encodeHandler.SignRegistration)
	r.POST("/boba/encode/sign_participation", encodeHandler.SignParticipation)
	r.GET("/boba/product/base_info", productContractHandler.BaseInfo)
	r.GET("/boba/product/list", productContractHandler.List)
	r.POST("/boba/update", productContractHandler.Update)
	r.GET("/boba/apr/bre_price", aprHandler.BrePrice)
	r.GET("/boba/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
