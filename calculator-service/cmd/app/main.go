package main

import (
	"fmt"
	"log"
	"net/http"

	"calculator-service/internal/config"
	"calculator-service/internal/handler"
	"calculator-service/internal/repository/kafka"
	"calculator-service/internal/service"
)

func main() {
	cfg := config.Load()
	
	producer := kafka.NewKafkaProducer(cfg.Kafka.Broker, cfg.Kafka.Topic)
	defer producer.Close()

	calcService := service.NewCalculatorService(producer)

	calcHandler := handler.NewCalculatorHandler(calcService)

	http.HandleFunc("/calculator", calcHandler.HandleCalculate)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Calculator API starting on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
