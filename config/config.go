package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

var (
	Addr    = "localhost:8080"
	BaseURL = "http://localhost:8080"
)

func init() {
	var config struct {
		Addr    string `env:"SERVER_ADDRESS"`
		BaseURL string `env:"BASE_URL"`
	}
	err := env.Parse(&config)
	if err != nil {
		log.Fatal(err)
	}
	if config.Addr != "" {
		Addr = config.Addr
		log.Print("ADDRESS not present. Use default")
	}
	if config.BaseURL != "" {
		BaseURL = config.BaseURL
		log.Print("BASE_URL not present. Use default")
	}
}
