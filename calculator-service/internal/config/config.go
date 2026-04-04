package config

import (
	"log"
	"os"
)

type Config struct {
	Kafka struct {
		Broker string
		Topic  string
	}
}

func Load() *Config {
	cfg := &Config{}

	cfg.Kafka.Broker = getEnv("KAFKA_BROKER", "localhost:9092")
	cfg.Kafka.Topic = getEnv("KAFKA_TOPIC", "calculations")

	return cfg
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		if fallback == "" {
			log.Fatalf("env %s is required", key)
		}
		return fallback
	}
	return val
}
