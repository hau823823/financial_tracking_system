// cmd/main.go
package main

import (
	"fintrack/internal/di"
	"log"
)

func main() {
	// 初始化 TransactionHandler
	transactionHandler, err := di.InitializeTransactionHandler()
	if err != nil {
		log.Fatalf("Failed to initialize transaction handler: %v", err)
	}

	// 初始化 MessageHandler
	messageHandler, err := di.InitializeMessageHandler()
	if err != nil {
		log.Fatalf("Failed to initialize message handler: %v", err)
	}

	// 啟動 HTTP 伺服器
	r := transactionHandler.SetupRouter() // 根據 transactionHandler 設定路由
	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	// 啟動 RabbitMQ 消費者，由 handler 負責整個消費和處理流程
	messageHandler.Start()

	// 保持主線程運行
	select {}
}
