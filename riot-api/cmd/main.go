package main

import (
	"riot-api/internal/account"
	"riot-api/internal/api"
	"riot-api/internal/config"
)

func main() {
	config.LoadEnvVariables()
	conf := config.NewConfig()

	// params := map[string]string{
	// 	"name": "danilairo",
	// 	"tag":  "000",
	// }
	// _ = api.GetAccountByName(conf.ValorantAPIKey, params, nil)

	// fmt.Println(account)

	account := account.Account{
		Tag:          "000",
		Name:         "DanilAiro",
		Card:         "e6529e9c-4a2b-c31c-7252-e185a8ce4a04",
		PUUID:        "95529167-bced-547e-bcac-cc245cf01487",
		Title:        "f3009eb7-4416-39e0-4b66-edb914c7f950",
		Region:       "eu",
		UpdatedAt:    "2026-09-28T15:40:06.429Z",
		Platforms:    []string{"pc", "CONSOLE"},
		AccountLevel: 91,
	}

	params := map[string]string{
		"affinity": account.Region,
		"platform": account.Platforms[0],
		"puuid":    account.PUUID,
		"match_id": "82b03a49-ab30-42a0-a8dc-ec8e444a8d1e",
	}

	// matchesToFetch := 6

	// for i := 2; i < matchesToFetch; i++ {
	// 	if i > 0 {
	// 		fmt.Println("Ожидание перед следующим запросом...")
	// 		time.Sleep(2 * time.Second) // Пауза в 2 секунды
	// 	}

	// 	fmt.Printf("\n=== ЗАГРУЗКА МАТЧА №%d (Смещение start: %d) ===\n", i+1, i)

	query := map[string]any{
		"force": true,
		"size":  1, // Всегда запрашиваем строго по 1 матчу
		// "start": i, // Сдвигаем индекс: 0 — последний матч, 1 — предпоследний и т.д.
	}

	// 	api.GetPlayerMatches(conf.ValorantAPIKey, params, query)
	// }

	// api.GetMatchInfo(conf.ValorantAPIKey, params)

	// api.GetPlayerMMR(conf.ValorantAPIKey, params)

	api.GetAccount(conf.ValorantAPIKey, params, query)
}
