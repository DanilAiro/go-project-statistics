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
	api.GetAccountByName(conf.ValorantAPIKey, params, nil)
}
