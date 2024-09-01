package api

import (
	"fintrack/internal/entity"
	"fintrack/internal/service"
	"net/http"
	"strconv"

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

	// 數據驗證
	if tx.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than zero"})
		return
	}
	if tx.Date.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
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

// 查詢交易紀錄，支持分頁和篩選
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	userID := c.Query("user_id")
	category := c.Query("category")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	transactions, err := h.Service.GetTransactions(userID, category, startDate, endDate, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// 匯入銀行或信用卡帳單
func (h *TransactionHandler) ImportReconcile(c *gin.Context) {
	// 此處處理文件上傳和格式驗證，省略具體實現
	err := h.Service.ImportTransactions(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to import transactions"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "Transactions imported successfully"})
}

// 生成並查詢財務報表
func (h *TransactionHandler) GetReports(c *gin.Context) {
	userID := c.Query("user_id")
	reportType := c.Query("report_type")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	report, err := h.Service.GenerateReport(c.Request.Context(), userID, reportType, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report"})
		return
	}

	c.JSON(http.StatusOK, report)
}
