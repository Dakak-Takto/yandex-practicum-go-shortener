package handler

import (
	"context"
	"net/http"
	"yandex-practicum-go-shortener/pkg/random"
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

		userID, exist := session.Values[userIDSessionValueName].(string)

		if !exist {
			userID = random.String(5)
		}

		ctx := context.WithValue(r.Context(), userIDctxKeyName, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
