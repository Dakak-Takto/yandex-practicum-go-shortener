package handler

import (
	"context"
	"net/http"

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
			session.Save(r, w)
		}

		ctx := context.WithValue(r.Context(), userIDctxKeyName, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
