package handlers

import (
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SetLocationHeader sets the Location header for created resources.
// Follows REST convention of providing URI to newly created resource.
func SetLocationHeader(c *gin.Context, id int64) {
	location := path.Join(c.Request.URL.Path, strconv.FormatInt(id, 10))
	c.Header("Location", location)
}
