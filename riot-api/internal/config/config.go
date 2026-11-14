package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ValorantAPIKey string
	LogPath        string
}

// заполнение конфига
func NewConfig() *Config {
	return &Config{
		ValorantAPIKey: os.Getenv("ValorantAPIKey"),
		LogPath:        os.Getenv("LogPath"),
	}
}

// логика выгрузки .env
func LoadEnvVariables() {
	_ = godotenv.Load()
}
