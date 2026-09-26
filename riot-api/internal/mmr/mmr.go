package mmr

import "riot-api/internal/account"

type MMR struct {
	Peak     peak            `json:"peak"`
	Account  account.Account `json:"account"`
	Current  current         `json:"current"`
	Seasonal seasonal        `json:"seasonal"`
}

type current struct {
	RR                    int32                `json:"rr"`
	Elo                   int32                `json:"elo"`
	Tier                  tier                 `json:"tier"`
	LastChange            int32                `json:"last_change"`
	LeaderboardPlacement  leaderboardPlacement `json:"leaderboard_placement"`
	GamesNeededForRating  int32                `json:"games_needed_for_rating"`
	RankProtectionShields int32                `json:"rank_protection_shields"`
}

type seasonal struct {
	Wins                 int32                `json:"wins"`
	Games                int32                `json:"games"`
	EndRR                int32                `json:"end_rr"`
	Season               season               `json:"season"`
	EndTier              tier                 `json:"end_tier"`
	ActWins              []actWins            `json:"act_wins"`
	RankingSchema        string               `json:"ranking_schema"`
	LeaderboardPlacement leaderboardPlacement `json:"leaderboard_placement"`
}

type peak struct {
	RR            int32  `json:"rr"`
	Tier          tier   `json:"tier"`
	Season        season `json:"season"`
	RankingSchema string `json:"ranking_schema"`
}

type tier struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type leaderboardPlacement struct {
	Rank      int32  `json:"rank"`
	UpdatedAt string `json:"updated_at"`
}

type actWins struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type season struct {
	ID    int32  `json:"id"`
	Short string `json:"short"`
}
