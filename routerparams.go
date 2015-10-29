package neko

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
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

func (c *routerParams) Json() *jsonParams {
	defer c.req.Body.Close()

	data, _ := ioutil.ReadAll(c.req.Body)
	objJson := &jsonParams{data: map[string]interface{}{}}
	objJson.source = string(data)
	json.Unmarshal(data, &objJson.data)

	return objJson
}

type jsonParams struct {
	source string
	data   map[string]interface{}
}

func (c *jsonParams) Get(name string) interface{} {
	if len(c.data) == 0 || c.data[name] == nil {
		return ""
	}
	return c.data[name]
}

func (c *jsonParams) GetString(name string) string {
	if len(c.data) == 0 || c.data[name] == nil {
		return ""
	}
	return toString(c.data[name])
}

func (c *jsonParams) GetInt32(name string) int32 {
	if len(c.data) == 0 || c.data[name] == nil {
		return 0
	}
	return toInt32(c.data[name])
}

func (c *jsonParams) GetUInt32(name string) uint32 {
	if len(c.data) == 0 || c.data[name] == nil {
		return 0
	}
	return toUint32(c.data[name])
}

func (c *jsonParams) GetFloat32(name string) float32 {
	if len(c.data) == 0 || c.data[name] == nil {
		return 0
	}
	return toFloat32(c.data[name])
}

func (c *jsonParams) GetFloat64(name string) float64 {
	if len(c.data) == 0 || c.data[name] == nil {
		return 0
	}
	return toFloat64(c.data[name])
}

func (c *jsonParams) String() string {
	return c.source
}
