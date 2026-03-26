package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"calculator-service/internal/config"
	"github.com/segmentio/kafka-go"
)

type Event struct {
	Expression string    `json:"expression"`
	Result     float64   `json:"result"`
	Timestamp  time.Time `json:"timestamp"`
}

type Producer struct {
	writer *kafka.Writer
}

func New(broker, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) Publish(expr string, result float64) {
	event := Event{
		Expression: expr,
		Result:     result,
		Timestamp:  time.Now(),
	}

	data, _ := json.Marshal(event)

	err := p.writer.WriteMessages(context.Background(),
		kafka.Message{Value: data},
	)

	if err != nil {
		log.Println("kafka error:", err)
	}
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
