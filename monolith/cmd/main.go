package main

import (
	"fmt"
	"log/slog"
	"monolith/internal/config"
	"monolith/internal/db"
	"os"
)

func main() {
	config.LoadEnvVariables()
	conf := config.NewConfig()

	logFile, err := os.OpenFile(conf.LogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		msg := fmt.Sprintf("Не удалось запустить logger: %v", err)
		slog.Warn(msg)
	} else {
		logger := slog.New(slog.NewJSONHandler(logFile, nil))
		slog.SetDefault(logger)
	}

	fmt.Println(conf.DBUrl)

	DB, err := db.ConnectDB(conf.DBUrl)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	DB.Close()
}
