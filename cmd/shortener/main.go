package main

import (
	"log"

	"github.com/caarlos0/env/v6"

	"yandex-practicum-go-shortener/internal/app"
	"yandex-practicum-go-shortener/internal/storage"
	"yandex-practicum-go-shortener/internal/storage/infile"
	"yandex-practicum-go-shortener/internal/storage/inmem"
)

var cfg struct {
	Addr            string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"` //TODO DELETE DEFAULT!
}

func main() {

	var err = env.Parse(&cfg)

	if err != nil {
		log.Println(err)
	}

	var store storage.Storage

	if cfg.FileStoragePath != "" {
		store = infile.New(cfg.FileStoragePath)
	} else {
		store = inmem.New()
	}

	app := app.New(
		app.WithStorage(store),
		app.WithBaseURL(cfg.BaseURL),
		app.WithAddr(cfg.Addr),
	)

	log.Fatal(app.Run())
}
