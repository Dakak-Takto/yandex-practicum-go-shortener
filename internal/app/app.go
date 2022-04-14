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
	store   storage.Storage
	baseURL string
	addr    string
}

func New(opts ...Option) Application {
	app := application{}
	for _, o := range opts {
		o(&app)
	}
	return &app
}

func (app *application) Run() error {

	router := chi.NewRouter()

	//Middlewares
	router.Use(middleware.Logger)
	router.Use(middleware.Compress(gzip.BestCompression, "application/*", "text/*"))
	router.Use(app.decompress)

	//Routes
	router.Get("/{key}", app.GetHandler)
	router.Post("/", app.LegacyPostHandler)
	router.Post("/api/shorten", app.PostHandler)

	//Run
	log.Printf("Run app on %s", app.addr)

	//Http server
	server := http.Server{}
	server.Addr = app.addr
	server.Handler = router

	return server.ListenAndServe()
}

//Application option declaration

type Option func(app *application)

//add storage to application
func WithStorage(storage storage.Storage) Option {
	return func(app *application) {
		app.store = storage
	}
}

//change application base_url
func WithBaseURL(baseURL string) Option {
	return func(app *application) {
		app.baseURL = baseURL
	}
}

//change http server addr
func WithAddr(addr string) Option {
	return func(app *application) {
		app.addr = addr
	}
}
