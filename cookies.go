package neko

import (
	"net/http"
)

type Cookie interface {
	// Set sets given cookie value to response header.
	Set(name, value string, other ...interface{})
	// Get returns given cookie value from request header.
	Get(name string) string
	// GetCookie returns http.Cookie object from request header.
	GetCookie(name string) *http.Cookie
	// Delete removes given cookie value.
	Delete(name string)
	// Clear deletes all values in the cookie.
	Clear()
}
