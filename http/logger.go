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
		log.Info(fmt.Sprintf("[%s] - %s, Used: %0.4f", method, c.Request.URL.RequestURI(), end.Sub(start).Seconds()))
		//println("end:", end.Unix() , " Used:", end.Sub(start))
	}
}