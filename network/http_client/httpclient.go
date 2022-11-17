package http_client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"example.com/login/network"
)

var debugNetwork bool = false

type httpclient struct {
	client *http.Client
	file   io.Writer
}

func New(client *http.Client, f io.Writer) *httpclient {
	if os.Getenv("DEBUG_NETWORK") == "1" {
		debugNetwork = true
	}
	return &httpclient{
		client: client,
		file:   f,
	}
}

func (hc *httpclient) PostFormData(body interface{}, uri string) network.Response {
	payloads, ok := body.(map[string]interface{})
	if !ok {
		return network.Response{Response: nil, Err: fmt.Errorf("cannot cast payloads to []map[string]interface{}")}
	}
	data := url.Values{}
	for key, value := range payloads {
		data.Add(key, fmt.Sprintf("%v", value))
	}
	httpheaders := map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
	}
	reqURL, err := url.Parse(uri)
	if err != nil {
		return network.Response{
			Err: err,
		}
	}
	req := &http.Request{
		Method: http.MethodPost,
		URL:    reqURL,
		Header: httpheaders,
		Body:   ioutil.NopCloser(strings.NewReader(data.Encode())),
	}
	return do(req, hc)
}

func (hc *httpclient) PostFormDataBytes(body interface{}, uri string, headers map[string]string) network.ResponseBytes {
	payloads, ok := body.(map[string]interface{})
	if !ok {
		return network.ResponseBytes{Response: nil, Err: fmt.Errorf("cannot cast payloads to []map[string]interface{}")}
	}

	data := url.Values{}
	for key, value := range payloads {
		_, ok := value.(string)
		if !ok {
			vbb, err := json.Marshal(value)
			if err != nil {
				return network.ResponseBytes{
					Err: err,
				}
			}
			value = string(vbb)
		}
		data.Add(key, fmt.Sprintf("%v", value))
	}

	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(data.Encode()))
	if err != nil {
		return network.ResponseBytes{
			Err: err,
		}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for key, val := range headers {
		req.Header.Set(key, val)
	}

	return doBytes(req, hc)
}

func (hc *httpclient) Post(body interface{}, uri string) network.Response {
	var bodyMessage []byte
	bodyMessage, _ = json.Marshal(body)
	httpheaders := map[string][]string{
		"Content-Type": {"application/json"},
	}
	// for i := range headers {
	// 	httpheaders[i] = []string{headers[i]}
	// }
	reqURL, err := url.Parse(uri)
	if err != nil {
		return network.Response{
			Err: err,
		}
	}
	req := &http.Request{
		Method: http.MethodPost,
		URL:    reqURL,
		Header: httpheaders,
		Body:   ioutil.NopCloser(bytes.NewReader(bodyMessage)),
	}
	return do(req, hc)
}

func (hc *httpclient) GetBytes(uri string, headers map[string]string) network.ResponseBytes {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return network.ResponseBytes{
			Err: err,
		}
	}
	req.Header.Set("Content-Type", "application/json")
	for key, val := range headers {
		req.Header.Set(key, val)
	}

	return doBytes(req, hc)
}

func (hc *httpclient) Get(uri string, headers map[string]string) network.Response {
	httpheaders := map[string][]string{
		"Content-Type": {"application/json"},
	}
	for i := range headers {
		httpheaders[i] = []string{headers[i]}
	}
	reqURL, err := url.Parse(uri)
	if err != nil {
		return network.Response{
			Err: err,
		}
	}
	req := &http.Request{
		Method: http.MethodGet,
		URL:    reqURL,
		Header: httpheaders,
	}
	return do(req, hc)
}
func (hc *httpclient) Delete(uri string, body interface{}, headers map[string]string) network.Response {
	bodyMessage, err := json.Marshal(body)
	if err != nil {
		return network.Response{
			Err: err,
		}
	}
	httpheaders := map[string][]string{
		"Content-Type": {"application/json"},
	}
	for i := range headers {
		httpheaders[i] = []string{headers[i]}
	}
	reqURL, err := url.Parse(uri)
	if err != nil {
		return network.Response{
			Err: err,
		}
	}
	req := &http.Request{
		Method: http.MethodDelete,
		URL:    reqURL,
		Header: httpheaders,
		Body:   ioutil.NopCloser(bytes.NewReader(bodyMessage)),
	}
	return do(req, hc)
}

