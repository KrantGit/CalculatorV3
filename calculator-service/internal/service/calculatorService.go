package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"calculator-service/internal/repository/kafka"
)

var (
	ErrDivisionByZero = errors.New("division by zero")
	ErrUnknownOp      = errors.New("unknown operation")
)

type CalculatorService struct {
	producer kafka.Producer
}

func NewCalculatorService(producer kafka.Producer) *CalculatorService {
	return &CalculatorService{producer: producer}
}

func (s *CalculatorService) Process(ctx context.Context, a, b float64, op string) (string, error) {
	result, err := s.calculate(a, b, op)
	if err != nil {
		return "", err
	}

	msg := fmt.Sprintf(`{"expression":"%.2f %s %.2f","result":%.2f}`, a, op, b, result)

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.producer.Send(timeoutCtx, []byte(msg)); err != nil {
		return "", fmt.Errorf("failed to send to queue: %w", err)
	}

	return msg, nil
}

func (s *CalculatorService) calculate(a, b float64, op string) (float64, error) {
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		return a / b, nil
	default:
		return 0, ErrUnknownOp
	}
}
