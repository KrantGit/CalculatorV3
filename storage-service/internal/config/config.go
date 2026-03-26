package config

import (
	"log"
	"os"
)

type Config struct {
	Port string
}

func Load() *Config {
	cfg := &Config{}
	cfg.Port = getEnv("PORT", "50051")
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
