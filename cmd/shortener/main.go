package main

import (
	"flag"
	"log"
	"os"

	"yandex-practicum-go-shortener/internal/app"
	"yandex-practicum-go-shortener/internal/storage"
	"yandex-practicum-go-shortener/internal/storage/infile"
	"yandex-practicum-go-shortener/internal/storage/inmem"
)

var cfg struct {
	Addr            string
	BaseURL         string
	FileStoragePath string
}

func init() {
	flag.StringVar(&cfg.Addr, "a", os.Getenv("SERVER_ADDRESS"), "host:port")
	flag.StringVar(&cfg.BaseURL, "b", os.Getenv("BASE_URL"), "ex: http://example.com")
	flag.StringVar(&cfg.FileStoragePath, "f", os.Getenv("FILE_STORAGE_PATH"), "ex: /path/to/file")
}

func main() {

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

/*
Поддержите конфигурирование сервиса с помощью флагов командной строки наравне с уже имеющимися переменными окружения:

    флаг -a, отвечающий за адрес запуска HTTP-сервера (переменная SERVER_ADDRESS);
    флаг -b, отвечающий за базовый адрес результирующего сокращённого URL (переменная BASE_URL);
    флаг -f, отвечающий за путь до файла с сокращёнными URL (переменная FILE_STORAGE_PATH).
*/
