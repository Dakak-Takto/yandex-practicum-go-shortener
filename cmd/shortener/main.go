package main

import (
	"crypto/aes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"path"
	"runtime"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/entity"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/handlers"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/repository/infile"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/repository/inmem"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/repository/postgresql"
	_usecase "github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/usecase/url"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/pkg/random"
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
	// init logger
	log := logrus.StandardLogger()

	// setup logger
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

	// parse env
	var err = env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	// parse flags
	flag.StringVar(&cfg.Addr, "a", cfg.Addr, "host:port")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "ex: http://example.com")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "ex: /path/to/file")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "dsn string for connection ")
	flag.Parse()

	//init storage instance
	var repo entity.URLRepository

	switch {

	case cfg.DatabaseDSN != "":
		log.Debug("use database. dsn:", cfg.DatabaseDSN)
		repo, err = postgresql.New(cfg.DatabaseDSN)
		if err != nil {
			log.Fatal(err)
		}

	case cfg.FileStoragePath != "":
		log.Debug("use file storage. file storage path:", cfg.FileStoragePath)
		repo, err = infile.New(cfg.FileStoragePath)
		if err != nil {
			log.Fatal(err)
		}

	default:

		log.Debug("use memory storage")
		repo, err = inmem.New()
		if err != nil {
			log.Fatal(err)
		}
	}

	// init router
	router := chi.NewMux()

	// init cookieStore
	secret, err := random.RandomBytes(aes.BlockSize * 4)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("cookie secret: %x", secret)
	cookieStore := sessions.NewCookieStore(secret)

	usecase := _usecase.New(repo, log)
	handler := handlers.New(usecase, log, cfg.BaseURL, cookieStore)
	handler.Register(router)

	// init http server
	server := http.Server{}
	server.Addr = cfg.Addr
	server.Handler = router

	//Run app
	log.Fatal(server.ListenAndServe())
}
