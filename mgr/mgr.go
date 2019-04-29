package mgr

import (
	"github.com/youpenglai/apix/http"
	"github.com/youpenglai/apix/middlewares"
)

func runManagerServer() {
	mgrServer := http.NewApiX()

	mgrServer.Use(middlewares.Server())

	mgrServer.Run("127.0.0.1:8080")
}

func init() {
	go runManagerServer()
}