package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"yandex-practicum-go-shortener/cmd/shortener/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

var testURL = "https://practicum.yandex.ru/learn/go-advanced/courses"
var server = CreateNewServer()

func TestHandlers(t *testing.T) {

	var shortLink string

	t.Run("Getting short link", func(t *testing.T) {
		reader := strings.NewReader(testURL)
		request := httptest.NewRequest(http.MethodPost, "/", reader)
		response, body := testRequest(t, server, request)
		defer response.Body.Close()
		require.Equal(t, response.StatusCode, http.StatusCreated)
		require.NotEmpty(t, body)
		shortLink = body
	})

	t.Run("Many requests", func(t *testing.T) {
		startTime := time.Now()
		gin.SetMode(gin.ReleaseMode)

		c := time.After(time.Second)

		rCount := testStress(t, c)

		t.Logf("Complete %d requests for %f second\n", rCount, time.Since(startTime).Seconds())
		t.Logf("Len storage: %d", storage.Len())
	})

	t.Run("Getting redirect", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, shortLink, nil)
		response, _ := testRequest(t, server, request)
		defer response.Body.Close()
		require.Equal(t, response.StatusCode, http.StatusTemporaryRedirect)
		require.Equal(t, response.Header.Get("Location"), testURL)
	})

	t.Run("Send method POST without body", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/", nil)
		response, _ := testRequest(t, server, request)
		defer response.Body.Close()
		require.Equal(t, response.StatusCode, http.StatusBadRequest)
	})

	t.Run("Get non-existent record", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/001122334455", nil)
		response, _ := testRequest(t, server, request)
		defer response.Body.Close()
		require.Equal(t, response.StatusCode, http.StatusNotFound)
	})
}

func testStress(t *testing.T, c <-chan time.Time) int {
	i := 1
	for {
		select {
		case <-c:
			return i
		default:
			reader := strings.NewReader(testURL)
			request := httptest.NewRequest(http.MethodPost, "/", reader)
			response, body := testRequest(t, server, request)
			defer response.Body.Close()
			require.Equal(t, response.StatusCode, http.StatusCreated)
			require.NotEmpty(t, body)
			i++
		}
	}
}

func testRequest(t *testing.T, server *gin.Engine, request *http.Request) (*http.Response, string) {
	w := httptest.NewRecorder()
	server.ServeHTTP(w, request)
	body, err := ioutil.ReadAll(w.Body)
	require.NoError(t, err)
	return w.Result(), string(body)
}
