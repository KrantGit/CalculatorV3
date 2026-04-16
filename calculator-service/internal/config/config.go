package config

import (
	"os"
)

type Config struct {
	Server ServerConfig
	Kafka  KafkaConfig
}

type ServerConfig struct {
	Port string
}

type KafkaConfig struct {
	Broker string
	Topic  string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "9091"),
		},
		Kafka: KafkaConfig{
			Broker: getEnv("KAFKA_BROKER", "kafka:29092"),
			Topic:  getEnv("KAFKA_TOPIC", "calculations"),
		},
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
