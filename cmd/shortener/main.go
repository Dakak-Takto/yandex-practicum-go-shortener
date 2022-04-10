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
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func main() {

	processEnv()
	processFlags()

	//Create storage instance
	var store storage.Storage
	var err error

	if cfg.FileStoragePath != "" {
		store, err = infile.New(cfg.FileStoragePath)
	} else {
		store, err = inmem.New()
	}
	if err != nil {
		log.Fatal(err)
	}

	//Create app instance
	app := app.New(
		app.WithStorage(store),
		app.WithBaseURL(cfg.BaseURL),
		app.WithAddr(cfg.Addr),
	)

	//Run app
	log.Fatal(app.Run())
}

func processEnv() {
	var err = env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func processFlags() {
	flag.StringVar(&cfg.Addr, "a", cfg.Addr, "host:port")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "ex: http://example.com")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "ex: /path/to/file")
	flag.Parse()
}
