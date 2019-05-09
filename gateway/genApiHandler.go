package gateway

import (
	apiXHttp "github.com/youpenglai/apix/http"
	"github.com/youpenglai/apix/apibuilder"
	"io/ioutil"
	"encoding/json"
	"strings"
)

type paramReader struct {
	ctx *apiXHttp.Context
	bodyCache map[string]interface{}
}

func (r *paramReader) Get(name, from string) (interface{}) {
	hasBody := false
	var bodyJSON map[string]interface{}
	if r.bodyCache == nil {
		method := strings.ToLower(r.ctx.Method())
		if method == "put" || method == "post" {
			hasBody = true
			defer r.ctx.Request.Body.Close()
			body, err := ioutil.ReadAll(r.ctx.Request.Body)
			if err != nil {
				return nil
			}
			if err = json.Unmarshal(body, &bodyJSON); err != nil {
				return nil
			}
		}
	} else {
		hasBody = true
		bodyJSON = r.bodyCache
	}

	switch from {
	case "body":
		if !hasBody {
			return nil
		}
		if v, ok := bodyJSON[name]; !ok {
			return nil
		} else {
			return v
		}
	case "path":
		params := r.ctx.Params()
		if v, ok := params[name]; !ok {
			return nil
		} else {
			return v
		}
	case "queries":
		params := r.ctx.Queries()
		if v, ok := params[name]; !ok {
			return nil
		} else {
			return v
		}
	case "header":
		return r.ctx.Header().Get(name)
	}

	return nil
}

func grpcHandler(ctx *apiXHttp.Context) {

}

func httpHandler(ctx *apiXHttp.Context) {
	
}

// Api代码生成
func GenApiHandle (block *apibuilder.ApiCodeBlock) (handler apiXHttp.Handler) {
	return func(ctx *apiXHttp.Context) {
		reader := &paramReader{ctx:ctx}
		block.BindParamReader(reader)
		params, err := block.ReadParams()
		if err != nil {
			// TODO: params err
		}

		err = params.Validation()

		ctx.WriteString(200, "Api Handle")
		//ctx.Next()
	}
}


