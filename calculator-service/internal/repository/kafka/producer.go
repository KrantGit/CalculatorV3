package kafka

import (
	"calculator-service/internal/entity"
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type Event struct {
	Expression entity.Input  `json:"expression"`
	Result     entity.Output `json:"result"`
	Timestamp  time.Time     `json:"timestamp"`
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

func (p *Producer) Publish(expr entity.Input, result entity.Output) error {

	event := Event{
		Expression: expr,
		Result:     result,
		Timestamp:  time.Now(),
	}

	data, _ := json.Marshal(event)

	err := p.writer.WriteMessages(context.Background(),
		kafka.Message{Value: data},
	)

	return err
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
