package render

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

type (
	Render interface {
		Render(interface{}, http.ResponseWriter) error
	}
	JSON  struct{}
	JSONP struct {
		Callback string
	}
	XML  struct{}
	TEXT struct{}
)

const (
	// ContentType header constant.
	ContentType = "Content-Type"
	// ContentJSON header value for JSON data.
	ContentJSON = "application/json"
	// ContentJSONP header value for JSONP data.
	ContentJSONP = "application/javascript"
	// ContentXML header value for XML data.
	ContentXML = "text/xml"
	// ContentPlain header value for Text data.
	ContentPlain = "text/plain"
)

// Render an JSON response.
func (c JSON) Render(data interface{}, w http.ResponseWriter) error {
	result, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set(ContentType, ContentJSON)
	w.Write(result)
	return nil
}

// Render an JSONP response.
func (c JSONP) Render(data interface{}, w http.ResponseWriter) error {
	result, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set(ContentType, ContentJSONP)
	w.Write([]byte(c.Callback + "("))
	w.Write(result)
	w.Write([]byte(");"))
	return nil
}

// Render an XML response.
func (c XML) Render(data interface{}, w http.ResponseWriter) error {
	result, err := xml.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set(ContentType, ContentXML)
	w.Write(result)
	return nil
}

// Render an Text response.
func (c TEXT) Render(data interface{}, w http.ResponseWriter) error {
	w.Header().Set(ContentType, ContentPlain)
	_, err := w.Write([]byte(data.(string)))
	return err
}
