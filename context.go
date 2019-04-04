package apix

import (
	"net/http"
	"encoding/json"
	"io"
	"path"
	"mime"
	"errors"
	"strconv"
	"strings"
	"net/url"
)

type Params map[string]string

var (
	ErrParamNotExists = errors.New("param not exists")
)

func (p Params) AddValue(key, value string) {
	p[key] = value
}

func (p Params) GetString(key string) (string, error) {
	v, exists := p[key]
	if !exists {
		return "", ErrParamNotExists
	}

	return v, nil
}

func (p Params) GetStringDefault(key string, defaultVal string) string {
	v, err := p.GetString(key)
	if err != nil {
		return defaultVal
	}

	return v
}

func (p Params) GetInt(key string) (int64, error) {
	v, err := p.GetString(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(v, 10, 0)
}

func (p Params) GetIntDefault(key string, defaultVal int64) int64 {
	v, err := p.GetInt(key)
	if err != nil {
		return defaultVal
	}

	return v
}

func (p Params) GetFloat(key string) (float64, error) {
	v, err := p.GetString(key)
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(v, 0)
	if err != nil {
		return 0, err
	}
	return f, nil
}

func (p Params) GetFloatDefault(key string, defaultVal float64) float64 {
	v, err := p.GetFloat(key)
	if err != nil {
		return defaultVal
	}

	return v
}

func NewParams() Params {
	return make(Params)
}

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	params         Params
	queries			Params

	Next func ()
	err error
	// TODO: add more
	writen int
}

func (c *Context) reset(w http.ResponseWriter, r *http.Request) {
	c.ResponseWriter = w
	c.Request = r
	c.writen = 0
	c.params = nil
	c.err = nil
}

func (c *Context) SetParams(params Params) {
	c.params = params
}

func (c *Context) SetError(err error) {
	c.err = err
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

func (c *Context) NoContent() {
	c.ResponseWriter.WriteHeader(204)
}

func (c *Context) ResponseURL() string {
	return c.Request.URL.Path
}

func (c *Context) Method() string {
	return c.Request.Method
}

func (c *Context) Params() Params {
	return c.params
}

// 标准库url包含了ParseQuery，看起来标准库实现的更高效
// 考虑到Params数据结构，暂时先这样实现
// TODO: 优化parseQueries的实现
func (c *Context) parseQueries() {
	rawQuery := c.Request.URL.RawQuery
	c.queries = NewParams()
	if len(rawQuery) == 0 {
		return
	}

	queriesItems := strings.Split(rawQuery, "&")
	for _, queryItem := range queriesItems {
		queryKV := strings.Split(queryItem, "=")
		l := len(queryKV)
		if l == 0 {
			continue
		}
		key, err := url.QueryUnescape(queryKV[0])
		if err != nil {
			continue
		}
		val := ""
		if l == 2 {
			val, err = url.QueryUnescape(queryKV[1])
			if err != nil {
				continue
			}
		}
		c.queries.AddValue(key, val)
	}
}

func (c *Context) Queries() Params {
	return c.queries
}