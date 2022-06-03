package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"

	"yandex-practicum-go-shortener/internal/short/model"
)

type handler struct {
	usecase     model.ShortUsecase
	cookieStore *sessions.CookieStore
}

func New(usecase model.ShortUsecase, cookieStore *sessions.CookieStore) *handler {
	return &handler{
		usecase:     usecase,
		cookieStore: cookieStore,
	}
}

func (h *handler) Register(router *chi.Mux) {
	router.Use(h.auth)
	router.Route("/", func(r chi.Router) {
		router.Get("/{key}", h.getRedirect)
		// router.Get("/ping", app.pingDatabase)
		router.Post("/", h.makeShort)
	})

	// router.Route("/api/shorten", func(r chi.Router) {
	// 	r.Post("/", app.postHandler)
	// 	r.Post("/batch", app.batchPostHandler)
	// })

	// router.Route("/api/user/urls", func(r chi.Router) {
	// 	r.Get("/", app.getUserURLs)
	// 	r.Delete("/", app.deleteHandler)
	// })

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
	return
}

type userShortResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (h *handler) getUserShorts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

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
			ShortURL:    fmt.Sprintf("%s/%s", "base_url", short.Key),
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
	userID := r.Context().Value("user_id").(string)

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

		result.Result = fmt.Sprintf("%s/%s", "base_url", short.Key)

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Printf("error encode result: %s", err)
		}
		return
	}

	locationURL, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	short, err := h.usecase.CreateNewShort(string(locationURL), userID)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(short.Location))
}
