package apix

import (
	"net/http"
	"encoding/json"
	"io"
	"path"
	"mime"
)

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request

	Next Handler
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

func (c *Context) Flush() {
	if flusher, ok := c.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (c *Context) Write(statusCode int, data []byte) error {
	c.ResponseWriter.WriteHeader(statusCode)
	n, err := c.ResponseWriter.Write(data)
	c.writen += n

	return err
}

func (c *Context) WriteString(statusCode int, str string) error {
	return c.Write(statusCode, []byte(str))
}

func (c *Context) SetHeader(name, value string) {
	c.ResponseWriter.Header().Set(name, value)
}

func (c *Context) JSON(statusCode int, obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	c.SetHeader("Content-Type", "application/json")
	return c.Write(statusCode, data)
}

type WriteFileHandler func(w io.Writer) error

func (c *Context) WriteFile(fileName string, writerHandler WriteFileHandler) error {
	ext := path.Ext(fileName)
	mimeType := mime.TypeByExtension(ext)

	c.ResponseWriter.WriteHeader(http.StatusOK)
	c.SetHeader("Content-Type", mimeType)

	if writerHandler == nil {
		panic("no writer handle")
	}

	writer, _ :=  c.ResponseWriter.(io.Writer)
	if err := writerHandler(writer); err != nil {
		return err
	}

	c.Flush()
	return nil
}