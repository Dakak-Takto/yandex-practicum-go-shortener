package app

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

//search exist short url in storage,return temporary redirect if found
func (app *application) getHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	url, err := app.store.GetByShort(key)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if url.Deleted {
		render.Status(r, http.StatusGone)
		render.PlainText(w, r, "")
	}
	http.Redirect(w, r, url.Original, http.StatusTemporaryRedirect)
}

func (app *application) getUserURLs(w http.ResponseWriter, r *http.Request) {

	uid, ok := r.Context().Value(ctxValueNameUID).(string)

	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, render.M{"error": "not authorized"})
		return
	}

	urls, err := app.store.SelectByUID(uid)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	for i := 0; i < len(urls); i++ {
		urls[i].Short = fmt.Sprintf("%s/%s", app.baseURL, urls[i].Short)
	}

	if urls == nil {
		render.NoContent(w, r)
		return
	}

	render.JSON(w, r, urls)
}

func (app *application) pingDatabase(w http.ResponseWriter, r *http.Request) {
	if err := app.store.Ping(); err != nil {
		render.Status(r, http.StatusInternalServerError)
	} else {
		render.Status(r, http.StatusOK)
	}
}