func (hc *httpclient) PostJSON(body interface{}, uri string, headers map[string]string) network.Response {
	bodyMessage, err := json.Marshal(body)
	if err != nil {
		return network.Response{
			Err: err,
		}
	}
	httpheaders := map[string][]string{
		"Content-Type": {"application/json"},
	}
	for i := range headers {
		httpheaders[i] = []string{headers[i]}
	}
	reqURL, err := url.Parse(uri)
	if err != nil {
		return network.Response{
			Err: err,
		}
	}
	req := &http.Request{
		Method: http.MethodPost,
		URL:    reqURL,
		Header: httpheaders,
		Body:   ioutil.NopCloser(bytes.NewReader(bodyMessage)),
	}
	return do(req, hc)
}

func (hc *httpclient) PostJSONBytes(body interface{}, uri string, headers map[string]string) network.ResponseBytes {
	bodyMessage, err := json.Marshal(body)
	if err != nil {
		return network.ResponseBytes{
			Err: err,
		}
	}
	httpheaders := map[string][]string{}
	for i := range headers {
		httpheaders[i] = []string{headers[i]}
	}
	reqURL, err := url.Parse(uri)
	if err != nil {
		return network.ResponseBytes{
			Err: err,
		}
	}
	req := &http.Request{
		Method: http.MethodPost,
		URL:    reqURL,
		Header: httpheaders,
		Body:   ioutil.NopCloser(bytes.NewReader(bodyMessage)),
	}
	return doBytes(req, hc)
}
func (hc *httpclient) PutJSONBytes(body interface{}, uri string, headers map[string]string) network.ResponseBytes {
	bodyMessage, err := json.Marshal(body)
	if err != nil {
		return network.ResponseBytes{
			Err: err,
		}
	}
	httpheaders := map[string][]string{
		"Content-Type": {"application/json"},
	}
	for i := range headers {
		httpheaders[i] = []string{headers[i]}
	}
	reqURL, err := url.Parse(uri)
	if err != nil {
		return network.ResponseBytes{
			Err: err,
		}
	}
	req := &http.Request{
		Method: http.MethodPut,
		URL:    reqURL,
		Header: httpheaders,
		Body:   ioutil.NopCloser(bytes.NewReader(bodyMessage)),
	}
	return doBytes(req, hc)
}

func do(req *http.Request, hc *httpclient) network.Response {
	resp, err := hc.client.Do(req)
	if err != nil {
		return network.Response{
			Err: err,
		}
	}

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return network.Response{
				Err: err,
			}
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	bodyMessage, err := ioutil.ReadAll(reader)
	if err != nil {
		return network.Response{
			Err: err,
		}
	}
	resp.Body.Close()
	response := map[string]interface{}{}
	err = json.Unmarshal(bodyMessage, &response)
	if err != nil {
		return network.Response{
			Err: err,
		}
	}

	if resp.StatusCode < 200 && resp.StatusCode > 300 {
		return network.Response{Response: response, Err: fmt.Errorf("%s", string(bodyMessage))}
	}
	return network.Response{Response: response, Err: nil}
}
func doBytes(req *http.Request, hc *httpclient) network.ResponseBytes {
	dumpReq, _ := httputil.DumpRequestOut(req, true)
	resp, err := hc.client.Do(req)
	if err != nil {
		return network.ResponseBytes{
			Err: err,
		}
	}
	dumpRes, _ := httputil.DumpResponse(resp, false)
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return network.ResponseBytes{
				Err: err,
			}
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	bodyMessage, err := ioutil.ReadAll(reader)
	if err != nil {
		return network.ResponseBytes{
			Err: err,
		}
	}
	resp.Body.Close()

	hc.file.Write(dumpReq)
	hc.file.Write(dumpRes)
	hc.file.Write(bodyMessage)

	if resp.StatusCode < 200 && resp.StatusCode > 300 {
		return network.ResponseBytes{Response: bodyMessage, Err: fmt.Errorf("%s", string(bodyMessage))}
	}
	return network.ResponseBytes{Response: bodyMessage, Err: nil}
}
