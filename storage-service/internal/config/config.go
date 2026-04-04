package config

import (
	"log"
	"os"
)

type Config struct {
	Host   string
	Port   string
	User   string
	Pass   string
	Dbname string
}

func Load() *Config {
	cfg := &Config{
		getEnv("DB_HOST", "postgres"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "dbpass"),
		getEnv("DB_NAME", "calculator"),
	}
	return cfg
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		if fallback == "" {
			log.Fatalf("env %s required", key)
		}
		return fallback
	}
	return val
}
