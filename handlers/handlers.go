package handlers

import (
	"io"
	"net/http"
	"net/url"
	"yandex-practicum-go-shortener/config"
	"yandex-practicum-go-shortener/storage"

	"github.com/gin-gonic/gin"
)

func GetHandler(c *gin.Context) {
	key := c.Param("key")
	url, err := storage.GetValueByKey(key)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func PostHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	parsedURL, err := url.ParseRequestURI(string(body))
	if err != nil {
		c.String(http.StatusBadRequest, "specify valid url")
		return
	}
	key := storage.SetValueReturnKey(parsedURL.String())

	c.String(http.StatusCreated, "%s://%s/%s", config.GetScheme(), config.GetAddr(), key)
}
