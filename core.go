package apix

import (
	"net/http"
	"sync"
)

const (
	ApiXName = "ApiX"
)

type ApiX struct {
	Router

	pool *sync.Pool
}


func buildHandleChain(ctx *Context, handler... Handler) {
	n := len(handler)
	cur := 0

	next := func () {
		if cur < n {
			nextHandle := handler[cur]
			cur++
			nextHandle(ctx)
		}
	}

	ctx.Next = next
}

func (apix *ApiX) handleHTTP(ctx *Context) {
	uri := ctx.ResponseURI()
	handlers, params, err := apix.match(uri, ctx.Method())
	if err != nil {
		errHandle(err, ctx)
		return
	}

	ctx.SetParams(params)
	buildHandleChain(ctx, handlers...)
	ctx.Next()
}

func errHandle(err error, ctx *Context) {

}

func (apix *ApiX) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, _ := apix.pool.Get().(*Context)

	ctx.ResponseWriter = w
	ctx.Request = r

	apix.handleHTTP(ctx)

	apix.pool.Put(ctx)
}

func (apix *ApiX) Run(bindAddr string) {
	http.ListenAndServe(bindAddr, apix)
}

func NewApiX() *ApiX {
	apix := &ApiX{}
	apix.pool = &sync.Pool{
		New: func() interface{} {
			return &Context{}
		},
	}

	apix.Use(server)

	return apix
}
