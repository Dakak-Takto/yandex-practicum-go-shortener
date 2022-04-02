package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

var testURL = "https://practicum.yandex.ru/learn/go-advanced/courses"
var server = CreateNewServer()

func TestHandlers(t *testing.T) {

	var shortLink string

	t.Run("Getting short link", func(t *testing.T) {
		reqBody := `{"url": "https://practicum.yandex.ru/learn/go-advanced/courses"}`
		reader := strings.NewReader(reqBody)
		request := httptest.NewRequest(http.MethodPost, "/api/shorten", reader)
		response, body := testRequest(t, server, request)
		defer response.Body.Close()
		require.Equal(t, response.StatusCode, http.StatusCreated)
		require.NotEmpty(t, body)

		var respJSON struct{ Result string }
		err := json.Unmarshal([]byte(body), &respJSON)
		require.NoError(t, err)
		require.NotEmpty(t, respJSON)
		shortLink = respJSON.Result
	})

	t.Run("Getting redirect", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, shortLink, nil)
		response, _ := testRequest(t, server, request)
		defer response.Body.Close()
		require.Equal(t, response.StatusCode, http.StatusTemporaryRedirect)
		require.Equal(t, response.Header.Get("Location"), testURL)
	})

}

func testRequest(t *testing.T, server *gin.Engine, request *http.Request) (*http.Response, string) {
	w := httptest.NewRecorder()
	server.ServeHTTP(w, request)
	body, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	return w.Result(), string(body)
}
