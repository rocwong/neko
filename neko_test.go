package neko

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func Test_New(t *testing.T) {
	Convey("Initialize a new instance", t, func() {
		So(New(), ShouldNotBeNil)
	})
	Convey("Initialize a ‘classic’ instance", t, func() {
		So(Classic("Neko"), ShouldNotBeNil)
	})
}

func Test_Run(t *testing.T) {
	Convey("Just test that Run doesn't bomb", t, func() {
		go New().Run(":3000")
	})
}

func Test_Handlers(t *testing.T) {
	Convey("Add custom handlers", t, func() {
		result := ""
		m := New()
		m.Use(func(ctx *Context) {
			result += "1"
			ctx.Writer.Before(func(w ResponseWriter) {
				result += "3"
			})
			ctx.SetHeader("x-before", "before")
			ctx.Next()

			result += "5"
			So(ctx.Writer.Size(), ShouldEqual, 5)
		})
		m.Use(func(ctx *Context) {
			ctx.Abort()
			result += "2"
			ctx.Text("abort", http.StatusBadRequest)
			result += "4"
		})
		m.GET("/", func(ctx *Context) {
			result += "0"
			ctx.Text("not found", http.StatusNotFound)
		})

		w := performRequest(m, "GET", "/")
		So(w.HeaderMap.Get("x-before"), ShouldEqual, "before")
		So(w.Code, ShouldEqual, http.StatusBadRequest)
		So(w.Body.String(), ShouldEqual, "abort")
		So(result, ShouldEqual, "12345")
	})
}

func Test_Version(t *testing.T) {
	Convey("Get version", t, func() {
		So(Version(), ShouldNotBeEmpty)
	})
}
