package neko

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func Test_Recovery(t *testing.T) {
	Convey("Recovery from panic", t, func() {
		m := New()
		m.Use(Recovery())
		m.Use(func(ctx *Context) {
			panic("here is a panic!")
		})
		So(func() { performRequest(m, "GET", "/") }, ShouldNotPanic)
		w := performRequest(m, "GET", "/")
		So(w.Code, ShouldEqual, http.StatusInternalServerError)
		So(w.HeaderMap.Get("Content-Type"), ShouldEqual, "text/html")
	})
}
