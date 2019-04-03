package apix

import "testing"

func TestNewApiX(t *testing.T) {
	t.Log("success")
}

func TestApiX_Run(t *testing.T) {
	apix := NewApiX()
	apix.Get("/", func(ctx *Context) {
		ctx.WriteString(200, "Hello, ApiX")
	})
	apix.Get("/mytest", func(ctx *Context) {
		ctx.WriteString(200, "Hello, ApiX MyTestHandler")
	})
	apix.Run("127.0.0.1:8080")
}
