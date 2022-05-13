package app

import (
	"net/http"

	"github.com/go-chi/render"
)

func (a application) deleteHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := r.Context().Value(ctxValueNameUID).(string)

	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, render.M{"error": "not authorized"})
		return
	}
	_ = uid
}
