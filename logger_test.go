package neko

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func Test_Logger(t *testing.T) {
	testLogger("GET", t)
	testLogger("POST", t)
	testLogger("DELETE", t)
	testLogger("PATCH", t)
	testLogger("PUT", t)
	testLogger("OPTIONS", t)
	testLogger("HEAD", t)

	m := New()
	m.Use(Logger())
	m.GET("/30x", func(ctx *Context) {
		ctx.Redirect("/")
	})
	m.GET("/5xx", func(ctx *Context) {
		ctx.Text("Bad Gateway", http.StatusBadGateway)
	})
	Convey("HTTP 30x", t, func() {
		w := performRequest(m, "GET", "/30x", "")
		So(w.Code, ShouldEqual, http.StatusFound)
	})
	Convey("HTTP 4xx", t, func() {
		w := performRequest(m, "GET", "/404", "")
		So(w.Code, ShouldEqual, http.StatusNotFound)
	})
	Convey("HTTP 5xx", t, func() {
		w := performRequest(m, "GET", "/5xx", "")
		So(w.Code, ShouldEqual, http.StatusBadGateway)
	})
}

func testLogger(method string, t *testing.T) {
	Convey(method+" Method Logger", t, func() {
		passed := false
		m := New()
		m.Use(Logger())
		switch method {
		case "GET":
			m.GET("", func(ctx *Context) { passed = true })
		case "POST":
			m.POST("", func(ctx *Context) { passed = true })
		case "DELETE":
			m.DELETE("", func(ctx *Context) { passed = true })
		case "PATCH":
			m.PATCH("", func(ctx *Context) { passed = true })
		case "PUT":
			m.PUT("", func(ctx *Context) { passed = true })
		case "OPTIONS":
			m.OPTIONS("", func(ctx *Context) { passed = true })
		case "HEAD":
			m.HEAD("", func(ctx *Context) { passed = true })
		}
		// RUN
		w := performRequest(m, method, "/", "")

		So(passed, ShouldBeTrue)
		So(w.Code, ShouldEqual, http.StatusOK)
	})
}
