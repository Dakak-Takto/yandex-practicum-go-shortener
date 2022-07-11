package app

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/hex"
	"io"
	"net/http"
	"strings"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/random"
)

type ctxValueTypeUID string

var (
	ctxValueNameUID ctxValueTypeUID = "uid"
	cookieNameToken                 = "token"
)

func (app *application) decompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			app.log.Debug("not contains Content-Encoding: gzip header. Continue.")
			next.ServeHTTP(w, r)
			return
		}

		app.log.Debug("Content-Encoding: gzip header. Try decompress.")
		gzReader, err := gzip.NewReader(r.Body)
		if err != nil {
			app.log.Debug("error read body:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = gzReader.Close()
		if err != nil {
			app.log.Debug("error close reader:", err)
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
				app.log.Debug("cookie: token not found.", err)
				return "", err
			}

			decoded := make(map[string]string)
			err = app.secureCookie.Decode("token", cookie.Value, &decoded)
			if err != nil {
				app.log.Debug("cookie: token decode failed.", err)
				return "", err
			}

			uid := decoded["uid"]

			app.log.Debug("cookie: uid:", uid)
			return uid, nil
		}()

		/*
		   если куков нет или они не прошли проверку, генерируем новые куки
		*/
		if err != nil {
			uidBytes, err := random.RandomBytes(8)
			if err != nil {
				app.log.Debugf("error generate token: %s", err)
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
				app.log.Debugf("error encode token: %s", err)
			}
		}
		/*
			порядок навели, передаем дальще
		*/
		ctx := context.WithValue(r.Context(), ctxValueNameUID, uid)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

type recorder struct {
	http.ResponseWriter
	response []byte
	Status   int
}

func (r *recorder) Write(b []byte) (int, error) {
	r.response = b
	return r.ResponseWriter.Write(b)
}

func (r *recorder) WriteHeader(statusCode int) {
	r.Status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (app *application) httpLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			app.log.Warn(err)
			next.ServeHTTP(w, r)
			return
		}
		defer r.Body.Close()
		reader := io.NopCloser(bytes.NewBuffer(body))

		r.Body = reader

		rec := &recorder{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}

		next.ServeHTTP(rec, r)

		app.log.Debugf("REQT: %s %s", r.Method, r.RequestURI)
		app.log.Debugf("%s", body)
		app.log.Debugf("RESP: %d %s %s", rec.Status, r.RequestURI, rec.response)

	})
}
