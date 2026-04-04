package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"storage-service/internal/config"
	"storage-service/internal/entity"

	"github.com/segmentio/kafka-go"
)

type Msg struct {
	Expression entity.Input  `json:"expression"`
	Result     entity.Output `json:"result"`
}

func main() {
	cfg := config.Load()

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Pass,
		cfg.Dbname,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to DB")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "calculations",
	})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Read error: %v", err)
			continue
		}

		var calc Msg
		if err := json.Unmarshal(msg.Value, &calc); err != nil {
			log.Printf("JSON error: %v", err)
			continue
		}

		_, err = db.Exec("INSERT INTO calculations (expression, result) VALUES ($1, $2)",
			calc.Expression, calc.Result)
		if err != nil {
			log.Printf("DB error: %v", err)
		} else {
			log.Printf("Saved: %s = %.2f", calc.Expression, calc.Result)
		}
	}
}
