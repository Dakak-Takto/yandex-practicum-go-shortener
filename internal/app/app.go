package app

import (
	"compress/gzip"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"yandex-practicum-go-shortener/internal/storage"
)

type Application interface {
	Run() error
}

type application struct {
	store           storage.Storage
	baseURL         string
	addr            string
	fileStoragePath string
}

func New(opts ...Option) Application {
	app := application{}
	for _, o := range opts {
		o(&app)
	}
	return &app
}

func (app *application) Run() error {

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Compress(gzip.BestCompression, "application/*", "text/*"))
	r.Use(app.decompress)

	r.Get("/{key}", app.GetHandler)
	r.Post("/", app.LegacyPostHandler)
	r.Post("/api/shorten", app.PostHandler)

	log.Printf("Run app on %s", app.addr)
	return http.ListenAndServe(app.addr, r)
}

type Option func(app *application)

func WithStorage(storage storage.Storage) Option {
	return func(app *application) {
		app.store = storage
	}
}
func WithBaseURL(baseURL string) Option {
	return func(app *application) {
		app.baseURL = baseURL
	}
}
func WithAddr(addr string) Option {
	return func(app *application) {
		app.addr = addr
	}
}

func WithFileStoragePath(fileStoragePath string) Option {
	return func(app *application) {
		app.fileStoragePath = fileStoragePath
	}
}
