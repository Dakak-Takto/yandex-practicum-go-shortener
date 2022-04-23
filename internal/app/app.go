package app

import (
	"compress/gzip"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/securecookie"

	"yandex-practicum-go-shortener/internal/storage"
)

type Application interface {
	Run() error
}

type application struct {
	store        storage.Storage
	baseURL      string
	addr         string
	secureCookie *securecookie.SecureCookie
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
	router.Use(app.debug)

	router.Use(middleware.Logger)
	router.Use(middleware.Compress(gzip.BestCompression, "application/*", "text/*"))
	router.Use(app.decompress)
	router.Use(app.SetCookie)

	//Routes
	router.Get("/{key}", app.GetHandler)
	router.Get("/ping", app.pingDatabase)
	router.Post("/", app.LegacyPostHandler)
	router.Post("/api/shorten", app.PostHandler)
	router.Get("/api/user/urls", app.getUserURLs)

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

func WithSecureCookie(s *securecookie.SecureCookie) Option {
	return func(app *application) {
		app.secureCookie = s
	}
}
