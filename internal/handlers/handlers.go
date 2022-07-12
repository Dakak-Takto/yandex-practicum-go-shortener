package handlers

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/entity"
)

const (
	contextValueKeyUserID string = "userID"
)

type handler struct {
	usecase     entity.URLUsecase
	log         *logrus.Logger
	baseURL     string
	cookieStore *sessions.CookieStore
}

func New(usecase entity.URLUsecase, log *logrus.Logger, baseURL string, cookieStore *sessions.CookieStore) *handler {
	return &handler{
		usecase:     usecase,
		log:         log,
		baseURL:     baseURL,
		cookieStore: cookieStore,
	}
}

func (h *handler) Register(mux *chi.Mux) {

	mux.Use(middleware.Compress(gzip.BestCompression, "application/*", "text/*"))
	mux.Use(h.decompress)
	mux.Use(h.SetCookie)
	mux.Use(h.httpLog)

	//Routes
	mux.Route("/", func(r chi.Router) {
		r.Get("/{short}", h.shortToOriginal)
		// r.Get("/ping", nil)
		r.Post("/", h.makeShortPlainText)
	})

	mux.Route("/api/shorten", func(r chi.Router) {
		r.Post("/", h.makeShort)
		r.Post("/batch", h.batchMakeShorts)
	})

	mux.Route("/api/user/urls", func(r chi.Router) {
		r.Get("/", h.getUserURLs)
		r.Delete("/", h.delete)
	})
}

func (h *handler) shortToOriginal(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "short")

	url, err := h.usecase.GetByShort(key)

	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			http.NotFound(w, r)
		}
	}

	if url.Deleted {
		http.Error(w, "", http.StatusGone)
		return
	}

	http.Redirect(w, r, url.Original, http.StatusTemporaryRedirect)
}

func (h *handler) getUserURLs(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(contextValueKeyUserID).(string)

	urls, err := h.usecase.UserURLs(uid)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			http.Error(w, "", http.StatusNoContent)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := make([]userURL_DTO, len(urls))

	for i := 0; i < len(urls); i++ {
		response[i] = userURL_DTO{
			Original: urls[i].Original,
			Short:    fmt.Sprintf("%s/%s", h.baseURL, urls[i].Short),
		}
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Warn(err)
	}
}

func (h *handler) makeShort(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(contextValueKeyUserID).(string)

	var req makeShortRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decode json", http.StatusBadRequest)
		return
	}

	var result resultDTO

	url, err := h.usecase.Create(req.URL, userID)
	if err != nil {
		if errors.Is(err, entity.ErrDuplicate) {
			url, err = h.usecase.GetByOriginal(url.Original)
			if err != nil {
				http.Error(w, "", http.StatusBadRequest)
				return
			}
			result.Result = fmt.Sprintf("%s/%s", h.baseURL, url.Short)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			if err := json.NewEncoder(w).Encode(&result); err != nil {
				return
			}
			return
		}
		http.Error(w, "error create short", http.StatusBadRequest)
		return
	}

	result.Result = fmt.Sprintf("%s/%s", h.baseURL, url.Short)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&result); err != nil {
		h.log.Warn(err)
	}
}

func (h *handler) makeShortPlainText(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(contextValueKeyUserID).(string)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Warn(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	original := string(body)

	url, err := h.usecase.Create(original, userID)
	if err != nil {
		if errors.Is(err, entity.ErrDuplicate) {
			url, err = h.usecase.GetByOriginal(url.Original)
			if err != nil {
				http.Error(w, "", http.StatusBadRequest)
				return
			}
			result := fmt.Sprintf("%s/%s", h.baseURL, url.Short)
			w.WriteHeader(http.StatusConflict)

			if _, err := w.Write([]byte(result)); err != nil {
				return
			}
			return
		}
		http.Error(w, "error create short", http.StatusBadRequest)
		return
	}

	result := fmt.Sprintf("%s/%s", h.baseURL, url.Short)

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(result)); err != nil {
		h.log.Warn(err)
	}
}

func (h *handler) delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(contextValueKeyUserID).(string)

	var shorts []string
	if err := json.NewDecoder(r.Body).Decode(&shorts); err != nil {
		h.log.Warn(err)
		return
	}

	go h.usecase.Delete(userID, shorts...)

	w.WriteHeader(http.StatusAccepted)
}

func (h *handler) batchMakeShorts(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(contextValueKeyUserID).(string)

	var req []batchItemRequestURLs

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warn(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	var batchResponseURLs []batchItemResponse

	for _, batchItem := range req {
		url, err := h.usecase.Create(batchItem.OriginalURL, userID)
		if err != nil {
			h.log.Warn(err)
			return
		}
		shortURL := fmt.Sprintf("%s/%s", h.baseURL, url.Short)

		batchResponseURLs = append(batchResponseURLs, batchItemResponse{
			CorellationID: batchItem.CorellationID,
			ShortURL:      shortURL,
		})
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&batchResponseURLs); err != nil {
		h.log.Warn(err)
	}
}
