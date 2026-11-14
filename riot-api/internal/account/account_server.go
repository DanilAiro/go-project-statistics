package account

import (
	"context"
	"encoding/json"

	"riot-api/internal/api"

	pb "api-contracts/gen/go/api/proto/account"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	apiKey string
	pb.UnimplementedAccountServiceServer
}

func NewAccountServer(apiKey string) *server {
	return &server{
		apiKey: apiKey,
	}
}

func (s *server) GetAccount(ctx context.Context, req *pb.AccountRequest) (*pb.Account, error) {
	puuid := req.GetPuuid()
	if puuid == "" {
		return nil, status.Error(codes.InvalidArgument, "PUUID пустой")
	}

	params := map[string]string{
		"puuid": puuid,
	}

	data, err := api.GetAccount(s.apiKey, params, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка запроса аккаунта: %v", err)
	}

	return decodeAccount(data)
}

func (s *server) GetAccountByName(ctx context.Context, req *pb.AccountByNameRequest) (*pb.Account, error) {
	name := req.GetName()
	tag := req.GetTag()

	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "Name пустой")
	}
	if tag == "" {
		return nil, status.Error(codes.InvalidArgument, "Tag пустой")
	}

	params := map[string]string{
		"name": name,
		"tag":  tag,
	}

	data, err := api.GetAccountByName(s.apiKey, params, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка запроса аккаунта: %v", err)
	}

	return decodeAccount(data)
}

func decodeAccount(data []byte) (*pb.Account, error) {
	var ac Account
	if err := json.Unmarshal(data, &ac); err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка разбора ответа: %v", err)
	}

	return ToProto(ac), nil
}

// ToProto конвертирует доменную структуру в proto-сообщение.
// Экспортирована, чтобы mmr-сервер мог переиспользовать её.
func ToProto(ac Account) *pb.Account {
	return &pb.Account{
		Tag:          ac.Tag,
		Name:         ac.Name,
		Card:         ac.Card,
		Puuid:        ac.PUUID,
		Title:        ac.Title,
		Region:       ac.Region,
		UpdatedAt:    ac.UpdatedAt,
		Platforms:    ac.Platforms,
		AccountLevel: ac.AccountLevel,
	}
}
