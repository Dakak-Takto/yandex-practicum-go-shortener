package main

import (
	"crypto/aes"
	"encoding/hex"
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/gorilla/securecookie"

	"yandex-practicum-go-shortener/internal/app"
	"yandex-practicum-go-shortener/internal/storage"
	"yandex-practicum-go-shortener/internal/storage/infile"
	"yandex-practicum-go-shortener/internal/storage/inmem"
)

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
		// return storage.Storage{}
	}

	if cfg.FileStoragePath != "" {
		log.Println("use file storage. file storage path:", cfg.FileStoragePath)
		return infile.New(cfg.FileStoragePath)
	}

	log.Println("use memory storage")
	return inmem.New()
}

func getSecureCookieInstance() *securecookie.SecureCookie {
	store, err := infile.New("secure_cookie.txt")
	if err != nil {
		panic(err)
	}

	var hashKey, blockKey []byte

	result, err := store.First("cookieHashKey")
	if err != nil {
		hashKey = securecookie.GenerateRandomKey(aes.BlockSize * 2)
		store.Insert("cookieHashKey", hex.EncodeToString(hashKey))
	} else {
		hashKey, err = hex.DecodeString(result.Value)
		if err != nil {
			panic(err)
		}
	}

	result, err = store.First("cookieBlockKey")
	if err != nil {
		blockKey = securecookie.GenerateRandomKey(aes.BlockSize * 2)
		store.Insert("cookieBlockKey", hex.EncodeToString(blockKey))
	} else {
		blockKey, err = hex.DecodeString(result.Value)
		if err != nil {
			panic(err)
		}
	}
	return securecookie.New(hashKey, blockKey)
}
