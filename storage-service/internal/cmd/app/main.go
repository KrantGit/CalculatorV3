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
	"storage-service/internal/handler"
	"storage-service/internal/migrator"
	"storage-service/internal/repository"
	"storage-service/internal/service"
)

func main() {
	cfg := config.Load()

	log.Println("🔄 Applying database migrations...")
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
		log.Fatalf("DB ping failed: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	repo := repository.NewPostgresRepo(db)
	calcService := service.NewCalculationService(repo)
	calcHandler := handler.NewCalculationHandler(calcService, cfg.Kafka.Broker, cfg.Kafka.Topic)
	defer calcHandler.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("Storage Worker is running. Press Ctrl+C to stop.")
	if err := calcHandler.Start(ctx); err != nil && err != context.Canceled {
		log.Fatalf("Handler fatal error: %v", err)
	}
	log.Println("Graceful shutdown complete")
}
