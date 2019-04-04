package apix

import "fmt"

func server(c *Context) {
	c.SetHeader("server", fmt.Sprintf("%s %s (%s)", ApiXName, ApiXVersion, OSName))
	c.Next()
}

//func etag(c *Context) {
//	c.Next(c)
//}