package neko

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_Cookie(t *testing.T) {
	Convey("Set and get cookie", t, func() {
		m := New()

		testSet(m)
		testGet(m)
	})
}

func Test_SecureCookie(t *testing.T) {
	Convey("Set and get secure cookie", t, func() {
		m := New()

		testSetSecureCookie(m)
		testGetSecureCookie(m)
	})
}

func testSet(m *Engine) {
	m.GET("/", func(ctx *Context) {
		ctx.SetCookie("basic", "neko")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	m.ServeHTTP(w, req)
	So(w.Header().Get("Set-Cookie"), ShouldEqual, "basic=neko; Path=/")
}

func testGet(m *Engine) {
	m.GET("/get", func(ctx *Context) {
		name := ctx.GetCookie("name")
		age := ctx.GetCookie("age")
		job := ctx.GetCookie("job")
		ctx.Text(name + ":" + age + job)
	})
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/get", nil)
	req.Header.Set("Cookie", "name=neko; age=3")
	m.ServeHTTP(w, req)
	So(w.Body.String(), ShouldEqual, "neko:3")
}

func testSetSecureCookie(m *Engine) {
	m.SetCookieSecret("secret123")
	m.GET("/set-secure", func(ctx *Context) {
		ctx.SetSecureCookie("full", "neko", 86400000000000, "/full", "abc.com", true, true)
		ctx.SetSecureCookie("full_32", "neko", int32(8640000), "/full", "abc.com", true, true)
		ctx.SetSecureCookie("full_64", "neko", int64(86400000000000), "/full", "abc.com", true, true)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/set-secure", nil)
	m.ServeHTTP(w, req)
	So(strings.HasPrefix(w.Header().Get("Set-Cookie"), "full=bmVrbw==|"), ShouldBeTrue)
}

func testGetSecureCookie(m *Engine) {
	m.GET("/get-secure", func(ctx *Context) {
		user, _ := ctx.GetSecureCookie("user")
		age, _ := ctx.GetSecureCookie("age")
		ctx.Text("hello " + user + age)
	})
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/get-secure", nil)
	req.Header.Set("Cookie", "user=bmVrbw==|1420100804032788408|6852e5511056060c41c991b6b228703c8ecae790; Path=/")
	m.ServeHTTP(w, req)
	So(w.Body.String(), ShouldEqual, "hello neko")
}
