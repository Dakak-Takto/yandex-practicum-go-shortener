package main

import (
	"crypto/aes"
	"database/sql"
	_ "database/sql"
	"encoding/json"
	_ "encoding/json"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"path"
	"runtime"

	"github.com/caarlos0/env/v6"
	"github.com/gorilla/securecookie"
	"github.com/sirupsen/logrus"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/app"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/random"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/storage"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/storage/database"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/storage/infile"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/storage/inmem"
)

var (
	_ json.Decoder
	_ sql.DB
)

var cfg struct {
	Addr            string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN" envDefault:"-"`
}

func main() {

	log := logrus.StandardLogger()
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "15:05:05",
		FullTimestamp:   true,
		ForceColors:     true,
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {

			return "", fmt.Sprintf(" %s:%d", path.Base(f.File), f.Line)
		},
	})
	log.SetReportCaller(true)
	log.SetLevel(logrus.DebugLevel)
	log.Debug("init logger")

	var err = env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&cfg.Addr, "a", cfg.Addr, "host:port")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "ex: http://example.com")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "ex: /path/to/file")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "dsn string for connection ")
	flag.Parse()

	//Create storage instance
	var store storage.Storage

	switch {
	case cfg.DatabaseDSN != "":
		log.Debug("use database. dsn:", cfg.DatabaseDSN)
		store, err = database.New(cfg.DatabaseDSN)
		if err != nil {
			log.Fatal(err)
		}

	case cfg.FileStoragePath != "":
		log.Debug("use file storage. file storage path:", cfg.FileStoragePath)
		store, err = infile.New(cfg.FileStoragePath)
		if err != nil {
			log.Fatal(err)
		}
	default:

		log.Debug("use memory storage")
		store, err = inmem.New()
		if err != nil {
			log.Fatal(err)
		}
	}

	secret, err := random.RandomBytes(aes.BlockSize * 4)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("cookie secret: %x", secret)
	secureCookie := securecookie.New(secret[:32], secret[32:])

	//Create app instance
	app := app.New(
		app.WithStorage(store),
		app.WithBaseURL(cfg.BaseURL),
		app.WithAddr(cfg.Addr),
		app.WithSecureCookie(secureCookie),
		app.WithLogger(log),
	)

	//pprof
	go func() {
		http.ListenAndServe("localhost:8008", nil)
	}()

	//Run app
	log.Fatal(app.Run())
}
