package network

import (
	"encoding/json"
)

type Client interface {
	Get(uri string, headers map[string]string) Response
	GetBytes(uri string, headers map[string]string) ResponseBytes
	Delete(uri string, body interface{}, headers map[string]string) Response
	Post(body interface{}, url string) Response
	PostJSON(body interface{}, url string, headers map[string]string) Response
	PostJSONBytes(body interface{}, url string, headers map[string]string) ResponseBytes
	PostFormData(body interface{}, url string) Response
	PostFormDataBytes(body interface{}, url string, headers map[string]string) ResponseBytes
	PutJSONBytes(body interface{}, uri string, headers map[string]string) ResponseBytes
	// PostForm(body []byte, url string) Response
	// PostFormDataBytes(body interface{}, url string, headers map[string]string) ResponseBytes
	// Redirect(req *fasthttp.Request, url string) Response
}

type Response struct {
	Response map[string]interface{}
	Err      error
}

func (r *Response) String() string {
	bb, _ := json.Marshal(r.Response)
	return string(bb)
}

type ResponseBytes struct {
	Response []byte
	Err      error
}

func (r *ResponseBytes) String() string {
	return string(r.Response)
}
