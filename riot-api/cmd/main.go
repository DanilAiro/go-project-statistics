package main

import (
	"log"
	"net"
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

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Не удалось запустить listener: %v", err)
	}

	grpcServer := grpc.NewServer()
	accountServer := account.NewAccountServer(conf.ValorantAPIKey)
	matchServer := match.NewMatchServer(conf.ValorantAPIKey)
	mmrServer := mmr.NewMMRServer(conf.ValorantAPIKey)

	accountpb.RegisterAccountServiceServer(grpcServer, accountServer)
	matchpb.RegisterMatchServiceServer(grpcServer, matchServer)
	mmrpb.RegisterMMRServiceServer(grpcServer, mmrServer)

	log.Println("Сервер запущен")
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка запуска gRPC сервера: %v", err)
	}
}
