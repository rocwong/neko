package neko

import (
	"net/http"
)

type Cookie interface {
	// Set sets given cookie value to response header.
	Set(name, value string, other ...interface{})
	// Get returns given cookie value from request header.
	Get(name string) string
	// GetCookies returns http.Cookie object from request header.
	GetCookies(name string) *http.Cookie
	// Delete removes given cookie value.
	Delete(name string)
	// Clear deletes all values in the cookie.
	Clear()
}

func GetCookies(res http.ResponseWriter, req *http.Request) Cookie {
	return &cookie{
		res: res,
		req: req,
	}
}

type cookie struct {
	res http.ResponseWriter
	req *http.Request
}

func (c *cookie) Set(name, value string, other ...interface{}) {
	cookie := http.Cookie{}
	cookie.Name = name
	cookie.Value = value

	c.res.Header().Set("Set-Cookie", cookie.String())
}

func (c *cookie) Get(name string) string {
	return c.GetCookies(name).Value
}

func (c *cookie) GetCookies(name string) *http.Cookie {
	cookie, err := c.req.Cookie(name)
	if err != nil {
		return &http.Cookie{}
	}
	return cookie
}

func (c *cookie) Delete(name string) {

}

func (c *cookie) Clear() {

}
