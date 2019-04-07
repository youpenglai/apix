package http

import "testing"

func TestNewApiX(t *testing.T) {
	t.Log("success")
}

func TestApiX_Run(t *testing.T) {
	apix := NewApiX()
	apix.Use(NewLogger(LoggerOpts{}))
	apix.Get("/", func(ctx *Context) {
		ctx.WriteString(200, "Hello, ApiX")
	})
	apix.Get("/mytest", func(ctx *Context) {
		ctx.WriteString(200, "Hello, ApiX MyTestHandler")
	})
	apix.Get("/params/:id", func(ctx *Context) {
		println("handles")
		params := ctx.Params()
		querys := ctx.Queries()
		panic("test")
		ctx.JSON(200, map[string]interface{}{
			"params": params,
			"queries": querys,
		})
	})
	apix.Run("127.0.0.1:8080")
}
