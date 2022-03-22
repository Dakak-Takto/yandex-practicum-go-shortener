package handlers

import (
	"io"
	"net/http"
	"strings"
	"yandex-practicum-go-shortener/storage"

	"github.com/gin-gonic/gin"
)

func GetHandler(c *gin.Context) {
	query := c.Param("key")
	link, err := storage.Get(query)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, link)
}

func PostHandler(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	newURLValue := string(body)
	newURLValue = strings.TrimSpace(newURLValue)
	if newURLValue == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	key, err := storage.Save(newURLValue)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	scheme := "http"
	if tls := c.Request.TLS; tls != nil {
		scheme = "https"
	}
	c.String(http.StatusCreated, "%s://%s/%s", scheme, c.Request.Host, key)
}
