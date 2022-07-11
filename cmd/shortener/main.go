package main

import (
	"flag"
	"fmt"
	"net/http"
	"path"
	"runtime"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	_handler "github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/short/handlers"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/short/model"
	_repo "github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/short/repo"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/short/usecase"
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

	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&cfg.Addr, "a", cfg.Addr, "host:port")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "ex: http://example.com")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "ex: /path/to/file")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "dsn string for connection ")
	flag.Parse()

	var repo model.ShortRepository
	var err error

	switch {
	case cfg.DatabaseDSN != "":
		log.Println("use database. dsn:", cfg.DatabaseDSN)
		repo, err = _repo.NewPostgresRepository("postgres://postgres:postgres@localhost/praktikum?sslmode=disable")
	case cfg.FileStoragePath != "":
		log.Println("use file storage. file storage path:", cfg.FileStoragePath)
		repo, err = _repo.NewFileRepository(cfg.FileStoragePath)
	default:
		log.Println("use memory storage")
		repo = _repo.NewMemoryRepository()
	}

	if err != nil {
		log.Fatal(err)
	}

	usecase := usecase.New(repo, log)

	cookieStore := sessions.NewCookieStore([]byte("secret"))
	handler := _handler.New(usecase, cfg.BaseURL, cookieStore, log)

	router := chi.NewMux()
	handler.Register(router)

	log.Fatal(
		http.ListenAndServe(cfg.Addr, router),
	)
}
