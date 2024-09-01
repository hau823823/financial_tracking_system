package db

import (
	"fintrack/internal/entity"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	SaveTransaction(tx entity.Transaction) error
	GetTransactions(userID string, startDate, endDate string) ([]entity.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) SaveTransaction(tx entity.Transaction) error {
	return r.db.Create(&tx).Error
}

func (r *transactionRepository) GetTransactions(userID string, startDate, endDate string) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := r.db.Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).Find(&transactions).Error
	return transactions, err
}