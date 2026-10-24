package match

import (
	"context"
	"encoding/json"

	"riot-api/internal/api"

	pb "api-contracts/gen/go/api/proto/match"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

type server struct {
	apiKey string
	pb.UnimplementedMatchServiceServer
}

func NewMatchServer(apiKey string) *server {
	return &server{
		apiKey: apiKey,
	}
}

func (s *server) GetPlayerMatches(
	ctx context.Context,
	req *pb.GetPlayerMatchesRequest,
) (*pb.GetPlayerMatchesResponse, error) {
	puuid := req.GetPuuid()
	affinity := req.GetAffinity()
	platform := req.GetPlatform()

	if puuid == "" {
		return nil, status.Error(codes.InvalidArgument, "PUUID пустой")
	}

	if affinity == "" {
		return nil, status.Error(codes.InvalidArgument, "Affinity пустой")
	}

	if platform == "" {
		return nil, status.Error(codes.InvalidArgument, "Platform пустой")
	}

	params := map[string]string{
		"puuid":    puuid,
		"affinity": affinity,
		"platform": platform,
	}

	query := make(map[string]any)

	if req.Mode != nil {
		query["mode"] = req.GetMode()
	}

	if req.Map != nil {
		query["map"] = req.GetMap()
	}

	if req.Size != nil {
		query["size"] = req.GetSize()
	}

	if req.Start != nil {
		query["start"] = req.GetStart()
	}

	data, err := api.GetPlayerMatches(s.apiKey, params, query)
	if err != nil {
		return nil, err
	}

	var matches []Match
	if err := json.Unmarshal(data, &matches); err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"не удалось декодировать список матчей: %v",
			err,
		)
	}

	pbMatches := make([]*pb.Match, 0, len(matches))
	for i := range matches {
		pbMatch, err := matchToProto(&matches[i])
		if err != nil {
			return nil, err
		}

		pbMatches = append(pbMatches, pbMatch)
	}

	return &pb.GetPlayerMatchesResponse{
		Matches: pbMatches,
	}, nil
}

func (s *server) GetMatchInfo(
	ctx context.Context,
	req *pb.GetMatchInfoRequest,
) (*pb.GetMatchInfoResponse, error) {
	affinity := req.GetAffinity()
	matchID := req.GetMatchId()

	if affinity == "" {
		return nil, status.Error(codes.InvalidArgument, "Affinity пустой")
	}

	if matchID == "" {
		return nil, status.Error(codes.InvalidArgument, "Match ID пустой")
	}

	params := map[string]string{
		"affinity": affinity,
		"match_id": matchID,
	}

	data, err := api.GetMatchInfo(s.apiKey, params)
	if err != nil {
		return nil, err
	}

	var m Match
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"не удалось декодировать матч: %v",
			err,
		)
	}

	pbMatch, err := matchToProto(&m)
	if err != nil {
		return nil, err
	}

	return &pb.GetMatchInfoResponse{
		Match: pbMatch,
	}, nil
}

func matchToProto(m *Match) (*pb.Match, error) {
	if m == nil {
		return nil, status.Error(codes.Internal, "матч пустой")
	}

	pbKills := make([]*pb.Kill, 0, len(m.Kills))
	for i := range m.Kills {
		pbKills = append(pbKills, killToProto(&m.Kills[i]))
	}

	pbTeams := make([]*pb.Team, 0, len(m.Teams))
	for i := range m.Teams {
		pbTeams = append(pbTeams, teamToProto(&m.Teams[i]))
	}

	pbRounds := make([]*pb.Round, 0, len(m.Rounds))
	for i := range m.Rounds {
		pbRounds = append(pbRounds, roundToProto(&m.Rounds[i]))
	}

	pbPlayers := make([]*pb.Player, 0, len(m.Players))
	for i := range m.Players {
		pbPlayers = append(pbPlayers, playerToProto(&m.Players[i]))
	}

	pbCoaches := make([]*pb.Coach, 0, len(m.Coaches))
	for i := range m.Coaches {
		pbCoaches = append(pbCoaches, coachToProto(&m.Coaches[i]))
	}

	pbObservers := make([]*pb.Observer, 0, len(m.Observers))
	for i := range m.Observers {
		pbObservers = append(pbObservers, observerToProto(&m.Observers[i]))
	}

	pbMetadata, err := metadataToProto(&m.Metadata)
	if err != nil {
		return nil, err
	}

	return &pb.Match{
		Kills:     pbKills,
		Teams:     pbTeams,
		Rounds:    pbRounds,
		Players:   pbPlayers,
		Coaches:   pbCoaches,
		Metadata:  pbMetadata,
		Observers: pbObservers,
	}, nil
}

