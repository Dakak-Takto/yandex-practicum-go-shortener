package main

import (
	"log"
	"yandex-practicum-go-shortener/config"
	"yandex-practicum-go-shortener/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	server := CreateNewServer()
	addr := config.GetAddr()
	log.Printf("Start server on %s", addr)
	log.Fatal(server.Run(addr))
}

func CreateNewServer() *gin.Engine {
	server := gin.Default()
	server.GET("/:key", handlers.GetHandler)
	server.POST("/", handlers.PostHandler)
	return server
}
