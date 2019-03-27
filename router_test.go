package apix

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