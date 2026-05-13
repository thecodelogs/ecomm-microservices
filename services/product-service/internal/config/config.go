package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL        string
	Port         string
	PasetoSecret string
}

func Load() *Config {
	_ = godotenv.Load(".env")

	cfg := &Config{
		DBURL:        getEnv("DB_URL", ""),
		Port:         getEnv("PORT", "50052"),
		PasetoSecret: getEnv("PASETO_SECRET", "change-this-to-32-characters!!"),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
