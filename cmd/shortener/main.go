package main

import (
	"flag"
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

	flag.StringVar(&cfg.Addr, "a", cfg.Addr, "host:port")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "ex: http://example.com")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "ex: /path/to/file")
	flag.Parse()

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
