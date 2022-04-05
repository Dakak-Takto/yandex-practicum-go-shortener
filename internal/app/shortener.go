package app

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

//Parse original URL from request and write short link
func (app *application) Shortener(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var request = struct {
		URL string `json:"url"`
	}{}

	err = json.Unmarshal(body, &request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	parsedURL, err := url.ParseRequestURI(request.URL)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no valid url found"})
		return
	}

	key, err := app.repository.Create(*parsedURL)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var baseURL = "http://localhost:8080" //TODO

	result := baseURL + "/" + key
	c.JSON(http.StatusCreated, gin.H{"result": result})
}
