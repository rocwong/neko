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
	Convey("Get Params By Get Method", t, func() {
		performRequest(m, "GET", "/params/neko?say=hello&name=golang", "")
	})
	Convey("Get Params By Post Method", t, func() {
		performRequest(m, "POST", "/params/neko", "say=hello&name=golang")
	})
}
