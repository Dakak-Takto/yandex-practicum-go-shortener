package main

import (
	"log"
	"yandex-practicum-go-shortener/config"
	"yandex-practicum-go-shortener/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	server := CreateNewServer()
	log.Fatal(server.Run(config.Addr))
}

func CreateNewServer() *gin.Engine {
	server := gin.Default()
	server.GET("/:key", handlers.GetHandler)
	server.POST("/api/shorten", handlers.PostHandler)
	return server
}
