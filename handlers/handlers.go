package handlers

import (
	"fmt"
	"hash/crc32"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"yandex-practicum-go-shortener/config"
	"yandex-practicum-go-shortener/storage"

	"github.com/gin-gonic/gin"
)

func GetHandler(c *gin.Context) {
	query := c.Param("key")
	link, err := storage.Get(query)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, link)
}

func PostHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	URL := string(body)
	URL = strings.TrimSpace(URL)
	if URL == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	key := generateKey(URL)
	err = storage.Set(key, URL)
	if err != nil {
		c.String(http.StatusInternalServerError, "%s", err)
		return
	}
	c.String(http.StatusCreated, "%s://%s/%s", config.GetScheme(), config.GetAddr(), key)
}

func generateKey(str string) string {
	b := []byte(str)
	hash := crc32.ChecksumIEEE(b)
	hash += uint32(time.Now().UnixMicro()) + rand.Uint32()

	return fmt.Sprintf("%x", hash)
}
