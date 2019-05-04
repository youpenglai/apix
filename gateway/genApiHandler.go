package gateway

import (
	apiXHttp "github.com/youpenglai/apix/http"
	"github.com/youpenglai/apix/apibuilder"
)

func genApiParamsHandler(params []*apibuilder.ApiParam) (handler apiXHttp.Handler, err error) {
	return
}


// Api代码生成
func GenApiHandles (doc *apibuilder.ApiDoc) (handler apiXHttp.Handler, err error) {
	return func(ctx *apiXHttp.Context) {

	}, nil
}


