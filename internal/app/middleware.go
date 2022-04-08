package app

import (
	"compress/gzip"
	"net/http"
	"strings"
)

func (app *application) decompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// не забывайте потом закрыть *gzip.Reader
		defer gz.Close()

		r.Body = gz

		next.ServeHTTP(w, r)

	})
}
