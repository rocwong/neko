package neko

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_RouterParams(t *testing.T) {
	m := New()
	m.GET("/params/:name", func(ctx *Context) {
		So(ctx.Params.ByGet("name"), ShouldEqual, "neko")
		So(ctx.Params.ByGet("say"), ShouldEqual, "hello")
	})
	m.POST("/params/:name", func(ctx *Context) {
		So(ctx.Params.ByGet("name"), ShouldEqual, "neko")
		So(ctx.Params.ByPost("say"), ShouldEqual, "hello")
	})

	m.POST("/json/:name", func(ctx *Context) {
		dataJson := ctx.Params.Json()
		So(ctx.Params.ByGet("name"), ShouldEqual, "neko")
		So(dataJson.Get("say"), ShouldEqual, "hello")
		So(dataJson.String(), ShouldEqual, `{"say": "hello"}`)
	})

	m.POST("/json-empty", func(ctx *Context) {
		dataJson := ctx.Params.Json()
		So(dataJson.String(), ShouldEqual, "")
		So(dataJson.Get("empty"), ShouldEqual, "")
	})

	Convey("Get Params By Query String", t, func() {
		performRequest(m, "GET", "/params/neko?say=hello&name=golang", "")
	})
	Convey("Get Params By Form Post", t, func() {
		performRequest(m, "POST", "/params/neko", "say=hello&name=golang")
	})
	Convey("Get Params By Json Data", t, func() {
		performRequest(m, "POST|JSON", "/json/neko", `{"say": "hello"}`)
		performRequest(m, "POST|JSON", "/json-empty", "")
	})
}
