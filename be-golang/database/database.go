package database

import (
	"fmt"
	"launchpad/config"
	"time"

	log "github.com/sirupsen/logrus"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitializeDB() *gorm.DB {
	// 读取数据库配置
	dbUser := config.AppConfig.Database.User         // 数据库用户名
	dbPassword := config.AppConfig.Database.Password // 数据库密码
	dbHost := config.AppConfig.Database.Host         // 数据库主机地址
	dbPort := config.AppConfig.Database.Port         // 数据库端口
	dbName := config.AppConfig.Database.Name         // 数据库名称

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	gormLogger := &LogrusLogger{Log: log.StandardLogger()}

	// 重试连接数据库，最多重试30次，每次间隔2秒
	var db *gorm.DB
	var err error
	maxRetries := 30
	retryInterval := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 禁用表名复数
			},
			Logger: gormLogger.LogMode(logger.Info), // 设置日志级别为 Info，显示所有 SQL 查询
		})
		if err == nil {
			// 测试连接是否真正可用
			sqlDB, pingErr := db.DB()
			if pingErr != nil {
				err = pingErr
			} else {
				if pingErr = sqlDB.Ping(); pingErr == nil {
					log.Info("Database connection established successfully")
					break
				} else {
					err = pingErr
				}
			}
		}

		if err != nil {
			log.Warnf("Failed to connect database (attempt %d/%d): %v, retrying in %v...", i+1, maxRetries, err, retryInterval)
			if i < maxRetries-1 {
				time.Sleep(retryInterval)
			}
		}
	}

	if err != nil {
		log.Fatalf("Failed to connect database after %d attempts: %v", maxRetries, err)
	}

	// 配置数据库连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to configure database pool: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)           // 空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接可复用的最大时间

	DB = db

	log.Info("Database connection established and configured")
	return DB
}
