package api

import (
	"fintrack/internal/entity"
	"fintrack/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	Service service.TransactionService
}

func NewTransactionHandler(s service.TransactionService) *TransactionHandler {
	return &TransactionHandler{Service: s}
}

// 接收用戶的交易記錄並將其發送至 RabbitMQ
func (h *TransactionHandler) AddTransaction(c *gin.Context) {
	var tx entity.Transaction
	if err := c.ShouldBindJSON(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 非同步處理寫入資料庫
	err := h.Service.AddTransaction(tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Transaction received and will be processed"})
}
