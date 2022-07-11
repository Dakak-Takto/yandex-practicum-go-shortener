package app

import (
	"compress/gzip"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/securecookie"
	"github.com/sirupsen/logrus"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/random"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/storage"
)

type Application interface {
	Run() error
}

type application struct {
	store        storage.Storage
	baseURL      string
	addr         string
	secureCookie *securecookie.SecureCookie
	log          *logrus.Logger
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
	router.Use(middleware.Compress(gzip.BestCompression, "application/*", "text/*"))
	router.Use(app.decompress)
	router.Use(app.SetCookie)
	router.Use(app.httpLog)

	//Routes
	router.Route("/", func(r chi.Router) {
		router.Get("/{key}", app.getHandler)
		router.Get("/ping", app.pingDatabase)
		router.Post("/", app.legacyPostHandler)
	})

	router.Route("/api/shorten", func(r chi.Router) {
		r.Post("/", app.postHandler)
		r.Post("/batch", app.batchPostHandler)
	})

	router.Route("/api/user/urls", func(r chi.Router) {
		r.Get("/", app.getUserURLs)
		r.Delete("/", app.deleteHandler)
	})

	//Run
	app.log.Infof("Run app on %s", app.addr)

	//Http server
	server := http.Server{}
	server.Addr = app.addr
	server.Handler = router

	return server.ListenAndServe()
}

/*generating unique key in cycle. If key will be exists in storage len be increase by one for each iteration*/
func (app *application) generateKey(startLenght int) string {
	var n = startLenght

	for {
		short := random.String(n)
		if _, err := app.store.GetByShort(short); err == nil {
			n = n + 1
			continue
		} else {
			return short
		}
	}

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

func WithLogger(log *logrus.Logger) Option {
	return func(app *application) {
		app.log = log
	}
}
