package config

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type Config struct {
	Env      string `mapstructure:"ENV"`
	Port     string `mapstructure:"PORT"`
	Database struct {
		Host     string `mapstructure:"HOST"`
		Port     string `mapstructure:"PORT"`
		User     string `mapstructure:"USER"`
		Password string `mapstructure:"PASSWORD"`
		Name     string `mapstructure:"NAME"`
	} `mapstructure:"DATABASE"`
	Owner struct {
		PrivateKey string `mapstructure:"PRIVATE_KEY"`
		NetworkUrl string `mapstructure:"NETWORK_URL"`
	} `mapstructure:"OWNER"`
	Sales struct {
		SalesFactoryAddress string `mapstructure:"SALE_FACTORY_ADDRESS"`
	} `mapstructure:"SALES"`
}

var AppConfig Config

// LoadConfig 加载配置文件
func LoadConfig() {
	env := os.Getenv("GO_ENV")
	if env != "" {
		viper.SetConfigName("config." + env) // 配置文件名（不带扩展名）
	} else {
		viper.SetConfigName("config") // 配置
	}

	viper.SetConfigType("yml") // 配置文件类型

	configPath := findRoot() + "/config"
	viper.AddConfigPath(configPath) // 添加配置文件路径

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("Error reading config file, %s", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Errorf("Unable to decode into struct, %v", err)
	}
	applyEnvOverrides()
}

func applyEnvOverrides() {
	AppConfig.Env = envOrDefault("APP_ENV", AppConfig.Env)
	AppConfig.Port = envOrDefault("PORT", AppConfig.Port)
	AppConfig.Database.Host = envOrDefault("DB_HOST", AppConfig.Database.Host)
	AppConfig.Database.Port = envOrDefault("DB_PORT", AppConfig.Database.Port)
	AppConfig.Database.User = envOrDefault("DB_USER", AppConfig.Database.User)
	AppConfig.Database.Password = envOrDefault("DB_PASSWORD", AppConfig.Database.Password)
	AppConfig.Database.Name = envOrDefault("DB_NAME", AppConfig.Database.Name)
	AppConfig.Owner.PrivateKey = cleanPrivateKey(envOrDefault("OWNER_PRIVATE_KEY", AppConfig.Owner.PrivateKey))
	AppConfig.Owner.NetworkUrl = envOrDefault("OWNER_NETWORK_URL", AppConfig.Owner.NetworkUrl)
	AppConfig.Sales.SalesFactoryAddress = envOrDefault("SALES_FACTORY_ADDRESS", AppConfig.Sales.SalesFactoryAddress)
}

func envOrDefault(key string, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func cleanPrivateKey(value string) string {
	return strings.TrimPrefix(strings.TrimSpace(value), "0x")
}

func findRoot() string {
	if os.Getenv("GO_ENV") == "production" {
		exePath, err := os.Executable()
		if err != nil {
			panic(err)
		}
		return filepath.Dir(exePath)
	} else {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		// 不断向上查找 go.mod 文件
		for {
			if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
				return dir // 找到 go.mod 文件，返回当前目录
			}

			// 向上一级目录
			parent := filepath.Dir(dir)
			if parent == dir {
				// 已经到根目录，停止
				break
			}
			dir = parent
		}

		panic("could not find go.mod")
	}

}
