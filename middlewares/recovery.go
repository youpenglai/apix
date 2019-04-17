package middlewares

import (
	"github.com/youpenglai/apix/http"
)

func stack() {}

func Recovery() http.Handler{
	return func(ctx *http.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctx.WriteString(500, "InternalServerError")
				// track stack
			}
		}()
		ctx.Next()
	}
}
