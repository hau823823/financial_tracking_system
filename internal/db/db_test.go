package db_test

import (
	"fintrack/internal/db"
	"fintrack/internal/entity"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New() // 使用 go-sqlmock 模擬數據庫
	assert.NoError(t, err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock
}

func TestSaveTransaction(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	client := &db.MySQLClient{DB: gormDB}

	transaction := entity.Transaction{
		ID:          "1",
		UserID:      "user123",
		Date:        time.Now(),
		Amount:      100.0,
		Category:    "INCOME",
		Description: "Salary",
		Source:      "MANUAL",
		Reconciled:  false,
	}

	// 設置預期的 INSERT SQL 行為
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `transactions`").
		WithArgs(transaction.ID, transaction.UserID, transaction.Date, transaction.Amount, transaction.Category, transaction.Description, transaction.Source, transaction.Reconciled).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := client.SaveTransaction(transaction)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFilteredTransactions(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	client := &db.MySQLClient{DB: gormDB}

	mockTransactions := []entity.Transaction{
		{ID: "1", UserID: "user123", Date: time.Now(), Amount: 100.0, Category: "INCOME", Description: "Salary", Source: "MANUAL", Reconciled: false},
		{ID: "2", UserID: "user123", Date: time.Now(), Amount: 50.0, Category: "EXPENSE", Description: "Groceries", Source: "CREDIT_CARD", Reconciled: true},
	}

	// 使用 regexp.QuoteMeta 包裹查詢語句，避免特殊字符被誤解，並匹配 LIMIT 子句
	query := regexp.QuoteMeta("SELECT * FROM `transactions` WHERE user_id = ? LIMIT ?")

	// 設置預期的 SELECT SQL 行為，包含 user_id 和 LIMIT 參數
	mock.ExpectQuery(query).
		WithArgs("user123", 10). // 添加 LIMIT 的參數
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "date", "amount", "category", "description", "source", "reconciled"}).
			AddRow(mockTransactions[0].ID, mockTransactions[0].UserID, mockTransactions[0].Date, mockTransactions[0].Amount, mockTransactions[0].Category, mockTransactions[0].Description, mockTransactions[0].Source, mockTransactions[0].Reconciled).
			AddRow(mockTransactions[1].ID, mockTransactions[1].UserID, mockTransactions[1].Date, mockTransactions[1].Amount, mockTransactions[1].Category, mockTransactions[1].Description, mockTransactions[1].Source, mockTransactions[1].Reconciled))

	// 調用 GetFilteredTransactions，包含分頁設置
	transactions, err := client.GetFilteredTransactions("user123", "", "", "", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, len(mockTransactions), len(transactions))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTransactions(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	client := &db.MySQLClient{DB: gormDB}

	mockTransactions := []entity.Transaction{
		{ID: "1", UserID: "user123", Date: time.Now(), Amount: 100.0, Category: "INCOME", Description: "Salary", Source: "MANUAL", Reconciled: false},
	}

	// 使用 regexp.QuoteMeta 包裹查詢語句，避免特殊字符被誤解
	query := regexp.QuoteMeta("SELECT * FROM `transactions` WHERE user_id = ? AND date BETWEEN ? AND ?")

	// 設置預期的 SELECT SQL 行為
	mock.ExpectQuery(query).
		WithArgs("user123", "2023-01-01", "2023-12-31").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "date", "amount", "category", "description", "source", "reconciled"}).
			AddRow(mockTransactions[0].ID, mockTransactions[0].UserID, mockTransactions[0].Date, mockTransactions[0].Amount, mockTransactions[0].Category, mockTransactions[0].Description, mockTransactions[0].Source, mockTransactions[0].Reconciled))

	transactions, err := client.GetTransactions("user123", "2023-01-01", "2023-12-31")
	assert.NoError(t, err)
	assert.Equal(t, len(mockTransactions), len(transactions))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteTransactionByID(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	client := &db.MySQLClient{DB: gormDB}

	// 設置預期的 DELETE SQL 行為
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `transactions` WHERE id = ?").
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := client.DeleteTransactionByID("1")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
