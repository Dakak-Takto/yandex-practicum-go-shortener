package app

import (
	"compress/gzip"
	"context"
	"encoding/hex"
	"net/http"
	"strings"
	"yandex-practicum-go-shortener/internal/random"
)

var (
	ctxValueNameUid = "uid"
	cookieNameToken = "token"
)

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

		/*
			пробуем разобрать куки
		*/
		uid, err := func() (string, error) {
			cookie, err := r.Cookie("token")
			if err != nil {
				app.logger.Printf("cookie: token not found.", err)
				return "", err
			}

			decoded := make(map[string]string)
			err = app.secureCookie.Decode("token", cookie.Value, &decoded)
			if err != nil {
				app.logger.Printf("cookie: token decode failed.", err)
				return "", err
			}

			uid := decoded["uid"]

			app.logger.Printf("cookie: uid:", uid)
			return uid, nil
		}()

		/*
		   если куков нет или они не прошли проверку, генерируем новые куки
		*/
		if err != nil {
			uidBytes, err := random.RandomBytes(8)
			if err != nil {
				app.logger.Printf("error generate token: %s", err)
				http.Error(w, "something went wrong", http.StatusInternalServerError)
				return
			}
			uid = hex.EncodeToString(uidBytes)

			cookies := map[string]string{
				"uid": uid,
			}

			if encoded, err := app.secureCookie.Encode("token", cookies); err == nil {
				cookie := &http.Cookie{
					Name:     cookieNameToken,
					Value:    encoded,
					Path:     "/",
					HttpOnly: true,
				}
				http.SetCookie(w, cookie)
			} else {
				app.logger.Printf("error encode token: %s", err)
			}
		}
		/*
			порядок навели, передаем дальще
		*/
		ctx := context.WithValue(r.Context(), ctxValueNameUid, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
		return

	})
}
