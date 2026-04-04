package main

import (
	"calculator-service/internal/config"
	"calculator-service/internal/grpcclient"
	"calculator-service/internal/handler"
	"calculator-service/internal/repository/kafka"
	"calculator-service/internal/service"
	"net/http"
)

func main() {
	cfg := config.Load()

	producer := kafka.New(cfg.Kafka.Broker, cfg.Kafka.Topic)

	calculatorService := service.New()

	calculatorHandler := handler.New(calculatorService)

	http.HandleFunc("/calculator", calculatorHandler.)
	if err := http.ListenAndServe(":9091", nil); err != nil {
		panic(err)
	}

}
