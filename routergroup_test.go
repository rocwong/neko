package neko

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Router(t *testing.T) {
	testRouteOK("GET", t)
	testRouteOK("POST", t)
	testRouteOK("DELETE", t)
	testRouteOK("PATCH", t)
	testRouteOK("PUT", t)
	testRouteOK("OPTIONS", t)
	testRouteOK("HEAD", t)
	testGroup(t)
	testStatic(t)
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
func testRouteOK(method string, t *testing.T) {
	Convey(method+" Method", t, func() {
		passed := false
		m := New()
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
		w := performRequest(m, method, "/")

		So(passed, ShouldBeTrue)
		So(w.Code, ShouldEqual, http.StatusOK)
	})
}

func testGroup(t *testing.T) {
	Convey("Group routing", t, func() {
		passedGroup, passedGroup2, passedNest := false, false, false
		m := New()
		v1 := m.Group("/v1", func(router *RouterGroup) {
			router.GET("/test", func(ctx *Context) { passedGroup = true })
			router.Group("/sub", func(sub *RouterGroup) {
				sub.GET("/test", func(ctx *Context) { passedNest = true })
			})
		})
		v1.GET("/", func(ctx *Context) { passedGroup2 = true })

		performRequest(m, "GET", "/v1/test")
		So(passedGroup, ShouldBeTrue)

		performRequest(m, "GET", "/v1/")
		So(passedGroup2, ShouldBeTrue)

		performRequest(m, "GET", "/v1/sub/test")
		So(passedNest, ShouldBeTrue)

		w := performRequest(m, "GET", "/v2/test")
		So(w.Code, ShouldEqual, http.StatusNotFound)
	})
}
func testStatic(t *testing.T) {
	Convey("Static serves", t, func() {
		m := New()
		So(func() { m.Static("", "test/") }, ShouldPanic)
		m.Static("/static", "test/")
		w := performRequest(m, "GET", "/static/test.css")
		So(w.Code, ShouldEqual, http.StatusOK)
		w = performRequest(m, "GET", "/static/test1.css")
		So(w.Code, ShouldEqual, http.StatusNotFound)
	})
}
