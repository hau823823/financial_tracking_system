package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AppConfig struct {
	DB          *gorm.DB
	RabbitMQURL string
	RedisAddr   string
	RedisPass   string
}

func LoadConfig() (*AppConfig, error) {
	err := godotenv.Load("configs/app.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	config := &AppConfig{
		DB:          db,
		RabbitMQURL: os.Getenv("RABBITMQ_URL"),
		RedisAddr:   os.Getenv("REDIS_ADDR"),
		RedisPass:   os.Getenv("REDIS_PASS"),
	}

	return config, nil
}
