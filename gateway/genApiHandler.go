package gateway

import (
	apiXHttp "github.com/youpenglai/apix/http"
	"github.com/youpenglai/apix/apibuilder"
)

// Api代码生成
func GenApiHandle (doc *apibuilder.ApiEntry) (handler apiXHttp.Handler, err error) {
	return func(ctx *apiXHttp.Context) {

	}, nil
}


