package app

import (
	"github.com/gin-gonic/gin"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/storage"
)

type Application interface {
	Shortener(c *gin.Context)
	Redirector(c *gin.Context)
	Run() error
}

type application struct {
	repository storage.Repository
}

var _ Application = (*application)(nil)

//Create and return application instance
func New(r storage.Repository) Application {
	var app = application{
		repository: r,
	}

	return &app
}

//Setup and run webserver
func (a *application) Run() error {
	var router = gin.Default()
	router.POST("/api/shorten", a.Shortener)
	router.GET("/:key", a.Redirector)
	return router.Run("localhost:8080")
}
