package main

import (
	"crypto/aes"
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/caarlos0/env/v6"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/gorilla/securecookie"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/app"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/storage"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/storage/database"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/storage/infile"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/storage/inmem"
)

var _ json.Number        //использование известной библиотеки кодирования JSON
var _ sql.IsolationLevel //использование библиотеки database/sql

var cfg struct {
	Addr            string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN" envDefault:"-"`
}

func main() {

	processEnv()
	processFlags()

	//Create storage instance
	store, err := getStorageInstance()
	if err != nil {
		panic(err)
	}

	var secureCookie = getSecureCookieInstance()

	//Create app instance
	app := app.New(
		app.WithStorage(store),
		app.WithBaseURL(cfg.BaseURL),
		app.WithAddr(cfg.Addr),
		app.WithSecureCookie(secureCookie),
	)
	go runPProfHttpServer("localhost:8081")
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
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "dsn string for connection ")
	flag.Parse()
}

func getStorageInstance() (storage.Storage, error) {

	if cfg.DatabaseDSN != "" {
		log.Println("use database. dsn:", cfg.DatabaseDSN)
		return database.New(cfg.DatabaseDSN)
	}

	if cfg.FileStoragePath != "" {
		log.Println("use file storage. file storage path:", cfg.FileStoragePath)
		return infile.New(cfg.FileStoragePath)
	}

	log.Println("use memory storage")
	return inmem.New()
}

func getSecureCookieInstance() *securecookie.SecureCookie {
	hashKey := securecookie.GenerateRandomKey(aes.BlockSize * 2)
	blockKey := securecookie.GenerateRandomKey(aes.BlockSize * 2)
	return securecookie.New(hashKey, blockKey)
}

func runPProfHttpServer(addr string) {
	log.Printf("run pprof http server on %s\n", addr)
	http.ListenAndServe(addr, nil)
}
