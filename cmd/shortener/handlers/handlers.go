package handlers

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"yandex-practicum-go-shortener/cmd/shortener/storage"

	"github.com/gin-gonic/gin"
)

var addr string = "localhost:8080"
var Links = storage.CreateNew()

func GetHandler(c *gin.Context) {
	query := c.Param("short")
	link, err := Links.Get(query)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}
	c.Redirect(http.StatusTemporaryRedirect, link)
}

func PostHandler(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	long := string(body)
	if long == "" {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	md5sum := md5.Sum(body)
	short := fmt.Sprintf("%x", md5sum[:])[:6]
	Links.Set(short, long)
	log.Printf("Saved short %v = %v", short, long)
	c.Status(http.StatusCreated)
	c.String(http.StatusCreated, "http://%s/%v", addr, short)
}
