package main

import (
	"log"

	"github.com/caarlos0/env/v6"

	"yandex-practicum-go-shortener/internal/app"
	"yandex-practicum-go-shortener/internal/storage"
)

var cfg struct {
	Addr    string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

func main() {

	var err = env.Parse(&cfg)

	if err != nil {
		log.Println(err)
	}

	app := app.New(
		app.WithStorage(storage.New()),
		app.WithBaseURL(cfg.BaseURL),
		app.WithAddr(cfg.Addr),
	)

	log.Fatal(app.Run())
}
