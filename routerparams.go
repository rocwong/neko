package neko

import (
	"github.com/julienschmidt/httprouter"
	"mime/multipart"
	"net/http"
)

type routerParams struct {
	req    *http.Request
	params httprouter.Params
}

func (c *routerParams) ByGet(name string) string {
	val := c.params.ByName(name)
	if val == "" {
		val = c.req.URL.Query().Get(name)
	}
	return val
}

func (c *routerParams) ByPost(name string) string {
	return c.req.FormValue(name)
}

func (c *routerParams) File(name string) (multipart.File, *multipart.FileHeader, error) {
	return c.req.FormFile(name)
}
