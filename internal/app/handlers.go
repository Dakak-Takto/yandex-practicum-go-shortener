package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const (
	keyLenghtStart = 8
)

//search exist short url in storage,return temporary redirect if found
func (app *application) GetHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	url, err := app.store.First(key)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.Value, http.StatusTemporaryRedirect)
}

func (app *application) getUserURLs(w http.ResponseWriter, r *http.Request) {

	uid, ok := r.Context().Value(uidContext("uid")).(uidContext)

	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, render.M{"error": "not authorized"})
		return
	}

	log.Printf("getUserURLs handler. uid: %s", uid)

	urls := app.store.Get(string(uid))

	if urls == nil {
		render.NoContent(w, r)
		return
	}

	type item struct {
		Short    string `json:"short_url"`
		Original string `json:"original_url"`
	}

	var result []item
	for _, u := range urls {
		u, err := app.store.First(u.Value)
		if err != nil {
			continue
		}
		result = append(result, item{
			Short:    fmt.Sprintf("%s/%s", app.baseURL, u.Key),
			Original: u.Value,
		})
	}
	render.JSON(w, r, result)
}

//accept json, make short url, write in storage, return short url
func (app *application) PostHandler(w http.ResponseWriter, r *http.Request) {

	uid, ok := r.Context().Value(uidContext("uid")).(uidContext)

	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, render.M{"error": "not authorized"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var request struct {
		URL string `json:"url"`
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{"error": "no valid json found"})
		return
	}

	parsedURL, err := url.ParseRequestURI(request.URL)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{"error": "no valid url found"})
	}

	app.store.Lock()
	defer app.store.Unlock()

	key := app.generateKey(keyLenghtStart)

	app.store.Insert(key, parsedURL.String())
	log.Printf("save short %s -> %s", key, parsedURL.String())
	app.store.Insert(string(uid), key)
	log.Printf("save %s, %s", uid, key)

	result := fmt.Sprintf("%s/%s", app.baseURL, key)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, render.M{"result": result})
}

//accept text/plain body with url, make short url, write in storage, return short url in body
func (app *application) LegacyPostHandler(w http.ResponseWriter, r *http.Request) {

	uid, ok := r.Context().Value(uidContext("uid")).(uidContext)

	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, render.M{"error": "not authorized"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parsedURL, err := url.ParseRequestURI(string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app.store.Lock()
	defer app.store.Unlock()

	key := app.generateKey(keyLenghtStart)

	app.store.Insert(key, parsedURL.String())
	result := fmt.Sprintf("%s/%s", app.baseURL, key)
	app.store.Insert(string(uid), key)

	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, result)
}
