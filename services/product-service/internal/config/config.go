package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL string
	Port  string
}

func Load() *Config {
	_ = godotenv.Load(".env")

	cfg := &Config{
		DBURL: getEnv("DB_URL", ""),
		Port:  getEnv("PORT", "50052"),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
