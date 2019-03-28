package apix

import (
	"net/http"
	"sync"
)


type ApiX struct {
	Router

	pool *sync.Pool
}

func (apix *ApiX) handleHTTP(ctx *Context) {

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

	return apix
}
