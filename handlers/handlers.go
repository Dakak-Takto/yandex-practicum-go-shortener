package handlers

import (
	"encoding/json"
	"fmt"
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

	var request = struct {
		Url string
	}{}

	err = json.Unmarshal(body, &request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	parsedURL, err := url.ParseRequestURI(request.Url)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no valid url found"})
		return
	}
	key := storage.SetValueReturnKey(parsedURL.String())

	result := fmt.Sprintf("%s://%s/%s", config.Scheme, config.Addr, key)

	c.JSON(http.StatusCreated, gin.H{"result": result})
}
