package main

import (
	"log"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/app"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/storage/infile"
)

const dumpFileName = "dumpfile.json"

func main() {

	repository := infile.Load(dumpFileName)
	app := app.New(repository)

	log.Fatal(
		app.Run(),
	)
}
