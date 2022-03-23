package main

import (
	"yandex-practicum-go-shortener/config"
	"yandex-practicum-go-shortener/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	server := CreateNewServer()
	server.SetTrustedProxies(nil)
	server.Run(config.GetAddr())
}

func CreateNewServer() *gin.Engine {
	server := gin.Default()
	server.GET("/:key", handlers.GetHandler)
	server.POST("/", handlers.PostHandler)
	return server
}
