package apix

func server(c *Context) {
	c.SetHeader("server", ApiXName)
	c.Next()
}

//func etag(c *Context) {
//	c.Next(c)
//}