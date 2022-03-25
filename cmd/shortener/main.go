package main

import (
	"log"
	"yandex-practicum-go-shortener/config"

	"github.com/gin-gonic/gin"
)

func main() {
	server := CreateNewServer()
	log.Fatal(server.Run(config.Addr))
}

func CreateNewServer() *gin.Engine {
	server := gin.Default()
	return server
}
