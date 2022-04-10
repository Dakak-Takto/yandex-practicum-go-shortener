package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type jsonResponse map[string]string

const (
	keyLenghtStart = 8
)

//search exist short url in storage,return temporary redirect if found
func (app *application) GetHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	url, err := app.store.Get(key)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

//accept json, make short url, write in storage, return short url
func (app *application) PostHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var request = struct {
		URL string `json:"url"`
	}{}

	err = json.Unmarshal(body, &request)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, jsonResponse{"error": "no valid json found"})
		return
	}

	parsedURL, err := url.ParseRequestURI(request.URL)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, jsonResponse{"error": "no valid url found"})
	}

	app.store.Lock()
	defer app.store.Unlock()

	key := app.generateKey(keyLenghtStart)

	app.store.Set(key, parsedURL.String())

	result := fmt.Sprintf("%s/%s", app.baseURL, key)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, jsonResponse{"result": result})
}

//accept text/plain body with url, make short url, write in storage, return short url in body
func (app *application) LegacyPostHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parsedURL, err := url.ParseRequestURI(string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app.store.Lock()
	defer app.store.Unlock()

	key := app.generateKey(keyLenghtStart)

	app.store.Set(key, parsedURL.String())
	result := fmt.Sprintf("%s/%s", app.baseURL, key)

	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, result)
}
