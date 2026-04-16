package kafka

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type CalculationMessage struct {
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
}

type Consumer struct {
	reader *kafka.Reader
	db     *sql.DB
}

func NewConsumer(broker, topic string, db *sql.DB) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:   []string{broker},
			Topic:     topic,
			Partition: 0,
			MinBytes:  10e3,
			MaxBytes:  10e6,
		}),
		db: db,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Println("🔄 Kafka consumer started, waiting for messages...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping consumer...")
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("⚠️ ReadMessage error: %v", err)
				time.Sleep(time.Second)
				continue
			}

			var calcMsg CalculationMessage
			if err := json.Unmarshal(msg.Value, &calcMsg); err != nil {
				log.Printf("⚠️ JSON unmarshal error: %v", err)
				continue
			}

			log.Printf("📩 Received: %s = %.2f", calcMsg.Expression, calcMsg.Result)

			if err := c.saveToDB(calcMsg); err != nil {
				log.Printf("❌ DB save error: %v", err)
			} else {
				log.Printf("✅ Saved to DB: %s", calcMsg.Expression)
			}
		}
	}
}

func (c *Consumer) saveToDB(msg CalculationMessage) error {
	_, err := c.db.Exec(
		"INSERT INTO calculations (expression, result) VALUES ($1, $2)",
		msg.Expression, msg.Result,
	)
	return err
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
