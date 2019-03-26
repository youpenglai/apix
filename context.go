package apix

import "net/http"

type Context struct {
	Writer http.ResponseWriter
	Request *http.Request

	// TODO: add more
	writen int
}

func (c *Context) reset(w http.ResponseWriter, r *http.Request) {
	c.Writer = w
	c.Request = r
	c.writen = 0
}

func (c *Context) GetResponseWriter() http.ResponseWriter {
	return c.Writer
}

func (c *Context) Write(data []byte) error {
	n, err := c.Writer.Write(data)
	c.writen += n

	return err
}

func (c *Context) WriteString(str string) error {
	return c.Write([]byte(str))
}
