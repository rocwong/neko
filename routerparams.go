package neko

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"mime/multipart"
	"net/http"
	"io/ioutil"
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

func (c *routerParams) Json() *jsonParams {
	defer c.req.Body.Close()

	data, _ := ioutil.ReadAll(c.req.Body)
	objJson := &jsonParams{ data: map[string]string{}}
	objJson.source = string(data)
	json.Unmarshal(data, &objJson.data);

	return objJson
}

type jsonParams struct {
	source string
	data map[string]string
}

func (c *jsonParams) Get(name string) string {
	if len(c.data) == 0 {
		return ""
	}
	return c.data[name]
}

func (c *jsonParams) String() string {
	return c.source
}
