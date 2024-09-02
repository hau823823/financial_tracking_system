// internal/service/transaction_service.go
package service

import (
	"encoding/json"
	"errors"
	"fintrack/internal/db"
	"fintrack/internal/entity"
	"log"
)

type MessageService interface {
	ProcessTransaction(body []byte) error
	ValidateTransaction(tx *entity.Transaction) error
}

// MessageService 負責處理交易的業務邏輯
type messageService struct {
	dbClient db.DBClient
}

// NewMessageService 創建並返回 MessageService 實例
func NewMessageService(dbClient db.DBClient) MessageService {
	return &messageService{dbClient: dbClient}
}

// ProcessTransaction 處理 RabbitMQ 消息，解析並執行業務邏輯
func (s *messageService) ProcessTransaction(body []byte) error {
	var transaction entity.Transaction

	// 解析消息，將 JSON 轉換為交易實體
	if err := json.Unmarshal(body, &transaction); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return err
	}

	// 執行交易數據的驗證邏輯
	if err := s.ValidateTransaction(&transaction); err != nil {
		log.Printf("Transaction validation failed: %v", err)
		return err
	}

	// 設定對帳狀態
	transaction.Reconciled = true

	// 保存到資料庫
	if err := s.dbClient.SaveTransaction(transaction); err != nil {
		log.Printf("Failed to save transaction: %v", err)
		return err
	}

	log.Printf("Successfully processed transaction: %v", transaction)
	return nil
}

// validateTransaction 驗證交易記錄的正確性
func (s *messageService) ValidateTransaction(tx *entity.Transaction) error {
	if tx.Amount <= 0 {
		return errors.New("invalid transaction amount")
	}
	return nil
}
