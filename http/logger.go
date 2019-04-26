package http

import (
	"time"
	"github.com/youpenglai/goutils/logger"
	ApixLogger "github.com/youpenglai/apix/logger"
	"fmt"
)

var (
	log = logger.GetLogger(ApixLogger.PrefixAccess)
)

func NewLogger() Handler{
	return func(c *Context) {
		start := time.Now()
		c.Next()
		end := time.Now()

		method := c.Method()

		// TODO: add ipaddr and status code here
		log.Info(fmt.Sprintf("[%s] %0.4fS - %5s", method, end.Sub(start).Seconds(), c.Request.URL.RequestURI()))
		//println("end:", end.Unix() , " Used:", end.Sub(start))
	}
}