package db

import (
	"errors"
	"fintrack/internal/entity"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DBClient 定義資料庫客戶端接口
type DBClient interface {
	SaveTransaction(tx entity.Transaction) error
	GetFilteredTransactions(userID, category, startDate, endDate string, page, pageSize int) ([]entity.Transaction, error)
	GetTransactions(userID, startDate, endDate string) ([]entity.Transaction, error)
	DeleteTransactionByID(txID string) error
}

// MySQLClient 實現 DBClient 接口
type MySQLClient struct {
	DB *gorm.DB
}

// NewDBClient 創建並返回一個 DBClient 實例，並使用 MySQL 作為數據庫
func NewDBClient(dsn string) (DBClient, error) {
	return NewMySQLClient(dsn)
}

// NewMySQLClient 創建一個新的 MySQLClient
func NewMySQLClient(dsn string) (*MySQLClient, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自動遷移數據庫模型
	err = db.AutoMigrate(&entity.Transaction{})
	if err != nil {
		return nil, err
	}

	return &MySQLClient{DB: db}, nil
}

// SaveTransactions 批量保存交易紀錄
func (c *MySQLClient) SaveTransaction(tx entity.Transaction) error {
	return c.DB.Create(&tx).Error
}

// 查詢交易數據，支持分頁、篩選
func (c *MySQLClient) GetFilteredTransactions(userID, category, startDate, endDate string, page, pageSize int) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	query := c.DB.Where("user_id = ?", userID)

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
func (c *MySQLClient) GetTransactions(userID, startDate, endDate string) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := c.DB.Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).Find(&transactions).Error
	return transactions, err
}

// DeleteTransactionByID 根據交易 ID 刪除交易紀錄
func (c *MySQLClient) DeleteTransactionByID(txID string) error {
	result := c.DB.Delete(&entity.Transaction{}, "id = ?", txID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("transaction not found")
	}
	return nil
}
