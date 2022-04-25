package app

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"yandex-practicum-go-shortener/internal/storage"

	"github.com/go-chi/render"
)

const (
	keyLenghtStart = 8
)

//accept json, make short url, write in storage, return short url
func (app *application) PostHandler(w http.ResponseWriter, r *http.Request) {

	uid, ok := r.Context().Value(ctxValueNameUid).(string)
	app.logger.Print("UID:", uid)

	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, render.M{"error": "not authorized"})
		return
	}

	var request struct {
		URL string `json:"url"`
	}

	err := render.DecodeJSON(r.Body, &request)
	if err != nil {
		app.logger.Print("error unmarshal json:", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{"error": "no valid json found"})
		return
	}

	parsedURL, err := url.ParseRequestURI(request.URL)
	if err != nil {
		app.logger.Print("error parse url:", request.URL, err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{"error": "no valid url found"})
	}

	app.store.Lock()
	defer app.store.Unlock()

	key := app.generateKey(keyLenghtStart)
	app.logger.Print("generated new key:", key)

	err = app.store.Save(key, parsedURL.String(), uid)
	if err != nil {
		if errors.Is(err, storage.ErrDuplicate) {
			app.logger.Print("database unique violation error", err)

			existUrl, _ := app.store.GetByOriginal(parsedURL.String())
			render.Status(r, http.StatusConflict)
			result := fmt.Sprintf("%s/%s", app.baseURL, existUrl.Short)
			render.JSON(w, r, render.M{"result": result})
			return
		}
		app.logger.Print(err)
	}
	app.logger.Printf("url saved: URL: '%s', key '%s'", uid, key)

	result := fmt.Sprintf("%s/%s", app.baseURL, key)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, render.M{"result": result})
}

func (app *application) batchPostHandler(w http.ResponseWriter, r *http.Request) {

	uid, ok := r.Context().Value(ctxValueNameUid).(string)
	app.logger.Print("UID:", uid)

	if !ok {
		app.logger.Print("UID not found in request")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, render.M{"error": "not authorized"})
		return
	}

	var batchRequestURLs []struct {
		CorellationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	err := render.DecodeJSON(r.Body, &batchRequestURLs)
	if err != nil {
		app.logger.Print("error decode json:", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{"error": "bad request"})
	}

	type responseURLs struct {
		CorellationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}

	var batchResponseURLs []responseURLs

	for _, batchItem := range batchRequestURLs {
		originalURL, err := url.ParseRequestURI(batchItem.OriginalURL)
		if err != nil {
			app.logger.Print("error parse url:", batchItem.OriginalURL, err)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, render.M{"error": "error parse url:" + batchItem.OriginalURL})
			return
		}
		key := app.generateKey(keyLenghtStart)
		app.logger.Print("generated key:", key)

		app.store.Save(key, originalURL.String(), uid)
		app.logger.Printf("url saved: URL: '%s', key '%s'", originalURL.String(), key)

		shortURL := fmt.Sprintf("%s/%s", app.baseURL, key)

		batchResponseURLs = append(batchResponseURLs, responseURLs{
			CorellationID: batchItem.CorellationID,
			ShortURL:      shortURL,
		})
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, batchResponseURLs)
}

//accept text/plain body with url, make short url, write in storage, return short url in body
func (app *application) LegacyPostHandler(w http.ResponseWriter, r *http.Request) {

	uid, ok := r.Context().Value(ctxValueNameUid).(string)

	app.logger.Printf("UID: %s", uid)

	if !ok {
		app.logger.Printf("UID not found in request")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, render.M{"error": "not authorized"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		app.logger.Printf("Error read body: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parsedURL, err := url.ParseRequestURI(string(body))
	if err != nil {
		app.logger.Printf("Error parse URL: %s; Err: %s", string(body), err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app.store.Lock()
	defer app.store.Unlock()

	key := app.generateKey(keyLenghtStart)
	app.logger.Printf("generated key: %s", key)

	err = app.store.Save(key, parsedURL.String(), uid)
	if err != nil {
		if errors.Is(err, storage.ErrDuplicate) {
			app.logger.Print("database unique violation error", err)

			existUrl, _ := app.store.GetByOriginal(parsedURL.String())
			render.Status(r, http.StatusConflict)
			result := fmt.Sprintf("%s/%s", app.baseURL, existUrl.Short)
			render.PlainText(w, r, result)
			return
		}
		app.logger.Print(err)
	}

	app.logger.Printf("URL saved: %s -> %s", parsedURL.String(), key)

	result := fmt.Sprintf("%s/%s", app.baseURL, key)

	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, result)
}
