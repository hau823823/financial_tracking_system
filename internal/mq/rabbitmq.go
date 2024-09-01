package mq

import (
	"encoding/json"
	"fintrack/internal/db"
	"fintrack/internal/entity"
	"log"

	"github.com/streadway/amqp"
)

// MQ 生產者接口
type MQProducer interface {
	SendMessage(body []byte) error
	Close()
}

// RabbitMQ 生產者實現
type RabbitMQProducer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
}

func NewRabbitMQProducer(url, queueName string) (*RabbitMQProducer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQProducer{
		connection: conn,
		channel:    ch,
		queueName:  queueName,
	}, nil
}

func (p *RabbitMQProducer) SendMessage(body []byte) error {
	err := p.channel.Publish(
		"",
		p.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Printf("Failed to send message to RabbitMQ: %v", err)
	}
	return err
}

func (p *RabbitMQProducer) Close() {
	p.channel.Close()
	p.connection.Close()
}

// RabbitMQ 消費者實現
type RabbitMQConsumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
	repo       db.TransactionRepository
}

func NewRabbitMQConsumer(url, queueName string, repo db.TransactionRepository) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQConsumer{
		connection: conn,
		channel:    ch,
		queueName:  queueName,
		repo:       repo,
	}, nil
}

func (c *RabbitMQConsumer) ConsumeMessages() {
	msgs, err := c.channel.Consume(
		c.queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to start consuming messages: %v", err)
		return
	}

	for msg := range msgs {
		var tx entity.Transaction
		if err := json.Unmarshal(msg.Body, &tx); err != nil {
			log.Printf("Failed to unmarshal transaction message: %v", err)
			continue
		}

		// 將交易數據寫入資料庫
		if err := c.repo.SaveTransaction(tx); err != nil {
			log.Printf("Failed to save transaction to database: %v", err)
		}
	}
}

func (c *RabbitMQConsumer) Close() {
	c.channel.Close()
	c.connection.Close()
}
