package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"

	"storage-service/internal/entity"
	"storage-service/internal/service"
)

type CalculationHandler struct {
	service *service.CalculationService
	reader  *kafka.Reader
}

func NewCalculationHandler(
	svc *service.CalculationService,
	broker, topic string,
) *CalculationHandler {
	return &CalculationHandler{
		service: svc,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:   []string{broker},
			Topic:     topic,
			Partition: 0,
			MinBytes:  10e3,
			MaxBytes:  10e6,
		}),
	}
}

func (h *CalculationHandler) Start(ctx context.Context) error {
	log.Println("🔄 Kafka handler started, listening on topic...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Context cancelled, stopping handler")
			return ctx.Err()
		default:
			if err := h.processNextMessage(ctx); err != nil {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				log.Printf("Processing error: %v, retrying in 1s...", err)
				time.Sleep(time.Second)
			}
		}
	}
}

func (h *CalculationHandler) processNextMessage(ctx context.Context) error {
	msg, err := h.reader.ReadMessage(ctx)
	if err != nil {
		return fmt.Errorf("kafka read failed: %w", err)
	}

	var calc entity.Calculation
	if err := json.Unmarshal(msg.Value, &calc); err != nil {
		return fmt.Errorf("json unmarshal failed: %w", err)
	}

	log.Printf("Received: %s = %.2f", calc.Expression, calc.Result)

	if err := h.service.ProcessCalculation(ctx, &calc); err != nil {
		log.Printf("Service error: %v", err)
		return err
	}

	log.Printf("Processed: %s", calc.Expression)
	return nil
}

func (h *CalculationHandler) Close() error {
	return h.reader.Close()
}
