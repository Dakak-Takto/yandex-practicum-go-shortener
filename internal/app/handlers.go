package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

const (
	keyLenghtStart = 5
)

func (app *application) GetHandler(c *gin.Context) {
	key := c.Param("key")
	url, err := app.store.Get(key)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (app *application) PostHandler(c *gin.Context) {
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no valid json found"})
		return
	}

	parsedURL, err := url.ParseRequestURI(request.URL)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no valid url found"})
	}

	app.store.Lock()
	defer app.store.Unlock()

	key := app.generateKey(keyLenghtStart)

	app.store.Set(key, parsedURL.String())

	result := fmt.Sprintf("%s/%s", app.baseURL, key)
	c.JSON(http.StatusCreated, gin.H{"result": result})
}

func (app *application) LegacyPostHandler(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	parsedURL, err := url.ParseRequestURI(string(body))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	app.store.Lock()
	defer app.store.Unlock()

	key := app.generateKey(keyLenghtStart)

	app.store.Set(key, parsedURL.String())
	result := fmt.Sprintf("%s/%s", app.baseURL, key)
	c.String(http.StatusCreated, result)
}
