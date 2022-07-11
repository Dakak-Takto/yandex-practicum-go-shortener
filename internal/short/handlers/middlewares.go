package handler

import (
	"compress/gzip"
	"context"
	"net/http"
	"strings"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/pkg/random"
)

type userIDctxKey string

const userIDctxKeyName userIDctxKey = "user_id"

const (
	sessionCookieName      string = "session"
	userIDSessionValueName string = "user_id"
)

func (h *handler) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := h.cookieStore.Get(r, sessionCookieName)
		if err != nil && !session.IsNew {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		userID, ok := session.Values[userIDSessionValueName].(string)

		if !ok {
			userID = random.String(5)
			session.Values[userIDSessionValueName] = userID
			err := session.Save(r, w)
			if err != nil {
				h.log.Warn(err)
			}
			h.log.Debugf("new user")
		}

		h.log.Debugf("userID: %s", userID)

		ctx := context.WithValue(r.Context(), userIDctxKeyName, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *handler) decompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			h.log.Debug("not contains Content-Encoding: gzip header. Continue.")
			next.ServeHTTP(w, r)
			return
		}

		h.log.Debug("Content-Encoding: gzip header. Try decompress.")
		gzReader, err := gzip.NewReader(r.Body)
		if err != nil {
			h.log.Debug("error read body:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = gzReader.Close()
		if err != nil {
			h.log.Debug("error close reader:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.Body = gzReader

		next.ServeHTTP(w, r)
	})
}
