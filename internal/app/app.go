package app

import (
	"log"

	"github.com/gin-gonic/gin"

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
	gin.SetMode("test")

	server := gin.New()
	server.Use(gin.Logger())
	server.GET("/:key", app.GetHandler)
	server.POST("/", app.LegacyPostHandler)
	server.POST("/api/shorten", app.PostHandler)

	log.Printf("Run app on %s", app.addr)
	return server.Run(app.addr)
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
