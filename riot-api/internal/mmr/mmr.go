package mmr

import (
	"fmt"
	"riot-api/internal/account"
	"strings"
)

type MMR struct {
	Peak     peak            `json:"peak"`
	Account  account.Account `json:"account"`
	Current  current         `json:"current"`
	Seasonal []seasonal      `json:"seasonal"`
}

func (m MMR) String() string {
	var b strings.Builder
	b.Grow(1024) // Увеличиваем размер буфера, так как сезонов может быть несколько

	// 1. Шапка профиля
	b.WriteString("=== VALORANT MMR PROFILE ===\n")
	fmt.Fprintf(&b, "Player:  %s#%s (Level %d)\n", m.Account.Name, m.Account.Tag, m.Account.AccountLevel)
	fmt.Fprintf(&b, "Region:  %s\n", strings.ToUpper(m.Account.Region))
	b.WriteString("----------------------------\n")

	// 2. Текущий ранг (Current)
	b.WriteString("CURRENT RANK:\n")
	fmt.Fprintf(&b, "  Rank:       %s\n", m.Current.Tier.Name)
	fmt.Fprintf(&b, "  Rating:     %d RR (Elo: %d)\n", m.Current.RR, m.Current.Elo)

	// Показываем изменение рейтинга за последний матч (с красивым знаком +/-)
	lastChangeSign := ""
	if m.Current.LastChange > 0 {
		lastChangeSign = "+"
	}
	fmt.Fprintf(&b, "  Last Match: %s%d RR\n", lastChangeSign, m.Current.LastChange)

	// Выводим позицию в Лидерборде, если игрок находится в Радиантах/Имморталах
	if m.Current.LeaderboardPlacement.Rank > 0 {
		fmt.Fprintf(&b, "  Leaderboard: #%d\n", m.Current.LeaderboardPlacement.Rank)
	}
	b.WriteString("----------------------------\n")

	// 3. Пиковый ранг (Peak)
	b.WriteString("PEAK RANK (ALL TIME):\n")
	fmt.Fprintf(&b, "  Max Rank:   %s\n", m.Peak.Tier.Name)
	fmt.Fprintf(&b, "  Max Rating: %d RR\n", m.Peak.RR)
	fmt.Fprintf(&b, "  Achieved:   %s\n", m.Peak.Season.Short)
	b.WriteString("----------------------------\n")

	// 4. Статистика сезонов (Seasonal)
	b.WriteString("SEASONAL STATS:\n")
	if len(m.Seasonal) == 0 {
		b.WriteString("  No seasonal data available\n")
	} else {
		for _, s := range m.Seasonal {
			// Защита: пропускаем пустые или неинициализированные акты
			if s.Season.Short == "" && s.Games == 0 {
				continue
			}

			// Считаем Win Rate для каждого акта отдельно
			var winRate float64
			if s.Games > 0 {
				winRate = (float64(s.Wins) / float64(s.Games)) * 100
			}

			fmt.Fprintf(&b, "  Act %-5s | Matches: %-3d | Wins: %-3d | Win Rate: %.1f%%\n",
				s.Season.Short, s.Games, s.Wins, winRate,
			)
		}
	}
	b.WriteString("==============================")

	return b.String()
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
	ID    string `json:"id"`
	Short string `json:"short"`
}
