package neko

import (
	"encoding/json"
	"encoding/xml"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Context struct {
	Writer   ResponseWriter
	Req      *http.Request
	Cookies  Cookie
	Session  Session
	Params   httprouter.Params
	Engine   *Engine
	writer   writer
	handlers []HandlerFunc
	index    int8
	HtmlEngine
}

// Next should be used only in the middlewares.
// It executes the pending handlers in the chain inside the calling handler.
func (c *Context) Next() {
	c.index++
	s := int8(len(c.handlers))
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// SetHeader sets a response header.
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

// Redirect returns a HTTP redirect to the specific location. default for 302
func (c *Context) Redirect(location string, status ...int) {
	c.SetHeader("Location", location)
	if status != nil {
		http.Redirect(c.Writer, c.Req, location, status[0])
	} else {
		http.Redirect(c.Writer, c.Req, location, 302)
	}
}

// Serializes the given struct as JSON into the response body in a fast and efficient way.
// It also sets the Content-Type as "application/json".
func (c *Context) Json(data interface{}) error {
	c.SetHeader("Content-Type", "application/json")
	encoder := json.NewEncoder(c.Writer)
	return encoder.Encode(data)
}

// Serializes the given struct as XML into the response body in a fast and efficient way.
// It also sets the Content-Type as "application/xml".
func (c *Context) Xml(data interface{}) error {
	c.SetHeader("Content-Type", "application/xml")
	encoder := xml.NewEncoder(c.Writer)
	return encoder.Encode(data)
}

// Writes the given string into the response body and sets the Content-Type to "text/plain".
func (c *Context) Text(data string) {
	c.SetHeader("Content-Type", "text/plain")
	c.Writer.Write([]byte(data))
}

// ClientIP returns more real IP address.
func (c *Context) ClientIP() string {
	clientIP := c.Req.Header.Get("X-Real-IP")
	if len(clientIP) == 0 {
		clientIP = c.Req.Header.Get("X-Forwarded-For")
	}
	if len(clientIP) == 0 {
		clientIP = c.Req.RemoteAddr
	}
	return clientIP
}

func (c *Engine) createContext(w http.ResponseWriter, req *http.Request, params httprouter.Params, handlers []HandlerFunc) *Context {
	ctx := c.pool.Get().(*Context)
	ctx.Writer = &ctx.writer
	ctx.Req = req
	ctx.Cookies = GetCookies(w, req)
	ctx.Params = params
	ctx.handlers = handlers
	ctx.writer.reset(w)
	ctx.index = -1
	return ctx
}
