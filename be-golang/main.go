package main

import (
	"launchpad/api/handler"
	"launchpad/config"
	"launchpad/database"
	"launchpad/logger"
	"launchpad/repository"
	"launchpad/route"
	"launchpad/service"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

// @title LaunchPad API
// @version 1.0
// @description This is LaunchPad API server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api/v1

func main() {
	// 加载 .env 文件中的环境变量
	godotenv.Load()
	// 初始化日志记录器
	logger.InitLogger()
	// 加载配置文件
	config.LoadConfig()
	// 初始化数据库
	db := database.InitializeDB()
	// 实例化 repository
	productContractRepo := repository.NewProductContractRepository(db)

	// 实例化 service
	productContractService := service.NewProductContractService(productContractRepo)
	encodeService := service.NewEncodeService()

	//实例化 合约查询初始化
	saleContractService := service.NewSaleContractService(productContractService)
	go func() {
		if err := saleContractService.StartSaleFactoryListen(); err != nil {
			log.Errorf("sale factory listener stopped: %v", err)
		}
	}()

	// 实例化 handler
	productContractHandler := handler.NewProductContractHandler(productContractService)
	encodeHandler := handler.NewEncodeHandler(encodeService)
	helloHandler := handler.NewHelloHandler()
	aprHandler := handler.NewAprHandler()

	r := route.SetupRouter(
		encodeHandler,
		helloHandler,
		productContractHandler,
		aprHandler,
	)

	r.Run(":" + config.AppConfig.Port) // 启动服务
}
