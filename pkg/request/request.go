package request

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Request ...
type Request struct {
	url    string
	method string
	header http.Header

	body   io.Reader
	result interface{}

	err error
}

// New ...
func New() *Request {
	r := &Request{
		header: http.Header{},
	}
	r.header.Set("User-Agent", userAgent)
	r.header.Set("Content-Type", "application/json")
	return r
}

// Method ...
func (a *Request) Method(method string) *Request {
	a.method = method
	return a
}

// Url ...
func (a *Request) Url(url string) *Request {
	a.url = url
	return a
}

// Get ...
func (a *Request) Get(url string) *Request {
	a.method = http.MethodGet
	a.url = url
	return a
}

// Post ...
func (a *Request) Post(url string) *Request {
	a.method = http.MethodPost
	a.url = url
	return a
}

// Put ...
func (a *Request) Put(url string) *Request {
	a.method = http.MethodPut
	a.url = url
	return a
}

// Delete ...
func (a *Request) Delete(url string) *Request {
	a.method = http.MethodDelete
	a.url = url
	return a
}

// RawBody ...
func (a *Request) RawBody(body string) *Request {
	a.body = bytes.NewReader([]byte(body))
	return a
}

// Body ...
func (a *Request) Body(body interface{}) *Request {
	bs, err := json.Marshal(body)
	if err != nil {
		a.err = err
	} else {
		a.body = bytes.NewReader(bs)
	}
	return a
}

// Header ...
func (a *Request) Header(h http.Header) *Request {
	for k, v := range h {
		a.header[k] = v
	}
	return a
}

// SetHeader ...
func (a *Request) SetHeader(key, val string) *Request {
	a.header.Set(key, val)
	return a
}

// Result ...
func (a *Request) Result(result interface{}) *Request {
	a.result = result
	return a
}

// Do ...
func (a *Request) Do() (*Response, error) {
	if a.err != nil {
		return nil, a.err
	}
	req, err := http.NewRequest(a.method, a.url, a.body)
	if err != nil {
		return nil, err
	}
	req.Header = a.header

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := bodyHandler(resp)
	if err != nil {
		return nil, err
	}

	if a.result != nil {
		err = json.Unmarshal(respBody, a.result)
		if err != nil {
			return nil, fmt.Errorf("unmarshal error: %s, data: %s", err.Error(), string(respBody))
		}
	}
	response := &Response{
		Request:    req,
		Response:   resp,
		StatusCode: resp.StatusCode,
		Content:    respBody,
	}
	return response, nil
}

func bodyHandler(resp *http.Response) ([]byte, error) {
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		if reader, err = gzip.NewReader(bytes.NewBuffer(respBody)); err != nil {
			return nil, err
		}
	case "deflate":
		// deflate should be zlib
		// http://www.gzip.org/zlib/zlib_faq.html#faq38
		if reader, err = zlib.NewReader(bytes.NewBuffer(respBody)); err != nil {
			// try RFC 1951 deflate
			// http: //www.open-open.com/lib/view/open1460866410410.html
			reader = flate.NewReader(bytes.NewBuffer(respBody))
		}
	}
	if reader == nil {
		return respBody, nil
	}

	defer reader.Close()
	respBody, err = ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

// Response ...
type Response struct {
	Request    *http.Request
	Response   *http.Response
	StatusCode int
	Content    []byte
}

// OK ...
func (a *Response) OK() bool {
	return a.StatusCode < 300
}

// String deprecated
func (a *Response) String() string {
	return string(a.Content)
}

// Text ...
func (a *Response) Text() string {
	return string(a.Content)
}
