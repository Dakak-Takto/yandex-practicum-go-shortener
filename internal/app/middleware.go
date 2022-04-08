package app

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type gzipWriter struct {
	gin.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Gzip() gin.HandlerFunc {
	return func(c *gin.Context) {

		if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}
		gz, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		defer gz.Close()
		c.Writer.Header().Set("Content-Encoding", "gzip")

		c.Writer = gzipWriter{ResponseWriter: c.Writer, Writer: gz}

		if !strings.Contains(c.Request.Header.Get("Content-Encoding"), "gzip") {
			c.Next()
			return
		}

		reader, err := gzip.NewReader(c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		defer reader.Close()
		c.Request.Body = reader
		c.Next()
	}
}
