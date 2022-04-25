package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

//search exist short url in storage,return temporary redirect if found
func (app *application) GetHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	url, err := app.store.First(key)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.Original, http.StatusTemporaryRedirect)
}

func (app *application) getUserURLs(w http.ResponseWriter, r *http.Request) {

	uid, ok := r.Context().Value(uidContext("uid")).(uidContext)

	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, render.M{"error": "not authorized"})
		return
	}

	log.Printf("getUserURLs handler. uid: %s", uid)

	urls := app.store.GetByUID(uid.String())

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
