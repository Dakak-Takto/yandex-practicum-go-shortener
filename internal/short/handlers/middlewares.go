package handler

import (
	"context"
	"net/http"
)

func (h *handler) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := h.cookieStore.Get(r, "session")
		if err != nil && !session.IsNew {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		userID, exist := session.Values["user_id"].(string)

		if !exist {
			userID = "new user"
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
