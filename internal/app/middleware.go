package app

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type gzipWriter struct {
	gin.ResponseWriter
	writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

func (w gzipWriter) WriteString(s string) (int, error) {
	return w.writer.Write([]byte(s))
}

func gzipMiddleware(c *gin.Context) {

	if strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
		c.Header("Content-Encoding", "gzip")
		writer, err := gzip.NewWriterLevel(c.Writer, gzip.BestCompression)
		if err != nil {
			log.Fatal(err)
		}
		defer writer.Close()

		c.Writer = gzipWriter{ResponseWriter: c.Writer, writer: writer}
	}

	if strings.Contains(c.Request.Header.Get("Content-Encoding"), "gzip") {
		reader, err := gzip.NewReader(c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		defer reader.Close()
		body, _ := ioutil.ReadAll(reader)

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	c.Next()
}
