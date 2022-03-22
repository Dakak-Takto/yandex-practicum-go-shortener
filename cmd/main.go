package main

import (
	"yandex-practicum-go-shortener/handlers"

	"github.com/gin-gonic/gin"
)

var addr string = "localhost:8080"

func main() {
	server := CreateNewServer()
	server.SetTrustedProxies(nil)
	server.Run(addr)
}

func CreateNewServer() *gin.Engine {
	server := gin.New()
	server.GET("/:key", handlers.GetHandler)
	server.POST("/", handlers.PostHandler)
	return server
}
