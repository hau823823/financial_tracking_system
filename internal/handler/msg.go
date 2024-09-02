package handler

import (
	"fintrack/internal/mq"
	"fintrack/internal/service"
	"log"

	"github.com/streadway/amqp"
)

type MessageHandler struct {
	Service  service.MessageService
	Consumer mq.MQConsumer
}

func NewMessageHandler(s service.MessageService, consumer mq.MQConsumer) *MessageHandler {
	return &MessageHandler{Service: s, Consumer: consumer}
}

// Start 消費並處理 RabbitMQ 消息
func (h *MessageHandler) Start() {
	// 獲取 RabbitMQ 消費的消息管道
	msgs, err := h.Consumer.ConsumeMessages()
	if err != nil {
		log.Fatalf("Failed to start consuming messages: %v", err)
	}

	// 開始處理消息
	h.StartConsumingMessages(msgs)
}

// StartConsumingMessages 開始消費並處理 RabbitMQ 消息
func (h *MessageHandler) StartConsumingMessages(msgs <-chan amqp.Delivery) {
	msgs, err := h.Consumer.ConsumeMessages()
	if err != nil {
		log.Fatalf("Failed to start consuming messages: %v", err)
	}

	go func() {
		for d := range msgs {
			// 調用 service 處理具體業務邏輯
			if err := h.Service.ProcessTransaction(d.Body); err != nil {
				log.Printf("Failed to process message: %v", err)
			}
		}
	}()
}
