//go:build wireinject
// +build wireinject

package di

import (
	"fintrack/config"
	"fintrack/internal/cache"
	"fintrack/internal/db"
	"fintrack/internal/handler"
	"fintrack/internal/mq"
	"fintrack/internal/service"
	"sync"

	"github.com/google/wire"
)

var (
	configOnce sync.Once
	cfg        *config.Config

	dbOnce   sync.Once
	database *db.MySQLClient

	rabbitMQOnce     sync.Once
	rabbitMQProducer *mq.RabbitMQClient
	rabbitMQConsumer *mq.RabbitMQClient

	redisOnce  sync.Once
	redisCache *cache.Redis
)

func NewConfig() *config.Config {
	configOnce.Do(func() {
		cfg = config.NewConfig()
	})
	return cfg
}

func NewDBClient(cfg *config.Config) (db.DBClient, error) {
	var err error
	dbOnce.Do(func() {
		database, err = db.NewMySQLClient(cfg.DSN)
	})
	return database, err
}

func NewRedisCache(cfg *config.Config) cache.Cache {
	redisOnce.Do(func() {
		redisCache = cache.NewRedis(cfg.RedisConfig)
	})
	return redisCache
}

// NewRabbitMQProducer 確保 RabbitMQ 生產者僅初始化一次
func NewRabbitMQProducer(cfg *config.Config) mq.MQProducer {
	rabbitMQOnce.Do(func() {
		rabbitMQProducer, _ = mq.NewRabbitMQClient(cfg.RabbitMQConfig)
	})
	return rabbitMQProducer
}

// NewRabbitMQConsumer 初始化 RabbitMQ 消費者
func NewRabbitMQConsumer(cfg *config.Config) mq.MQConsumer {
	rabbitMQOnce.Do(func() {
		rabbitMQConsumer,  _ = mq.NewRabbitMQClient(cfg.RabbitMQConfig)
	})
	return rabbitMQConsumer
}

// ProviderSet 定義所有的依賴提供者
var ProviderSet = wire.NewSet(
	NewConfig,                     // 單例模式加載配置
	NewDBClient,                   // 單例模式初始化資料庫
	NewRedisCache,                 // 單例模式初始化 Redis
	NewRabbitMQProducer,           // 單例模式初始化 RabbitMQ 生產者
	NewRabbitMQConsumer,           // 單例模式初始化 RabbitMQ 消費者
	service.NewTransactionService, // 初始化業務邏輯層
	handler.NewTransactionHandler, // 初始化 API 處理層
	service.NewMessageService,     // 初始化業務邏輯層
	handler.NewMessageHandler,     // 初始化消息處理層
)

func InitializeTransactionHandler() (*handler.TransactionHandler, error) {
	wire.Build(ProviderSet)
	return &handler.TransactionHandler{}, nil
}

func InitializeMessageHandler() (*handler.MessageHandler, error) {
	wire.Build(ProviderSet)
	return &handler.MessageHandler{}, nil
}
