package neko

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

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

// SetCookie sets given cookie value to response header.
// ctx.SetCookie(name, value [, MaxAge, Path, Domain, Secure, HttpOnly])
func (c *Context) SetCookie(name, value string, others ...interface{}) {
	cookie := &http.Cookie{}
	cookie.Name = name
	cookie.Value = value

	if len(others) > 0 {
		switch v := others[0].(type) {
		case int:
			cookie.MaxAge = v
		case int64:
			cookie.MaxAge = int(v)
		case int32:
			cookie.MaxAge = int(v)
		}
	}

	// default "/"
	if len(others) > 1 {
		if v, ok := others[1].(string); ok && len(v) > 0 {
			cookie.Path = v
		}
	} else {
		cookie.Path = "/"
	}

	// default empty
	if len(others) > 2 {
		if v, ok := others[2].(string); ok && len(v) > 0 {
			cookie.Domain = v
		}
	}

	// default empty
	if len(others) > 3 {
		switch v := others[3].(type) {
		case bool:
			cookie.Secure = v
		}
	}

	// default false.
	if len(others) > 4 {
		if v, ok := others[4].(bool); ok && v {
			cookie.HttpOnly = true
		}
	}

	http.SetCookie(c.Writer, cookie)
}

// GetCookie returns given cookie value from request header.
func (c *Context) GetCookie(name string) string {
	cookie, err := c.Req.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

var cookieSecret string

// SetCookieSecret sets global default secure cookie secret.
func (m *Engine) SetCookieSecret(secret string) {
	cookieSecret = secret
}

// SetSecureCookie sets given cookie value to response header with default secret string.
func (ctx *Context) SetSecureCookie(name, value string, others ...interface{}) {
	ctx.SetBasicSecureCookie(cookieSecret, name, value, others...)
}

// GetSecureCookie returns given cookie value from request header with default secret string.
func (ctx *Context) GetSecureCookie(name string) (string, bool) {
	return ctx.GetBasicSecureCookie(cookieSecret, name)
}

// SetBasicSecureCookie sets given cookie value to response header with secret string.
func (ctx *Context) SetBasicSecureCookie(Secret, name, value string, others ...interface{}) {

	vs := base64.URLEncoding.EncodeToString([]byte(value))
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	hm := hmac.New(sha1.New, []byte(Secret))
	fmt.Fprintf(hm, "%s%s", vs, timestamp)
	sig := fmt.Sprintf("%02x", hm.Sum(nil))
	cookie := strings.Join([]string{vs, timestamp, sig}, "|")
	ctx.SetCookie(name, cookie, others...)
}

// GetBasicSecureCookie returns given cookie value from request header with secret string.
func (ctx *Context) GetBasicSecureCookie(Secret, name string) (string, bool) {
	val := ctx.GetCookie(name)
	if val == "" {
		return "", false
	}

	parts := strings.SplitN(val, "|", 3)
	if len(parts) != 3 {
		return "", false
	}
	vs := parts[0]
	timestamp := parts[1]
	sig := parts[2]

	hm := hmac.New(sha1.New, []byte(Secret))
	fmt.Fprintf(hm, "%s%s", vs, timestamp)
	if fmt.Sprintf("%02x", hm.Sum(nil)) != sig {
		return "", false
	}
	res, _ := base64.URLEncoding.DecodeString(vs)
	return string(res), true
}
