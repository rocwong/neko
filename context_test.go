package neko

import (
	"encoding/xml"
	"github.com/rocwong/neko/render"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

type ExampleXml struct {
	XMLName xml.Name `xml:"example"`
	One     string   `xml:"one,attr"`
	Two     string   `xml:"two,attr"`
}

func Test_Render(t *testing.T) {
	m := New()
	m.GET("/json", func(ctx *Context) {
		ctx.Json(JSON{"test": "json render"})
	})
	m.GET("/jsonp", func(ctx *Context) {
		ctx.Jsonp("callback", JSON{"test": "json render"})
	})
	m.GET("/xml", func(ctx *Context) {
		ctx.Xml(ExampleXml{One: "hello", Two: "xml"})
	})
	m.GET("/text", func(ctx *Context) {
		ctx.Text("not found", 404)
	})
	m.GET("/redirect/302", func(ctx *Context) {
		ctx.Redirect("/")
	})
	m.GET("/redirect/301", func(ctx *Context) {
		ctx.Redirect("/", 301)
	})

	w := performRequest(m, "GET", "/json", "")
	Convey("Json Render", t, func() {
		So(w.Code, ShouldEqual, http.StatusOK)
		So(w.Header().Get(render.ContentType), ShouldEqual, render.ContentJSON)
	})

	w = performRequest(m, "GET", "/jsonp", "")
	Convey("Jsonp Render", t, func() {
		So(w.Code, ShouldEqual, http.StatusOK)
		So(w.Header().Get(render.ContentType), ShouldEqual, render.ContentJSONP)
	})

	w = performRequest(m, "GET", "/xml", "")
	Convey("Xml Render", t, func() {
		So(w.Code, ShouldEqual, http.StatusOK)
		So(w.Header().Get(render.ContentType), ShouldEqual, render.ContentXML)
	})

	w = performRequest(m, "GET", "/text", "")
	Convey("Text Render", t, func() {
		So(w.Code, ShouldEqual, http.StatusNotFound)
		So(w.Header().Get(render.ContentType), ShouldEqual, render.ContentPlain)
	})

	w = performRequest(m, "GET", "/redirect/302", "")
	Convey("Redirect 302", t, func() {
		So(w.Code, ShouldEqual, http.StatusFound)
	})

	w = performRequest(m, "GET", "/redirect/301", "")
	Convey("Redirect 301", t, func() {
		So(w.Code, ShouldEqual, http.StatusMovedPermanently)
	})
}

func Test_GetSet(t *testing.T) {
	Convey("Get/Set Method", t, func() {
		m := New()
		m.GET("/test", func(ctx *Context) {
			So(ctx.Keys, ShouldBeNil)

			ctx.Set("foo", "bar")

			v, err := ctx.Get("foo")
			So(err, ShouldBeNil)
			So(v, ShouldEqual, "bar")

			v, err = ctx.Get("foo2")
			So(err, ShouldNotBeNil)

			v = ctx.MustGet("foo")
			So(v, ShouldEqual, "bar")

			So(func() { ctx.MustGet("foo3") }, ShouldPanic)
		})
		// First Visit
		performRequest(m, "GET", "/test", "")
		// The Other Visit
		performRequest(m, "GET", "/test", "")
	})
}
