package service

import (
	"context"
	"encoding/json"
	"errors"
	"fintrack/internal/cache"
	"fintrack/internal/db"
	"fintrack/internal/entity"
	"fintrack/internal/mq"
	"log"
	"time"
)

type TransactionService interface {
	AddTransaction(tx entity.Transaction) error
	GenerateReport(ctx context.Context, userID string, period string) (interface{}, error)
}

type transactionService struct {
	repo     db.TransactionRepository
	cache    cache.Cache
	producer mq.MQProducer
}

func NewTransactionService(repo db.TransactionRepository, cache cache.Cache, producer mq.MQProducer) TransactionService {
	return &transactionService{repo: repo, cache: cache, producer: producer}
}

// 新增交易，將寫入操作委派給 RabbitMQ 進行異步處理
func (s *transactionService) AddTransaction(tx entity.Transaction) error {
	if tx.Amount == 0 {
		return errors.New("transaction amount cannot be zero")
	}

	// 將交易數據轉換為 JSON 並推送到 RabbitMQ
	message, err := json.Marshal(tx)
	if err != nil {
		return err
	}

	if err := s.producer.SendMessage(message); err != nil {
		log.Printf("Failed to send transaction message to RabbitMQ: %v", err)
		return err
	}

	return nil
}

// 生成財務報表
func (s *transactionService) GenerateReport(ctx context.Context, userID string, period string) (interface{}, error) {
	cacheKey := "report:" + userID + ":" + period

	// 從緩存中獲取報表
	report, err := s.cache.Get(ctx, cacheKey)
	if err == nil && report != "" {
		log.Println("Cache hit: returning cached report")
		return report, nil
	}

	// 從資料庫查詢並生成報表
	transactions, err := s.repo.GetTransactions(userID, "2023-01-01", "2023-12-31")
	if err != nil {
		return nil, err
	}

	generatedReport := map[string]interface{}{
		"user":    userID,
		"period":  period,
		"entries": transactions,
	}

	// 將生成的報表存入緩存
	_ = s.cache.Set(ctx, cacheKey, generatedReport, 24*time.Hour)

	return generatedReport, nil
}
