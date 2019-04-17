package http

import (
	"net/http"
	"sync"
	"runtime"
	"fmt"
	"context"
)

const (
	ApiXName = "ApiX"
	ApiXVersion = "0.0.1"
)

var OSName = runtime.GOOS

type ApiX struct {
	Router

	pool *sync.Pool
	server *http.Server
}

func buildHandleChain(ctx *Context, err error, handler... Handler) {
	if err != nil {
		handler = append(handler, errHandle)
	}
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
	uri := ctx.RequestURL()
	handlers, params, err := apix.match(uri, ctx.Method())
	ctx.SetError(err)
	ctx.SetParams(params)
	ctx.parseQueries()

	buildHandleChain(ctx, err, handlers...)
	ctx.Next()
}

func notFoundHandler(ctx *Context) {
	ctx.WriteString(http.StatusNotFound,
		fmt.Sprintf("%s %s (%s)\n404: url not found",
			ApiXName, ApiXVersion, OSName))
}

func notAllowHandler(ctx *Context) {
	ctx.WriteString(http.StatusMethodNotAllowed,
		fmt.Sprintf("%s %s (%s)\nMethod: %s not allowed",
			ApiXName, ApiXVersion, OSName, ctx.Method()))
}

func internalErrorHandler(ctx *Context) {
	// print error stack
	ctx.WriteString(http.StatusInternalServerError,
		fmt.Sprintf("%s %s (%s)\nInternal server error: %s",
			ApiXName, ApiXVersion, OSName, ctx.err.Error()))
}

func requestErrorHandler(ctx *Context) {
	ctx.WriteString(http.StatusBadRequest,
		fmt.Sprintf("%s %s (%s)\nBad request",
			ApiXName, ApiXVersion, OSName))
}

func errHandle(ctx *Context) {
	switch ctx.err {
	case nil:
		return
	case ErrMethodNotFound:
		notAllowHandler(ctx)
	case ErrParamNotExists:
		requestErrorHandler(ctx)
	case ErrRouterNotFound:
		notFoundHandler(ctx)
	default:
		internalErrorHandler(ctx)
	}
}

func (apix *ApiX) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, _ := apix.pool.Get().(*Context)

	ctx.ResponseWriter = w
	ctx.Request = r

	apix.handleHTTP(ctx)

	apix.pool.Put(ctx)
}

func (apix *ApiX) Run(bindAddr string) (err error) {
	apix.server = &http.Server{Addr:bindAddr, Handler:apix}
	return apix.server.ListenAndServe()
}

func (apix *ApiX) Shutdown() {
	if (apix.server != nil) {
		ctx := context.Background()
		apix.server.Shutdown(ctx)
	}
}

func NewApiX() *ApiX {
	apix := &ApiX{}
	apix.pool = &sync.Pool{
		New: func() interface{} {
			return &Context{}
		},
	}

	apix.Use(NewLogger())

	return apix
}