func metadataToProto(m *metadata) (*pb.Metadata, error) {
	premier, err := anyToProto(m.Premier)
	if err != nil {
		return nil, err
	}

	pbPenalties := make([]*pb.PartyRRPenalty, 0, len(m.PartyRRPenaltys))
	for _, p := range m.PartyRRPenaltys {
		pbPenalties = append(pbPenalties, &pb.PartyRRPenalty{
			PartyId: p.PartyID,
			Penalty: p.Penalty,
		})
	}

	return &pb.Metadata{
		Map: &pb.GameMap{
			Id:   m.Map.ID,
			Name: m.Map.Name,
		},
		Queue: &pb.Queue{
			Id:       m.Queue.ID,
			Name:     m.Queue.Name,
			ModeType: m.Queue.ModeType,
		},
		Region:           m.Region,
		MatchId:          m.MatchID,
		Cluster:          m.Cluster,
		Premier:          premier,
		Platform:         m.Platform,
		StartedAt:        m.StartedAt,
		GameVersion:      m.GameVersion,
		IsCompleted:      m.IsCompleted,
		GameLengthInMs:   m.GameLengthInMS,
		PartyRrPenalties: pbPenalties,
	}, nil
}

func anyToProto(v any) (*anypb.Any, error) {
	if v == nil {
		return nil, nil
	}

	raw, err := json.Marshal(v)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"не удалось сериализовать premier в JSON: %v",
			err,
		)
	}

	var value map[string]any
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"premier должен быть JSON-объектом: %v",
			err,
		)
	}

	structValue, err := structpb.NewStruct(value)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"не удалось преобразовать premier в protobuf Struct: %v",
			err,
		)
	}

	return anypb.New(structValue)
}

func playerToProto(p *player) *pb.Player {
	return &pb.Player{
		Tag:  p.Tag,
		Name: p.Name,
		Tier: &pb.Tier{
			Id:   p.Tier.ID,
			Name: p.Tier.Name,
		},
		Agent: &pb.Agent{
			Id:   p.Agent.ID,
			Name: p.Agent.Name,
		},
		Puuid:   p.PUUID,
		Stats:   statsToProto(&p.Stats),
		TeamId:  p.TeamID,
		Economy: economyToProto(&p.Economy),
		PartyId: p.PartyID,
		Behavior: &pb.Behavior{
			AfkRounds: p.Behavior.AfkRounds,
			FriendlyFire: &pb.FriendlyFire{
				Incoming: p.Behavior.FriendlyFire.Incoming,
				Outgoing: p.Behavior.FriendlyFire.Outgoing,
			},
			RoundsInSpawn: p.Behavior.RoundsInSpawn,
		},
		Platform: p.Platform,
		AbilityCasts: &pb.AbilityCasts{
			Ability1: p.AbilityCasts.Ability1,
			Ability2: p.AbilityCasts.Ability2,
			Grenade:  p.AbilityCasts.Grenade,
			Ultimate: p.AbilityCasts.Ultimate,
		},
		AccountLevel: p.AccountLevel,
		Customization: &pb.Customization{
			Card:                 p.Customization.Card,
			Title:                p.Customization.Title,
			PreferredLevelBorder: p.Customization.PreferredLevelBorder,
		},
		SessionPlaytimeInMs: p.SessionPlaytimeInMS,
	}
}

func statsToProto(s *stats) *pb.Stats {
	return &pb.Stats{
		Kills: s.Kills,
		Score: s.Score,
		Damage: &pb.Damage{
			Dealt:    s.Damage.Dealt,
			Received: s.Damage.Received,
		},
		Assists:   s.Assists,
		Legshots:  s.Legshots,
		Bodyshots: s.Bodyshots,
		Headshots: s.Headshots,
	}
}

func economyToProto(e *economy) *pb.Economy {
	return &pb.Economy{
		Spent: &pb.Spent{
			Overall: e.Spent.Overall,
			Average: e.Spent.Average,
		},
		LoadoutValue: &pb.LoadoutValue{
			Overall: e.LoadoutValue.Overall,
			Average: e.LoadoutValue.Average,
		},
	}
}

