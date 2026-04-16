package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type Producer interface {
	Send(ctx context.Context, msg []byte) error
	Close() error
}

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(broker, topic string) Producer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *KafkaProducer) Send(ctx context.Context, msg []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{Value: msg})
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
