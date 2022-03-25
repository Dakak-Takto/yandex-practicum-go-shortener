package main

import (
	"log"
	"yandex-practicum-go-shortener/config"

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
	return server
}
