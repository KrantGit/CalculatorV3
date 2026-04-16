package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"storage-service/internal/config"
	"storage-service/internal/migrator"
	"storage-service/internal/repository/kafka"
)

func main() {
	cfg := config.Load()

	log.Println("Applying database migrations")
	if err := migrator.RunMigrations(cfg.MigrateURL(), cfg.Migrations.Path); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	consumer := kafka.NewConsumer(cfg.Kafka.Broker, cfg.Kafka.Topic, db)
	defer consumer.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go consumer.Start(ctx)

	log.Println("DB Service is running. Press Ctrl C to stop.")
	<-ctx.Done()

	log.Println("Graceful shutdown complete")
}
