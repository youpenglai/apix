package apix

type Middleware interface {
	Handle(c *Context)
}
