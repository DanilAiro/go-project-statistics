package mmr

import (
	"context"
	"encoding/json"

	"riot-api/internal/account"
	"riot-api/internal/api"

	pb "api-contracts/gen/go/api/proto/mmr"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	apiKey string
	pb.UnimplementedMMRServiceServer
}

func NewMMRServer(apiKey string) *server {
	return &server{
		apiKey: apiKey,
	}
}

func (s *server) GetMMR(ctx context.Context, req *pb.MMRRequest) (*pb.MMR, error) {
	puuid := req.GetPuuid()
	affinity := req.GetAffinity()
	platform := req.GetPlatform()

	if puuid == "" {
		return nil, status.Error(codes.InvalidArgument, "PUUID пустой")
	}
	if affinity == "" {
		return nil, status.Error(codes.InvalidArgument, "Affinity пустой")
	}

	switch platform {
	case "":
		platform = "pc"
	case "pc", "console":
	default:
		return nil, status.Error(codes.InvalidArgument, "Platform должен быть pc или console")
	}

	params := map[string]string{
		"affinity": affinity,
		"platform": platform,
		"puuid":    puuid,
	}

	data, err := api.GetPlayerMMR(s.apiKey, params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка запроса MMR: %v", err)
	}

	var m MMR
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка разбора ответа: %v", err)
	}

	return toProto(m), nil
}

func toProto(m MMR) *pb.MMR {
	seasonal := make([]*pb.Seasonal, 0, len(m.Seasonal))
	for _, s := range m.Seasonal {
		actWins := make([]*pb.ActWins, 0, len(s.ActWins))
		for _, a := range s.ActWins {
			actWins = append(actWins, &pb.ActWins{Id: a.ID, Name: a.Name})
		}

		seasonal = append(seasonal, &pb.Seasonal{
			Wins:                 s.Wins,
			Games:                s.Games,
			EndRr:                s.EndRR,
			Season:               seasonToProto(s.Season),
			EndTier:              tierToProto(s.EndTier),
			ActWins:              actWins,
			RankingSchema:        s.RankingSchema,
			LeaderboardPlacement: placementToProto(s.LeaderboardPlacement),
		})
	}

	return &pb.MMR{
		Peak: &pb.Peak{
			Rr:            m.Peak.RR,
			Tier:          tierToProto(m.Peak.Tier),
			Season:        seasonToProto(m.Peak.Season),
			RankingSchema: m.Peak.RankingSchema,
		},
		Account: account.ToProto(m.Account),
		Current: &pb.Current{
			Rr:                    m.Current.RR,
			Elo:                   m.Current.Elo,
			Tier:                  tierToProto(m.Current.Tier),
			LastChange:            m.Current.LastChange,
			LeaderboardPlacement:  placementToProto(m.Current.LeaderboardPlacement),
			GamesNeededForRating:  m.Current.GamesNeededForRating,
			RankProtectionShields: m.Current.RankProtectionShields,
		},
		Seasonal: seasonal,
	}
}

func tierToProto(t tier) *pb.Tier {
	return &pb.Tier{Id: t.ID, Name: t.Name}
}

func seasonToProto(s season) *pb.Season {
	return &pb.Season{Id: s.ID, Short: s.Short}
}

func placementToProto(l leaderboardPlacement) *pb.LeaderboardPlacement {
	return &pb.LeaderboardPlacement{Rank: l.Rank, UpdatedAt: l.UpdatedAt}
}
