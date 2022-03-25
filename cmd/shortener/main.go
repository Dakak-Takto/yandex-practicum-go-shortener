package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	server := CreateNewServer()
	addr := "localhost:8080"
	log.Printf("Start server on %s", addr)
	log.Fatal(server.Run(addr))
}

func CreateNewServer() *gin.Engine {
	server := gin.Default()
	// server.GET("/:key", handlers.GetHandler)
	// server.POST("/", handlers.PostHandler)
	return server

}
