package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	server := CreateNewServer()

	var shortLink string

	t.Run("Test correct query", func(t *testing.T) {
		// testing getting a short link
		reader := strings.NewReader("https://practicum.yandex.ru/learn/go-advanced/courses")
		request := httptest.NewRequest(http.MethodPost, "/", reader)
		response, body := testRequest(t, server, request)
		defer response.Body.Close()
		assert.Equal(t, response.StatusCode, http.StatusCreated)
		assert.NotEmpty(t, body)
		shortLink = body

		// test redirect
		request = httptest.NewRequest(http.MethodGet, shortLink, nil)
		response, _ = testRequest(t, server, request)
		defer response.Body.Close()
		assert.Equal(t, response.StatusCode, http.StatusTemporaryRedirect)
		assert.NotEmpty(t, response.Header.Get("Location"))
	})

	t.Run("Test incorrect query", func(t *testing.T) {
		//post without body
		request := httptest.NewRequest(http.MethodPost, "/", nil)
		response, _ := testRequest(t, server, request)
		defer response.Body.Close()
		assert.Equal(t, response.StatusCode, http.StatusBadRequest)

		//get non-existent record
		request = httptest.NewRequest(http.MethodGet, "/non-existent-link-record-09754564", nil)
		response, _ = testRequest(t, server, request)
		defer response.Body.Close()
		assert.Equal(t, response.StatusCode, http.StatusNotFound)
	})
}

func testRequest(t *testing.T, server *gin.Engine, request *http.Request) (*http.Response, string) {
	w := httptest.NewRecorder()
	server.ServeHTTP(w, request)
	body, err := ioutil.ReadAll(w.Body)
	assert.NoError(t, err)
	return w.Result(), string(body)
}
