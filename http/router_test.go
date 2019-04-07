package http

import (
	"testing"
)

func TestRouter_Group(t *testing.T) {
	r := &Router{}
	g := r.Group("/root/group")

	if g.name != "group" {
		t.Error("group test error")
		return
	}
	t.Log("group test success")
}

func TestRouter_Get(t *testing.T) {
	r := &Router{}
	g := r.Group("/root/group")
	g.Use(func(ctx *Context) {
		t.Log("middleware 1")
	}, func(ctx *Context) {
		t.Log("middleware 2")
	})

	g.Get(":id/info", func(ctx *Context) {
		t.Log("get function")
	})

	handlers, params, err := r.match("/root/group/1/info", "GET")
	if len(handlers) != 3 {
		t.Error("invalid handlers:", handlers)
		return
	}
	if len(params) == 0 {
		t.Error("params cannot be empty")
		return
	}
	if err != nil {
		t.Error(err.Error())
	}

	for _, h := range handlers {
		h(nil)
	}

	t.Log("params:", params)

	t.Log("test get success")
}