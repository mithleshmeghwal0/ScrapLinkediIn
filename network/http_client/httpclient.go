package http_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"example.com/login/network"
	"github.com/sirupsen/logrus"
)

var debugNetwork bool = false

type httpclient struct {
	client *http.Client
	log    *logrus.Entry
}

func New(client *http.Client, log *logrus.Entry) *httpclient {
	if os.Getenv("DEBUG_NETWORK") == "1" {
		debugNetwork = true
	}
	return &httpclient{
		client: client,
		log:    log,
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
	return hc.do(req)
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
	httpheaders := map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
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
		Method: http.MethodPost,
		URL:    reqURL,
		Header: httpheaders,
		Body:   ioutil.NopCloser(strings.NewReader(data.Encode())),
	}
	return hc.doBytes(req)
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
	return hc.do(req)
}

func (hc *httpclient) GetBytes(uri string, headers map[string]string) network.ResponseBytes {
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
		Method: http.MethodGet,
		URL:    reqURL,
		Header: httpheaders,
	}
	return hc.doBytes(req)
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
	return hc.do(req)
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
	return hc.do(req)
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
	return hc.do(req)
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
	return hc.doBytes(req)
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
	return hc.doBytes(req)
}

func (hc *httpclient) do(req *http.Request) network.Response {
	cookies := hc.client.Jar.Cookies(req.URL)
	for i := range cookies {
		req.AddCookie(cookies[i])
	}
	reqdump := printRequest(req)
	resp, err := hc.client.Do(req)
	if err != nil {
		return network.Response{
			Err: err,
		}
	}

	respdump := printResponse(resp)
	hc.printRR(reqdump, respdump)

	bodyMessage, err := ioutil.ReadAll(resp.Body)
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
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return network.Response{Response: response, Err: fmt.Errorf("%s", string(bodyMessage))}
	}
	return network.Response{Response: response, Err: nil}
}
func (hc *httpclient) doBytes(req *http.Request) network.ResponseBytes {
	cookies := hc.client.Jar.Cookies(req.URL)
	for i := range cookies {
		req.AddCookie(cookies[i])
	}
	reqdump := printRequest(req)
	resp, err := hc.client.Do(req)
	if err != nil {
		return network.ResponseBytes{
			Err: err,
		}
	}
	respdump := printResponse(resp)
	hc.printRR(reqdump, respdump)
	bodyMessage, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return network.ResponseBytes{
			Err: err,
		}
	}
	resp.Body.Close()
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return network.ResponseBytes{Response: bodyMessage, Err: fmt.Errorf("%s", string(bodyMessage))}
	}
	return network.ResponseBytes{Response: bodyMessage, Err: nil}
}

func printRequest(req *http.Request) (dump []byte) {
	if debugNetwork {
		dump, _ = httputil.DumpRequestOut(req, true)
	}
	return
}
func printResponse(resp *http.Response) (dump []byte) {
	if debugNetwork {
		dump, _ = httputil.DumpResponse(resp, false)
	}
	return
}

func (hc *httpclient) printRR(reqDump, respDump []byte) {
	if debugNetwork {
		fmt.Println(string(reqDump))
	}
}
