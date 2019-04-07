package http

import "time"

type LoggerOpts struct {
	Name string
	Level string
}

func NewLogger(opts LoggerOpts) Handler{
	return func(c *Context) {
		start := time.Now()
		println("start:", start.Unix())
		c.Next()
		end := time.Now()
		println("end:", end.Unix() , " Used:", end.Sub(start))
	}
}