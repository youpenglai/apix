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
			r.bodyCache = bodyJSON
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

var forwardFuncs = map[string]func(string, interface{} , map[string]interface{})([]byte, error){
	"grpc": func(service string, target interface{}, params map[string]interface{}) ([]byte, error) {
		p, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}

		grpcTarget := target.(*apibuilder.GRPCForward)

		result, err := CallService(service, grpcTarget.Method, p)
		return result, err
	},
	"http": func(service string, target interface{}, i map[string]interface{}) ([]byte, error) {
		return nil, nil
	},
	"redis": func(service string, target interface{}, keys map[string]interface{}) ([]byte, error) {
		p, err := json.Marshal(keys)
		if err != nil {
			return nil, err
		}

		result, err := CallService(service, "", p)
		return result, err
	},
}

type forwardImpl struct {}

func (fi *forwardImpl) ForwardTo(dest *apibuilder.ApiForwards, mapper map[string]interface{}) (ret []byte, err error) {
	ff, _ := forwardFuncs[dest.TargetType]
	ret, err = ff(dest.Service, dest.TargetInfo, mapper)
	return
}

// Api代码生成
func GenApiHandle (code *apibuilder.ApiCodeBlock) (handler apiXHttp.Handler) {
	return func(ctx *apiXHttp.Context) {
		reader := &paramReader{ctx:ctx}
		code.BindParamReader(reader)
		code.BindForwardImpl(&forwardImpl{})
		params, err := code.ReadParams()
		if err != nil {
			// TODO: params err
		}

		if err = params.Validation(); err != nil {
			// TODO: err process
		}

		var ret interface{}
		if ret, err = code.DoForwards(params); err != nil {
			ctx.JSON(500, map[string]interface{}{"success": false})
		} else {
			ctx.RawBytes(200,"application/json", ret.([]byte))
		}
	}
}


