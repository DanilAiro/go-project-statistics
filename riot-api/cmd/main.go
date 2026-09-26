package main

import (
	"riot-api/internal/api"
	"riot-api/internal/config"
)

func main() {
	config.LoadEnvVariables()
	conf := config.NewConfig()

	params := map[string]string{
		"name": "danilairo",
		"tag":  "000",
	}
	account := api.GetAccountByName(conf.ValorantAPIKey, params, nil)

	params = map[string]string{
		"affinity": account.Region,
		"platform": account.Platforms[0],
		"puuid":    account.PUUID,
	}
	query := map[string]any{
		"size": 1,
	}
	api.GetPlayerMatches(conf.ValorantAPIKey, params, query)
}
