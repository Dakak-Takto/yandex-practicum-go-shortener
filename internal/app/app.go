package app

import (
	"compress/gzip"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/gorilla/securecookie"
	"github.com/rs/zerolog"

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
	logger       zerolog.Logger
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
	app.logger = httplog.NewLogger("httplog", httplog.Options{LogLevel: "debug", JSON: false})
	router.Use(httplog.Handler(app.logger))

	router.Use(middleware.Compress(gzip.BestCompression, "application/*", "text/*"))
	router.Use(app.decompress)
	router.Use(app.SetCookie)

	//Routes
	router.Get("/{key}", app.getHandler)
	router.Get("/ping", app.pingDatabase)
	router.Post("/", app.legacyPostHandler)
	router.Post("/api/shorten", app.postHandler)
	router.Get("/api/user/urls", app.getUserURLs)
	router.Get("/api/user/urls", app.deleteHandler)
	router.Post("/api/shorten/batch", app.batchPostHandler)

	//Run
	app.logger.Printf("Run app on %s", app.addr)

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
