package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	var req = struct {
		Url string `json:"url"`
	}{}
	if err := json.Unmarshal(body, &req); err != nil {
		httpErrorJSON(c, http.StatusBadRequest, "bad request")
		log.Println(err)
		return
	}

	parsedURL, err := url.ParseRequestURI(req.Url)
	if err != nil {
		httpErrorJSON(c, http.StatusBadRequest, "url parsing error")
		return
	}
	key := storage.SetValueReturnKey(parsedURL.String())
	result := fmt.Sprintf("%s://%s/%s", config.Scheme, config.Addr, key)
	c.JSON(http.StatusCreated, gin.H{"result": result})
}

func httpErrorJSON(c *gin.Context, statusCode int, msg string) {
	c.JSON(statusCode, gin.H{"error": msg})
}
