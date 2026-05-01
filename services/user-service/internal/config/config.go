package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL         string
	PASETO_SECRET string
	Port          string
}

func Load() *Config {
	_ = godotenv.Load(".env")

	secret := getEnv("PASETO_SECRET", "change-this-to-32-characters!!")

	// PASETO v2.local requires exactly 32 bytes
	if len(secret) < 32 {
		log.Fatalf("PASETO_SECRET must be at least 32 characters for PASETO, got %d", len(secret))
	}

	cfg := &Config{
		DBURL:         getEnv("DB_URL", "postgres://postgres:postgres@localhost:5432/user_db?sslmode=disable"),
		PASETO_SECRET: secret[:32], // truncate to exactly 32 bytes
		Port:          getEnv("PORT", "50051"),
	}

	if cfg.PASETO_SECRET == "change-this-to-32-characters!!" {
		log.Println("WARNING: Using default PASETO secret. Change in production!")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
