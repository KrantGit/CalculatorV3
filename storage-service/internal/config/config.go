package config

import (
	"fmt"
	"os"
)

type Config struct {
	Kafka      KafkaConfig
	Database   DatabaseConfig
	Migrations MigrationsConfig
}

type KafkaConfig struct {
	Broker string
	Topic  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type MigrationsConfig struct {
	Path string
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host, c.Database.Port, c.Database.User,
		c.Database.Password, c.Database.Name, c.Database.SSLMode)
}

func (c *Config) MigrateURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Database.User, c.Database.Password,
		c.Database.Host, c.Database.Port,
		c.Database.Name, c.Database.SSLMode)
}

func Load() *Config {
	return &Config{
		Kafka: KafkaConfig{
			Broker: getEnv("KAFKA_BROKER", "kafka:29092"),
			Topic:  getEnv("KAFKA_TOPIC", "calculations"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "postgres"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "dbpass"),
			Name:     getEnv("DB_NAME", "calculator"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Migrations: MigrationsConfig{
			Path: getEnv("MIGRATIONS_PATH", "/app/migrations"),
		},
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
