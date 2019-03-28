package apix

import "net/http"

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request

	// TODO: add more
	writen int
}

func (c *Context) reset(w http.ResponseWriter, r *http.Request) {
	c.ResponseWriter = w
	c.Request = r
	c.writen = 0
}

func (c *Context) GetResponseWriter() http.ResponseWriter {
	return c.ResponseWriter
}

func (c *Context) Write(data []byte) error {
	n, err := c.ResponseWriter.Write(data)
	c.writen += n

	return err
}

func (c *Context) WriteString(str string) error {
	return c.Write([]byte(str))
}
