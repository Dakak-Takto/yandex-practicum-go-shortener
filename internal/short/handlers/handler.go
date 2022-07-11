package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/short/model"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/short/usecase"
)

type handler struct {
	usecase     model.ShortUsecase
	baseURL     string
	cookieStore *sessions.CookieStore
	log         *logrus.Logger
}

func New(usecase model.ShortUsecase, baseURL string, cookieStore *sessions.CookieStore, log *logrus.Logger) *handler {
	return &handler{
		usecase:     usecase,
		cookieStore: cookieStore,
		baseURL:     baseURL,
	}
}

func (h *handler) Register(router chi.Router) {
	router.Use(middleware.Compress(5, "text/*", "application/json"))

	router.Use(h.auth)
	router.Route("/", func(r chi.Router) {
		router.Get("/{key}", h.getRedirect)
		router.Get("/ping", h.pingDatabase)
		router.Post("/", h.makeShort)
	})

	router.Route("/api/shorten", func(r chi.Router) {
		r.Post("/", h.makeShort)
		r.Post("/batch", h.makeShortBatch)
	})

	router.Route("/api/user/urls", func(r chi.Router) {
		r.Get("/", h.getUserShorts)
		r.Delete("/", h.deleteShorts)
	})

}

func (h *handler) getRedirect(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	short, err := h.usecase.FindByKey(key)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	if short.Deleted {
		http.Error(w, "", http.StatusGone)
		return
	}
	http.Redirect(w, r, short.Location, http.StatusTemporaryRedirect)
}

type userShortResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (h *handler) getUserShorts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDctxKeyName).(string)

	shorts, err := h.usecase.GetUserShorts(userID)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	if shorts == nil {
		http.Error(w, "", http.StatusNoContent)
		return
	}

	var userShorts []userShortResponse

	for _, short := range shorts {
		userShorts = append(userShorts, userShortResponse{
			ShortURL:    fmt.Sprintf("%s/%s", h.baseURL, short.Key),
			OriginalURL: short.Location,
		})
	}

	if err := json.NewEncoder(w).Encode(userShorts); err != nil {
		log.Printf("error decode json: %s\n", err)
	}
}

type makeShortRequest struct {
	URL string `json:"url"`
}

func (h *handler) makeShort(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDctxKeyName).(string)

	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		var request makeShortRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if _, err := url.Parse(request.URL); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		short, err := h.usecase.CreateNewShort(request.URL, userID)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		var result struct {
			Result string `json:"result"`
		}

		result.Result = fmt.Sprintf("%s/%s", h.baseURL, short.Key)

		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Printf("error encode result: %s", err)
		}
		return
	}

	locationURL, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error read body:", err)
		return
	}

	short, err := h.usecase.CreateNewShort(string(locationURL), userID)
	if err != nil {
		if errors.Is(err, usecase.ErrDuplicate) {
			short, err := h.usecase.FindByLocation(string(locationURL))
			if err != nil {
				log.Println(err)
				http.Error(w, "", http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(short.Location))
			return
		}

		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s/%s", h.baseURL, short.Key)
}

type makeShortBatchRequest struct {
	OriginalURL   string `json:"original_url"`
	CorrelationID string `json:"correlation_id"`
}

func (h *handler) makeShortBatch(w http.ResponseWriter, r *http.Request) {
	var req []makeShortBatchRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (h *handler) deleteShorts(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) pingDatabase(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
