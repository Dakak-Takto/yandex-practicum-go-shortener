package app

import (
	"net/http"

	"github.com/go-chi/render"
)

func (app *application) deleteHandler(w http.ResponseWriter, r *http.Request) {
	uid, ok := r.Context().Value(ctxValueNameUID).(string)

	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, render.M{"error": "not authorized"})
		return
	}
	var shorts []string
	render.DecodeJSON(r.Body, &shorts)

	go app.store.Delete(uid, shorts...)

	render.Status(r, http.StatusAccepted)
	render.PlainText(w, r, "")
}
