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

	w := performRequest(m, "GET", "/json")
	Convey("Json Render", t, func() {
		So(w.Code, ShouldEqual, http.StatusOK)
		So(w.Header().Get(render.ContentType), ShouldEqual, render.ContentJSON)
	})

	w = performRequest(m, "GET", "/jsonp")
	Convey("Jsonp Render", t, func() {
		So(w.Code, ShouldEqual, http.StatusOK)
		So(w.Header().Get(render.ContentType), ShouldEqual, render.ContentJSONP)
	})

	w = performRequest(m, "GET", "/xml")
	Convey("Xml Render", t, func() {
		So(w.Code, ShouldEqual, http.StatusOK)
		So(w.Header().Get(render.ContentType), ShouldEqual, render.ContentXML)
	})

	w = performRequest(m, "GET", "/text")
	Convey("Text Render", t, func() {
		So(w.Code, ShouldEqual, http.StatusNotFound)
		So(w.Header().Get(render.ContentType), ShouldEqual, render.ContentPlain)
	})

	w = performRequest(m, "GET", "/redirect/302")
	Convey("Redirect 302", t, func() {
		So(w.Code, ShouldEqual, http.StatusFound)
	})

	w = performRequest(m, "GET", "/redirect/301")
	Convey("Redirect 301", t, func() {
		So(w.Code, ShouldEqual, http.StatusMovedPermanently)
	})
}

func Test_Version(t *testing.T) {
	Convey("Get version", t, func() {
		So(Version(), ShouldNotBeEmpty)
	})
}
