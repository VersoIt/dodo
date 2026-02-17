package config

import "os"

type Config struct {
	DatabaseURL string
	NatsURL     string
	AppPort     string
	JWTSecret   string
}

func NewConfig() *Config {
	getEnv := func(key, fallback string) string {
		if value, ok := os.LookupEnv(key); ok {
			return value
		}
		return fallback
	}

	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://pizza_user:pizza_password@localhost:5432/pizza_db?sslmode=disable"),
		NatsURL:     getEnv("NATS_URL", "nats://localhost:4222"),
		AppPort:     getEnv("APP_PORT", "8089"),
		JWTSecret:   getEnv("JWT_SECRET", "super-secret-key"),
	}
}
