package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DatabaseURL        string
	JWTSecret          string
	StripeSecretKey    string
	StripWebhookSecret string
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:               getEnv("PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		StripeSecretKey:    getEnv("STRIPE_SECRET_KEY", ""),
		StripWebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL URL is required")
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	return cfg
}
