package config

import (
	"os"
)

type Config struct {
	ServerPort      string
	DatabaseURL     string
	SpellcheckerURL string
}

func NewConfig() *Config {
	return &Config{
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/notes?sslmode=disable"),
		SpellcheckerURL: getEnv("SPELLCHECKER_URL", "https://speller.yandex.net/services/spellservice.json/checkText"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
