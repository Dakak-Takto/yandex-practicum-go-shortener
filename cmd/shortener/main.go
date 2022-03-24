package main

import (
	"log"
	"yandex-practicum-go-shortener/config"
	"yandex-practicum-go-shortener/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	server := CreateNewServer()
	server.SetTrustedProxies(nil)

	addr := config.GetAddr()
	log.Fatal(server.Run(addr))
}

func CreateNewServer() *gin.Engine {
	server := gin.Default()
	server.GET("/:key", handlers.GetHandler)
	server.POST("/", handlers.PostHandler)
	return server
}