func killToProto(k *kill) *pb.Kill {
	pbAssistants := make([]*pb.Assistant, 0, len(k.Assistants))
	for _, a := range k.Assistants {
		pbAssistants = append(pbAssistants, &pb.Assistant{
			Tag:   a.Tag,
			Name:  a.Name,
			Team:  a.Team,
			Puuid: a.PUUID,
		})
	}

	pbLocations := make([]*pb.PlayerLocation, 0, len(k.PlayerLocations))
	for i := range k.PlayerLocations {
		pbLocations = append(pbLocations, playerLocationToProto(&k.PlayerLocations[i]))
	}

	killer := playerShort(k.Killer)
	victim := playerShort(k.Victim)

	return &pb.Kill{
		Round:  k.Round,
		Killer: playerShortToProto(&killer),
		Victim: playerShortToProto(&victim),
		Weapon: &pb.Weapon{
			Id:   k.Weapon.ID,
			Name: k.Weapon.Name,
			Type: k.Weapon.Type,
		},
		Location: &pb.Location{
			X: k.Location.X,
			Y: k.Location.Y,
		},
		Assistants:        pbAssistants,
		PlayerLocations:   pbLocations,
		TimeInRoundInMs:   k.TimeInRoundInMS,
		TimeInMatchInMs:   k.TimeInMatchInMS,
		SecondaryFireMode: k.SecondaryFireMode,
	}
}

func roundToProto(r *round) *pb.Round {
	pbStats := make([]*pb.Stats, 0, len(r.Stats))
	for i := range r.Stats {
		pbStats = append(pbStats, statsToProto(&r.Stats[i]))
	}

	return &pb.Round{
		Id:          r.ID,
		Stats:       pbStats,
		Plant:       plantToProto(&r.Plant),
		Result:      r.Result,
		Defuse:      defuseToProto(&r.Defuse),
		Ceremony:    r.Ceremony,
		WinningTeam: r.WinningTeam,
	}
}

func teamToProto(t *team) *pb.Team {
	pbMembers := make([]string, 0, len(t.PremierRoster.Members))
	pbMembers = append(pbMembers, t.PremierRoster.Members...)

	return &pb.Team{
		Won: t.Won,
		Rounds: &pb.RoundShort{
			Won:  t.Rounds.Won,
			Lost: t.Rounds.Lost,
		},
		TeamId: t.TeamID,
		PremierRoster: &pb.PremierRoster{
			Id:      t.PremierRoster.ID,
			Tag:     t.PremierRoster.Tag,
			Name:    t.PremierRoster.Name,
			Members: pbMembers,
			Customization: &pb.PremierCustomization{
				Icon:           t.PremierRoster.Customization.Icon,
				Image:          t.PremierRoster.Customization.Image,
				PrimaryColor:   t.PremierRoster.Customization.PrimaryColor,
				SecondaryColor: t.PremierRoster.Customization.SecondaryColor,
				TertiaryColor:  t.PremierRoster.Customization.TetritaryColor,
			},
		},
	}
}

func coachToProto(c *coach) *pb.Coach {
	return &pb.Coach{
		Puuid:  c.PUUID,
		TeamId: c.TeamID,
	}
}

func observerToProto(o *observer) *pb.Observer {
	return &pb.Observer{
		Tag:                 o.Tag,
		Name:                o.Name,
		CardId:              o.CardID,
		Puuid:               o.PUUID,
		TitleId:             o.TitleID,
		PartyId:             o.PartyID,
		AccountLevel:        o.AccountLevel,
		SessionPlaytimeInMs: o.SessionPlaytimeInMS,
	}
}

func playerShortToProto(p *playerShort) *pb.PlayerShort {
	return &pb.PlayerShort{
		Tag:   p.Tag,
		Name:  p.Name,
		Team:  p.Team,
		Puuid: p.PUUID,
	}
}

func playerLocationToProto(pl *playerLocation) *pb.PlayerLocation {
	return &pb.PlayerLocation{
		Player: playerShortToProto(&pl.Player),
		Location: &pb.Location{
			X: pl.Location.X,
			Y: pl.Location.Y,
		},
		ViewRadians: pl.ViewRadians,
	}
}

func plantToProto(p *plant) *pb.Plant {
	pbLocations := make([]*pb.PlayerLocation, 0, len(p.PlayerLocations))
	for i := range p.PlayerLocations {
		pbLocations = append(pbLocations, playerLocationToProto(&p.PlayerLocations[i]))
	}

	return &pb.Plant{
		Player: playerShortToProto(&p.Player),
		Location: &pb.Location{
			X: p.Location.X,
			Y: p.Location.Y,
		},
		RoundTimeInMs:   p.RoundTimeInMS,
		PlayerLocations: pbLocations,
	}
}

func defuseToProto(d *defuse) *pb.Defuse {
	pbLocations := make([]*pb.PlayerLocation, 0, len(d.PlayerLocations))
	for i := range d.PlayerLocations {
		pbLocations = append(pbLocations, playerLocationToProto(&d.PlayerLocations[i]))
	}

	return &pb.Defuse{
		Player: playerShortToProto(&d.Player),
		Location: &pb.Location{
			X: d.Location.X,
			Y: d.Location.Y,
		},
		RoundTimeInMs:   d.RoundTimeInMS,
		PlayerLocations: pbLocations,
	}
}
