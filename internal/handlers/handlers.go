// Package handlers contain http handlers
package handlers

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	_url "net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/app"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/storage"
)

type handler struct {
	app      app.Application
	baseURL  string
	sessions *sessions.CookieStore
	log      *logrus.Logger
}

func New(app app.Application, baseURL string, sessions *sessions.CookieStore, log *logrus.Logger) *handler {
	return &handler{
		app:      app,
		baseURL:  baseURL,
		sessions: sessions,
		log:      log,
	}
}

func (h *handler) Register(router *chi.Mux) {

	router.Use(middleware.Compress(gzip.BestCompression, "application/*", "text/*"))
	router.Use(h.decompress)
	router.Use(h.SetCookie)

	router.Route("/", func(r chi.Router) {
		router.Get("/{key}", h.getHandler)
		router.Get("/ping", h.pingDatabase)
		router.Post("/", h.legacyPostHandler)
	})

	router.Route("/api/shorten", func(r chi.Router) {
		r.Post("/", h.postHandler)
		r.Post("/batch", h.batchPostHandler)
	})

	router.Route("/api/user/urls", func(r chi.Router) {
		r.Get("/", h.getUserURLs)
		r.Delete("/", h.deleteHandler)
	})

}

//accept json, make short url, write in storage, return short url
func (h *handler) postHandler(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(contextUserID).(string)
	h.log.Debug("UID:", userID)

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
		h.log.Warn("error unmarshal json:", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{"error": "no valid json found"})
		return
	}

	url, err := h.app.MakeShort(request.URL, userID)
	if err != nil {
		if errors.Is(err, storage.ErrDuplicate) {
			h.log.Print("database unique violation error", err)

			existURL, _ := h.app.GetByOriginal(request.URL)
			render.Status(r, http.StatusConflict)
			result := fmt.Sprintf("%s/%s", h.baseURL, existURL.Short)
			render.JSON(w, r, render.M{"result": result})
			return
		}
		h.log.Warn(err)
	}
	h.log.Debugf("url saved: URL: '%s', key '%s'", userID, url.Short)

	result := fmt.Sprintf("%s/%s", h.baseURL, url.Short)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, render.M{"result": result})
}

func (h *handler) batchPostHandler(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(contextUserID).(string)
	h.log.Print("UID:", userID)

	if !ok {
		h.log.Warn("UID not found in request")
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
		h.log.Warn("error decode json:", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{"error": "bad request"})
	}

	type responseURLs struct {
		CorellationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}

	var batchResponseURLs []responseURLs

	for _, batchItem := range batchRequestURLs {
		originalURL, err := _url.ParseRequestURI(batchItem.OriginalURL)
		if err != nil {
			h.log.Warn("error parse url:", batchItem.OriginalURL, err)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, render.M{"error": "error parse url:" + batchItem.OriginalURL})
			return
		}

		url, err := h.app.MakeShort(originalURL.String(), userID)
		if err != nil {
			http.Error(w, "error make short", http.StatusBadRequest)
			return
		}
		h.log.Debugf("url saved: URL: '%s', key '%s'", url.Original, url.Short)

		shortURL := fmt.Sprintf("%s/%s", h.baseURL, url.Short)

		batchResponseURLs = append(batchResponseURLs, responseURLs{
			CorellationID: batchItem.CorellationID,
			ShortURL:      shortURL,
		})
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, batchResponseURLs)
}

//accept text/plain body with url, make short url, write in storage, return short url in body
func (h *handler) legacyPostHandler(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(contextUserID).(string)

	h.log.Printf("UserID: %s", userID)

	if !ok {
		h.log.Warn("UID not found in request")
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, render.M{"error": "not authorized"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Warn("Error read body: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url, err := h.app.MakeShort(string(body), userID)
	if err != nil {
		if errors.Is(err, storage.ErrDuplicate) {
			h.log.Warn("database unique violation error", err)

			existURL, _ := h.app.GetByOriginal(string(body))
			render.Status(r, http.StatusConflict)
			result := fmt.Sprintf("%s/%s", h.baseURL, existURL.Short)
			render.PlainText(w, r, result)
			return
		}
		h.log.Warn(err)
	}

	h.log.Debug("URL saved: %s -> %s", url.Original, url.Short)

	result := fmt.Sprintf("%s/%s", h.baseURL, url.Short)

	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, result)
}

//search exist short url in storage,return temporary redirect if found
func (h *handler) getHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	url, err := h.app.GetByShort(key)
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

func (h *handler) getUserURLs(w http.ResponseWriter, r *http.Request) {

	uid, ok := r.Context().Value(contextUserID).(string)

	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, render.M{"error": "not authorized"})
		return
	}

	urls, err := h.app.SelectByUID(uid)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	for i := 0; i < len(urls); i++ {
		urls[i].Short = fmt.Sprintf("%s/%s", h.baseURL, urls[i].Short)
	}

	if urls == nil {
		render.NoContent(w, r)
		return
	}

	render.JSON(w, r, urls)
}

func (h *handler) pingDatabase(w http.ResponseWriter, r *http.Request) {
	if err := h.app.PingDatabase(); err != nil {
		render.Status(r, http.StatusInternalServerError)
	} else {
		render.Status(r, http.StatusOK)
	}
}

func (h *handler) deleteHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := r.Context().Value(contextUserID).(string)

	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, render.M{"error": "not authorized"})
		return
	}
	var shorts []string
	render.DecodeJSON(r.Body, &shorts)

	go h.app.Delete(uid, shorts...)

	render.Status(r, http.StatusAccepted)
	render.PlainText(w, r, "")
}
