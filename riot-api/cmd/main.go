package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"riot-api/internal/account"
	"riot-api/internal/config"
	"riot-api/internal/match"
	"riot-api/internal/mmr"

	accountpb "api-contracts/gen/go/api/proto/account"
	matchpb "api-contracts/gen/go/api/proto/match"
	mmrpb "api-contracts/gen/go/api/proto/mmr"

	"google.golang.org/grpc"
)

func main() {
	config.LoadEnvVariables()
	conf := config.NewConfig()

	if conf.LogPath == "" {
		conf.LogPath = "log"
	}

	logFile, err := os.OpenFile(conf.LogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		msg := fmt.Sprintf("Не удалось запустить logger: %v", err)
		slog.Warn(msg)
	} else {
		logger := slog.New(slog.NewJSONHandler(logFile, nil))
		slog.SetDefault(logger)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		msg := fmt.Sprintf("Не удалось запустить listener: %v", err)
		slog.Error(msg)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	accountServer := account.NewAccountServer(conf.ValorantAPIKey)
	matchServer := match.NewMatchServer(conf.ValorantAPIKey)
	mmrServer := mmr.NewMMRServer(conf.ValorantAPIKey)

	accountpb.RegisterAccountServiceServer(grpcServer, accountServer)
	matchpb.RegisterMatchServiceServer(grpcServer, matchServer)
	mmrpb.RegisterMMRServiceServer(grpcServer, mmrServer)

	slog.Info("Сервер запущен")
	if err = grpcServer.Serve(lis); err != nil {
		msg := fmt.Sprintf("Ошибка запуска gRPC сервера: %v", err)
		slog.Error(msg)
		os.Exit(1)
	}
}
