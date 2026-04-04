package service

import (
	"calculator-service/internal/entity"
	"calculator-service/internal/repository/kafka"
)

type Service struct {
	kafka *kafka.Producer
}

func New(kafka *kafka.Producer) *Service {
	return &Service{kafka}
}

func (s *Service) CalculatorService(input entity.Input) entity.Output {

	var result entity.Output

	switch input.Sign {
	case "+":
		result.Result = input.FirstNumber + input.SecondNumber
	case "-":
		result.Result = input.FirstNumber - input.SecondNumber
	case "*":
		result.Result = input.FirstNumber * input.SecondNumber
	case "/":
		if input.SecondNumber == 0 {
			result.Error = "Division by zero"
		} else {
			result.Result = input.FirstNumber / input.SecondNumber
		}
	default:
		result.Error = "Unknown sign: " + input.Sign
	}

	s.kafka.Publish(input, result)

	return result
}
