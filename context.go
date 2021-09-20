package simplegin

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HandlersChain []HandlerFunc

type Context struct {
	Writer   *responseWriter
	Req      *http.Request
	Path     string
	Method   string
	Params   map[string]string
	handlers HandlersChain // [middlewares1, middleware2, ..., handler]
	index    int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: &responseWriter{
			ResponseWriter: w,
		},
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (c *Context) Next() {
	for c.index++; c.index < len(c.handlers); c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, value ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)

	if _, err := c.Writer.Write([]byte(fmt.Sprintf(format, value...))); err != nil {
		panic(err)
	}
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)

	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		panic(err)
	}
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)

	if _, err := c.Writer.Write([]byte(html)); err != nil {
		panic(err)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)

	if _, err := c.Writer.Write(data); err != nil {
		panic(err)
	}
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}
