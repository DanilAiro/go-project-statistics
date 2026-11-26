package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	LogPath string
	DBUrl   string
}

// заполнение конфига
func NewConfig() *Config {
	logPath := os.Getenv("LOG_PATH")
	if logPath == "" {
		logPath = "log"
	}

	postgresUser := os.Getenv("POSTGRES_USER")
	if postgresUser == "" {
		log.Fatalln("нет имени пользователя для базы данных")
	}

	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	if postgresPassword == "" {
		log.Fatalln("нет пароля для базы данных")
	}

	postgresIP := os.Getenv("POSTGRES_IP")
	if postgresIP == "" {
		log.Fatalln("нет названия базы данных")
	}

	postgresPort := os.Getenv("POSTGRES_PORT")
	if postgresPort == "" {
		log.Fatalln("нет порта для базы данных")
	}

	postgresDB := os.Getenv("POSTGRES_DB")
	if postgresDB == "" {
		log.Fatalln("нет названия базы данных")
	}

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", postgresUser, postgresPassword, postgresIP, postgresPort, postgresDB)

	return &Config{
		LogPath: logPath,
		DBUrl:   dbUrl,
	}
}

// логика выгрузки .env
func LoadEnvVariables() {
	_ = godotenv.Load()
}
