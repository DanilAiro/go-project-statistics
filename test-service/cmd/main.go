package main

import (
	"context"
	"log"
	"time"

	pb "api-contracts/gen/go/api/proto/account"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться: %v", err)
	}
	defer conn.Close()

	client := pb.NewAccountServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := client.GetAccountByName(ctx, &pb.AccountByNameRequest{Tag: "000", Name: "danilairo"})
	if err != nil {
		log.Fatalf("Ошибка вызова метода GetAccountByName: %v", err)
	}

	log.Printf("Ответ сервера -> PUUID: %s, Name: %s, Tag: %s, Region: %s", resp.GetPuuid(), resp.GetName(), resp.GetTag(), resp.GetRegion())
}
