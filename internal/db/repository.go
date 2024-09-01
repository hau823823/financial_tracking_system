package db

import (
	"fintrack/internal/entity"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	SaveTransaction(tx entity.Transaction) error
	GetFilteredTransactions(userID, category, startDate, endDate string, page, pageSize int) ([]entity.Transaction, error)
	GetTransactions(userID, startDate, endDate string) ([]entity.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

// 儲存交易數據
func (r *transactionRepository) SaveTransaction(tx entity.Transaction) error {
	return r.db.Create(&tx).Error
}

// 查詢交易數據，支持分頁、篩選
func (r *transactionRepository) GetFilteredTransactions(userID, category, startDate, endDate string, page, pageSize int) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	query := r.db.Where("user_id = ?", userID)

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if startDate != "" && endDate != "" {
		query = query.Where("date BETWEEN ? AND ?", startDate, endDate)
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&transactions).Error
	return transactions, err
}

// 查詢指定範圍內的交易
func (r *transactionRepository) GetTransactions(userID, startDate, endDate string) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := r.db.Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).Find(&transactions).Error
	return transactions, err
}
