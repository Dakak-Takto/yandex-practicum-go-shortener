package handlers

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/pkg/random"
)

const (
	newUserIDLenght   int    = 10
	cookieSessionName string = `session`
	sessionUserIDKey  string = `userID`
)

func (h *handler) SetCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, err := h.cookieStore.Get(r, cookieSessionName)
		if err != nil {
			h.log.Warn(err)
		}

		if session.IsNew {
			h.log.Debug("is new session. generate new userID")
			session.Values["userID"] = random.String(newUserIDLenght)
			if err := session.Save(r, w); err != nil {
				h.log.Warn(err)
				http.Error(w, "", http.StatusBadRequest)
				return
			}
		}

		userID := session.Values[sessionUserIDKey].(string)

		h.log.Debugf("is old session. userID = %s", userID)
		ctx := context.WithValue(r.Context(), contextValueKeyUserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *handler) decompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gzReader, err := gzip.NewReader(r.Body)
		if err != nil {
			h.log.Warn("error read body:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = gzReader.Close()
		if err != nil {
			h.log.Warn("error close reader:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.Body = gzReader

		next.ServeHTTP(w, r)
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

func (h *handler) httpLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.log.Warn(err)
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

		h.log.Debugf("REQT: %s %s", r.Method, r.RequestURI)
		h.log.Debugf("%s", body)
		h.log.Debugf("RESP: %d %s %s", rec.Status, r.RequestURI, rec.response)
	})
}
