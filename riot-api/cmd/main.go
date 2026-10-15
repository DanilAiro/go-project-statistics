package main

import (
	"log"
	"net"
	"riot-api/internal/account"
	"riot-api/internal/config"

	pb "api-contracts/gen/go/api/proto/account"

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

	pb.RegisterAccountServiceServer(grpcServer, accountServer)

	log.Println("Сервер запущен")
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка запуска gRPC сервера: %v", err)
	}
}
