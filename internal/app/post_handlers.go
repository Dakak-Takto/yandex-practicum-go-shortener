package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/render"
)

const (
	keyLenghtStart = 8
)

//accept json, make short url, write in storage, return short url
func (app *application) PostHandler(w http.ResponseWriter, r *http.Request) {

	uid, ok := r.Context().Value(uidContext("uid")).(uidContext)

	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, render.M{"error": "not authorized"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var request struct {
		URL string `json:"url"`
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{"error": "no valid json found"})
		return
	}

	parsedURL, err := url.ParseRequestURI(request.URL)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{"error": "no valid url found"})
	}

	app.store.Lock()
	defer app.store.Unlock()

	key := app.generateKey(keyLenghtStart)

	app.store.Save(key, parsedURL.String(), uid.String())
	log.Printf("save short %s -> %s", key, parsedURL.String())

	result := fmt.Sprintf("%s/%s", app.baseURL, key)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, render.M{"result": result})
}

func (app *application) batchPostHandler(w http.ResponseWriter, r *http.Request) {

	uid, ok := r.Context().Value(uidContext("uid")).(uidContext)

	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, render.M{"error": "not authorized"})
		return
	}

	var batchRequestURLs []struct {
		CorellationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	render.DecodeJSON(r.Body, &batchRequestURLs)

	if batchRequestURLs == nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{"error": "bad request"})
		return
	}

	type responseURLs struct {
		CorellationID string `json:"corellation_id"`
		ShortURL      string `json:"short_url"`
	}
	var batchResponseURLs []responseURLs

	for _, batchItem := range batchRequestURLs {
		originalURL, err := url.ParseRequestURI(batchItem.OriginalURL)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, render.M{"error": "error parse url:" + batchItem.OriginalURL})
			return
		}
		key := app.generateKey(keyLenghtStart)

		log.Printf("save short %s -> %s", key, originalURL.String())
		app.store.Save(key, originalURL.String(), uid.String())

		shortURL := fmt.Sprintf("%s/%s", app.baseURL, key)

		batchResponseURLs = append(batchResponseURLs, responseURLs{
			ShortURL:      shortURL,
			CorellationID: batchItem.CorellationID,
		})
	}

	render.JSON(w, r, batchResponseURLs)
}

//accept text/plain body with url, make short url, write in storage, return short url in body
func (app *application) LegacyPostHandler(w http.ResponseWriter, r *http.Request) {

	uid, ok := r.Context().Value(uidContext("uid")).(uidContext)

	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, render.M{"error": "not authorized"})
		return
	}

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

	app.store.Save(key, parsedURL.String(), uid.String())
	result := fmt.Sprintf("%s/%s", app.baseURL, key)

	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, result)
}
