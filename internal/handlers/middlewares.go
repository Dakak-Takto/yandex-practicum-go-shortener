package handlers

import (
	"compress/gzip"
	"context"
	"net/http"
	"strings"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/random"
)

func (h *handler) decompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			h.log.Debug("not contains Content-Encoding: gzip header. Continue.")
			next.ServeHTTP(w, r)
			return
		}

		h.log.Print("Content-Encoding: gzip header. Try decompress.")
		gzReader, err := gzip.NewReader(r.Body)
		if err != nil {
			h.log.Warn("error read body:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = gzReader.Close()
		if err != nil {
			h.log.Warn("error close reader:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.Body = gzReader

		next.ServeHTTP(w, r)
	})
}

func (h *handler) SetCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, err := h.sessions.Get(r, sessionKeyName)
		if err != nil {
			if !session.IsNew {
				h.log.Warn(err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
		}

		userID := func() string {
			if uid, exist := session.Values[cookieUserIDkeyName]; !exist {
				return ""
			} else {
				if result, ok := uid.(string); ok {
					return result
				}
			}
			return ""
		}()

		if userID == "" {
			session.Values[cookieUserIDkeyName] = random.String(5)
			session.Save(r, w)
		}

		/*
			порядок навели, передаем дальще
		*/
		ctx := context.WithValue(r.Context(), contextUserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
