// приложение для сокращения URL
package main

import (
	"crypto/aes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"path"
	"runtime"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/app"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/handlers"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/random"
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

	//Create storage instance
	var store storage.Storage

	switch {
	case cfg.DatabaseDSN != "":
		store, err = database.New(cfg.DatabaseDSN)
		if err != nil {
			log.Fatal(err)
		}
	case cfg.FileStoragePath != "":
		store, err = infile.New(cfg.FileStoragePath)
		if err != nil {
			log.Fatal(err)
		}
	default:
		store, err = inmem.New()
		if err != nil {
			log.Fatal(err)
		}
	}

	// create securecookie
	secret, err := random.RandomBytes(aes.BlockSize * 4)
	if err != nil {
		log.Fatal(err)
	}
	session := sessions.NewCookieStore(secret)

	//Create app instance
	app := app.New(
		app.WithStorage(store),
		app.WithBaseURL(cfg.BaseURL),
		app.WithAddr(cfg.Addr),
		app.WithLogger(log),
	)

	//create router
	router := chi.NewRouter()

	// create and register handler
	handler := handlers.New(app, cfg.BaseURL, session, log)
	handler.Register(router)

	//HTTP server
	server := http.Server{}
	server.Addr = cfg.Addr
	server.Handler = router

	go runPProfHTTPServer("localhost:8081")

	log.Fatal(server.ListenAndServe())
}

func runPProfHTTPServer(addr string) {
	log.Printf("run pprof http server on %s\n", addr)
	http.ListenAndServe(addr, nil)
}
