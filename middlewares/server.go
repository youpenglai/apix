package middlewares

import (
	"fmt"
	"github.com/youpenglai/apix/http"
)

func Server() http.Handler {
	return func(c *http.Context) {
		c.SetHeader("server", fmt.Sprintf("%s %s (%s)", http.ApiXName, http.ApiXVersion, http.OSName))
		c.Next()
	}
}
