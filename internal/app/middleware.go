package app

import (
	"compress/gzip"
	"context"
	"encoding/hex"
	"net/http"
	"strings"
	"yandex-practicum-go-shortener/internal/random"
)

type uidContext string

func (u *uidContext) String() string {
	return string(*u)
}

func (app *application) decompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			app.logger.Print("not contains Content-Encoding: gzip header. Continue.")
			next.ServeHTTP(w, r)
			return
		}

		app.logger.Print("Content-Encoding: gzip header. Try decompress.")
		gzReader, err := gzip.NewReader(r.Body)
		if err != nil {
			app.logger.Print("error read body:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = gzReader.Close()
		if err != nil {
			app.logger.Print("error close reader:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		r.Body = gzReader

		next.ServeHTTP(w, r)
	})
}

func (app *application) SetCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if cookie, err := r.Cookie("token"); err == nil {
			value := make(map[string]string)
			if err := app.secureCookie.Decode("token", cookie.Value, &value); err == nil {
				ctx := context.WithValue(r.Context(), uidContext("uid"), uidContext(value["uid"]))
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			} else {
				app.logger.Printf("error decode token cookie: %s", err)

			}
		} else {
			app.logger.Printf("error take token cookie: %s", err)
		}

		uidBytes, err := random.RandomBytes(8)
		if err != nil {
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		uid := hex.EncodeToString(uidBytes)

		app.logger.Printf("New UID: %s", uid)

		value := map[string]string{
			"uid": uid,
		}

		if encoded, err := app.secureCookie.Encode("token", value); err == nil {
			cookie := &http.Cookie{
				Name:     "token",
				Value:    encoded,
				Path:     "/",
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)

			ctx := context.WithValue(r.Context(), uidContext("uid"), uidContext(uid))
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
	})
}
