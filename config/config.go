// config/config.go
package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)
type Config struct {
	Port           string
	Mode           string
	DSN            string
	RabbitMQConfig RabbitMQConfig
	RedisConfig    RedisConfig
}

// RabbitMQConfig 包含 RabbitMQ 的相關配置
type RabbitMQConfig struct {
	URL   string
	Queue string
}

// RedisConfig 包含 Redis 的相關配置
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// init 函數用於設置環境變量和讀取配置文件
func init() {
	// 設置默認時區
	err := os.Setenv("TZ", "UTC")
	if err != nil {
		panic(fmt.Errorf("fatal error configs file: set time zone to UTC: %w", err))
	}

	// 設置 viper 的配置源
	viper.AutomaticEnv()
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../../configs")

	// 讀取配置文件
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("讀取設定檔出現錯誤，原因為：%v", err))
	}
}

func NewConfig() *Config {
	return &Config{
		Port: viper.GetString("PORT"),
		Mode: viper.GetString("MODE"),
		DSN:  viper.GetString("DSN"),
		RabbitMQConfig: RabbitMQConfig{
			URL:   viper.GetString("RABBITMQ_URL"),
			Queue: viper.GetString("RABBITMQ_QUEUE"),
		},
		RedisConfig: RedisConfig{
			Addr:     viper.GetString("REDIS_ADDR"),
			Password: viper.GetString("REDIS_PASS"),
			DB:       viper.GetInt("REDIS_DB"),
		},
	}
}
