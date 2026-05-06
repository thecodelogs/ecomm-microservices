package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	UserServiceURL    string
	ProductServiceURL string
	JWTSecret         string
	Environment       string
}

func Load() *Config {
	_ = godotenv.Load(".env")

	cfg := &Config{
		Port:              getEnv("PORT", "8080"),
		UserServiceURL:    getEnv("USER_SERVICE_URL", "localhost:50051"),
		ProductServiceURL: getEnv("PRODUCT_SERVICE_URL", "localhost:50052"),
		JWTSecret:         getEnv("PASETO_SECRET", "change-this-to-32-characters!!"),
		Environment:       getEnv("ENV", "development"),
	}

	if len(cfg.JWTSecret) < 32 {
		log.Fatalf("JWT_SECRET must be at least 32 characters for PASETO, got %d", len(cfg.JWTSecret))
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
