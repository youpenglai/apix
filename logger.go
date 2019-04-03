package apix

type Logger struct {

}

type LoggerOpts struct {
	Name string
	Level string
}

func New(opts LoggerOpts) Handler{
	return func(c *Context) {
		c.Next()
	}
}