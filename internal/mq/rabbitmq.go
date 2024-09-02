package mq

import (
	"fintrack/config"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// MQProducer 定義生產者接口
type MQProducer interface {
	SendMessage(body []byte) error
	Close() error
}

// MQConsumer 定義 RabbitMQ 消費者接口
type MQConsumer interface {
	ConsumeMessages() (<-chan amqp.Delivery, error)
}

// RabbitMQClient 實現了 MQProducer 和 MQConsumer 接口
type RabbitMQClient struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

// NewMQProducer 創建並返回一個 MQProducer 實例，並使用 RabbitMQ 作為消息隊列
func NewMQProducer(config config.RabbitMQConfig) (MQProducer, error) {
	return NewRabbitMQClient(config)
}

func NewMQConsumer(config config.RabbitMQConfig) (MQProducer, error)  {
	return NewRabbitMQClient(config)
}

// NewRabbitMQProducer 創建並返回一個 RabbitMQ 實例
func NewRabbitMQClient(config config.RabbitMQConfig) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	q, err := ch.QueueDeclare(
		config.Queue, // 使用配置中的 QueueName
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &RabbitMQClient{
		connection: conn,
		channel:    ch,
		queue:      q,
	}, nil
}

// SendMessage 發送消息到 RabbitMQ
func (c *RabbitMQClient) SendMessage(body []byte) error {
	err := c.channel.Publish(
		"",           // exchange
		c.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Printf("Failed to send message to RabbitMQ: %v", err)
	}
	return err
}

// ConsumeMessages 消費 RabbitMQ 隊列中的消息
func (c *RabbitMQClient) ConsumeMessages() (<-chan amqp.Delivery, error) {
	msgs, err := c.channel.Consume(
		c.queue.Name, // queue
		"",                  // consumer
		true,                // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume messages: %w", err)
	}

	return msgs, nil
}

// Close 關閉 RabbitMQ 連接和通道
func (c *RabbitMQClient) Close() error {
	if err := c.channel.Close(); err != nil {
		log.Printf("Failed to close channel: %v", err)
		return err
	}
	if err := c.connection.Close(); err != nil {
		log.Printf("Failed to close connection: %v", err)
		return err
	}
	return nil
}
