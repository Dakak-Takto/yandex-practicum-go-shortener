package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

var (
	Addr, BaseURL string
)

func init() {
	var config struct {
		Addr    string `env:"ADDRESS"`
		BaseURL string `env:"BASE_URL"`
	}
	err := env.Parse(&config)
	if err != nil {
		log.Fatal(err)
	}
	Addr = config.Addr
	BaseURL = config.BaseURL
	log.Printf("config loaded. server_addr: `%s`; baseURL: `%s`", Addr, BaseURL)
}
