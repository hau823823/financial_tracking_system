package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fintrack/internal/entity"
	"fintrack/internal/handler"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTransactionService 用於模擬 TransactionService
type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) AddTransaction(tx entity.Transaction) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockTransactionService) GetTransactions(userID, category, startDate, endDate string, page, pageSize int) ([]entity.Transaction, error) {
	args := m.Called(userID, category, startDate, endDate, page, pageSize)
	return args.Get(0).([]entity.Transaction), args.Error(1)
}

func (m *MockTransactionService) ImportTransactions(data io.Reader) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockTransactionService) GenerateReport(ctx context.Context, userID, reportType, startDate, endDate string) (interface{}, error) {
	args := m.Called(ctx, userID, reportType, startDate, endDate)
	return args.Get(0), args.Error(1)
}

func TestAddTransaction(t *testing.T) {
	// 設置 Gin 模式為測試模式
	gin.SetMode(gin.TestMode)

	mockService := new(MockTransactionService)
	handler := handler.NewTransactionHandler(mockService)

	// 模擬交易數據
	tx := entity.Transaction{
		ID:          "1",
		UserID:      "user123",
		Date:        time.Now().UTC(), // 保持時間一致
		Amount:      100.0,
		Category:    "INCOME",
		Description: "Salary",
		Source:      "MANUAL",
	}

	// 使用 mock.MatchedBy 進行靈活比較
	mockService.On("AddTransaction", mock.MatchedBy(func(transaction entity.Transaction) bool {
		return transaction.ID == tx.ID &&
			transaction.UserID == tx.UserID &&
			transaction.Amount == tx.Amount &&
			transaction.Category == tx.Category &&
			transaction.Description == tx.Description &&
			transaction.Source == tx.Source
	})).Return(nil)

	// 構建 HTTP 請求
	reqBody, _ := json.Marshal(tx)
	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	// 設置測試路由並執行測試請求
	router := gin.New()
	router.POST("/transactions", handler.AddTransaction)
	router.ServeHTTP(w, req)

	// 檢查狀態碼
	assert.Equal(t, http.StatusAccepted, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetTransactions(t *testing.T) {
	mockService := new(MockTransactionService)
	handler := handler.NewTransactionHandler(mockService)

	transactions := []entity.Transaction{
		{ID: "1", UserID: "user123", Date: time.Now(), Amount: 100.0, Category: "INCOME", Description: "Salary"},
	}

	mockService.On("GetTransactions", "user123", "", "", "", 1, 10).Return(transactions, nil)

	req := httptest.NewRequest(http.MethodGet, "/transactions?user_id=user123", nil)
	w := httptest.NewRecorder()

	router := handler.SetupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
