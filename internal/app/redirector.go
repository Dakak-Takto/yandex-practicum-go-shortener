package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//Search original URL by key and write 301 redirect
func (a *application) Redirector(c *gin.Context) {
	key := c.Param("key")

	location, err := a.repository.Get(key)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, location.String())
}
