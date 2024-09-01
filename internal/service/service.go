package service

import (
	"context"
	"encoding/json"
	"errors"
	"fintrack/internal/cache"
	"fintrack/internal/db"
	"fintrack/internal/entity"
	"fintrack/internal/mq"
	"io"
	"log"
	"time"
)

type TransactionService interface {
	AddTransaction(tx entity.Transaction) error
	GetTransactions(userID, category, startDate, endDate string, page, pageSize int) ([]entity.Transaction, error)
	ImportTransactions(data io.Reader) error
	GenerateReport(ctx context.Context, userID, reportType, startDate, endDate string) (interface{}, error)
}

type transactionService struct {
	repo     db.DBClient
	cache    cache.Cache
	producer mq.MQProducer
}

func NewTransactionService(repo db.DBClient, cache cache.Cache, producer mq.MQProducer) TransactionService {
	return &transactionService{repo: repo, cache: cache, producer: producer}
}

// 新增交易紀錄，將寫入操作委派給 RabbitMQ 進行異步處理
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

// 查詢交易紀錄
func (s *transactionService) GetTransactions(userID, category, startDate, endDate string, page, pageSize int) ([]entity.Transaction, error) {
	// 分頁查詢交易記錄，支持篩選條件
	return s.repo.GetFilteredTransactions(userID, category, startDate, endDate, page, pageSize)
}

// 匯入帳單交易
func (s *transactionService) ImportTransactions(data io.Reader) error {
	// 實現匯入邏輯，包括文件解析、數據格式檢查、以及自動對帳
	// 此處省略具體實現，可根據業務需要擴展
	return nil
}

// 生成財務報表
func (s *transactionService) GenerateReport(ctx context.Context, userID, reportType, startDate, endDate string) (interface{}, error) {
	cacheKey := "report:" + userID + ":" + reportType + ":" + startDate + ":" + endDate

	// 從緩存中獲取報表
	report, err := s.cache.Get(ctx, cacheKey)
	if err == nil && report != "" {
		log.Println("Cache hit: returning cached report")
		return report, nil
	}

	// 根據日期範圍查詢並生成報表
	transactions, err := s.repo.GetTransactions(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	generatedReport := map[string]interface{}{
		"user":    userID,
		"period":  startDate + " - " + endDate,
		"entries": transactions,
	}

	// 將生成的報表存入緩存
	_ = s.cache.Set(ctx, cacheKey, generatedReport, 24*time.Hour)

	return generatedReport, nil
}
