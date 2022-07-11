package app

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const (
	keyLenghtStart = 8
)

//search exist short url in storage,return temporary redirect if found
func (app *application) getHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	url, err := app.store.GetByShort(key)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if url.Deleted {
		render.Status(r, http.StatusGone)
		render.PlainText(w, r, "")
	}
	http.Redirect(w, r, url.Original, http.StatusTemporaryRedirect)
}

func (app *application) getUserURLs(w http.ResponseWriter, r *http.Request) {

	uid, ok := r.Context().Value(ctxValueNameUID).(string)

	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, render.M{"error": "not authorized"})
		return
	}

	urls, err := app.store.SelectByUID(uid)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	for i := 0; i < len(urls); i++ {
		urls[i].Short = fmt.Sprintf("%s/%s", app.baseURL, urls[i].Short)
	}

	if urls == nil {
		render.NoContent(w, r)
		return
	}

	render.JSON(w, r, urls)
}

func (app *application) pingDatabase(w http.ResponseWriter, r *http.Request) {
	if err := app.store.Ping(); err != nil {
		render.Status(r, http.StatusInternalServerError)
	} else {
		render.Status(r, http.StatusOK)
	}
}

//accept json, make short url, write in storage, return short url
func (app *application) postHandler(w http.ResponseWriter, r *http.Request) {

	uid, ok := r.Context().Value(ctxValueNameUID).(string)
	app.log.Debug("UID:", uid)

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
		app.log.Warn("error unmarshal json:", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{"error": "no valid json found"})
		return
	}

	parsedURL, err := url.ParseRequestURI(request.URL)
	if err != nil {
		app.log.Warn("error parse url:", request.URL, err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{"error": "no valid url found"})
	}

	app.store.Lock()
	defer app.store.Unlock()

	key := app.generateKey(keyLenghtStart)
	app.log.Warn("generated new key:", key)

	err = app.store.Save(key, parsedURL.String(), uid)
	if err != nil {
		if errors.Is(err, storage.ErrDuplicate) {
			app.log.Warn("database unique violation error", err)

			existURL, _ := app.store.GetByOriginal(parsedURL.String())
			render.Status(r, http.StatusConflict)
			result := fmt.Sprintf("%s/%s", app.baseURL, existURL.Short)
			render.JSON(w, r, render.M{"result": result})
			return
		}
		app.log.Warn(err)
	}
	app.log.Debugf("url saved: URL: '%s', key '%s'", uid, key)

	result := fmt.Sprintf("%s/%s", app.baseURL, key)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, render.M{"result": result})
}

func (app *application) batchPostHandler(w http.ResponseWriter, r *http.Request) {

	uid, ok := r.Context().Value(ctxValueNameUID).(string)
	app.log.Debug("UID:", uid)

	if !ok {
		app.log.Warn("UID not found in request")
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
		app.log.Warn("error decode json:", err)
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
			app.log.Warn("error parse url:", batchItem.OriginalURL, err)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, render.M{"error": "error parse url:" + batchItem.OriginalURL})
			return
		}
		key := app.generateKey(keyLenghtStart)
		app.log.Debug("generated key:", key)

		err = app.store.Save(key, originalURL.String(), uid)
		if err != nil {
			log.Println(err)
		}
		app.log.Debugf("url saved: URL: '%s', key '%s'", originalURL.String(), key)

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
func (app *application) legacyPostHandler(w http.ResponseWriter, r *http.Request) {

	uid, ok := r.Context().Value(ctxValueNameUID).(string)

	app.log.Debugf("UID: %s", uid)

	if !ok {
		app.log.Debug("UID not found in request")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, render.M{"error": "not authorized"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		app.log.Warnf("Error read body: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parsedURL, err := url.ParseRequestURI(string(body))
	if err != nil {
		app.log.Warnf("Error parse URL: %s; Err: %s", string(body), err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app.store.Lock()
	defer app.store.Unlock()

	key := app.generateKey(keyLenghtStart)
	app.log.Warnf("generated key: %s", key)

	err = app.store.Save(key, parsedURL.String(), uid)
	if err != nil {
		if errors.Is(err, storage.ErrDuplicate) {
			app.log.Warn("database unique violation error", err)

			existURL, _ := app.store.GetByOriginal(parsedURL.String())
			render.Status(r, http.StatusConflict)
			result := fmt.Sprintf("%s/%s", app.baseURL, existURL.Short)
			render.PlainText(w, r, result)
			return
		}
		app.log.Warn(err)
	}

	app.log.Debugf("URL saved: %s -> %s", parsedURL.String(), key)

	result := fmt.Sprintf("%s/%s", app.baseURL, key)

	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, result)
}

func (app *application) deleteHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := r.Context().Value(ctxValueNameUID).(string)

	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, render.M{"error": "not authorized"})
		return
	}
	var shorts []string
	if err := render.DecodeJSON(r.Body, &shorts); err != nil {
		app.log.Warn(err)
	}

	go app.store.Delete(uid, shorts...)

	render.Status(r, http.StatusAccepted)
	render.PlainText(w, r, "")
}
