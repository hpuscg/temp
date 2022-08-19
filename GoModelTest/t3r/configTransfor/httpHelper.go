package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var (
	httpClient *http.Client
)

func initHttp() {

	keepAliveTimeout := 60 * time.Second
	timeout := 10 * time.Second

	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: keepAliveTimeout,
		}).Dial,
		// Proxy:               proxy,
		MaxIdleConns:        5,
		MaxIdleConnsPerHost: 2,
	}
	httpClient = &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
}

type HttpRequest struct {
	method   string
	url      string
	headers  map[string]string
	params   map[string]interface{}
	body     string
	fileData io.Reader
}

func NewHttpRequest() *HttpRequest {
	return &HttpRequest{}
}

func (h *HttpRequest) SetMethod(method string) *HttpRequest {
	h.method = method
	return h
}

func (h *HttpRequest) GetMethod() string {
	return h.method
}

func (h *HttpRequest) SetURL(url string) *HttpRequest {
	h.url = url
	return h
}

func (h *HttpRequest) GetURL() string {
	return h.url
}

func (h *HttpRequest) SetHeader(key, value string) *HttpRequest {
	if h.headers == nil {
		h.headers = make(map[string]string)
	}
	h.headers[key] = value
	return h
}

func (h *HttpRequest) GetHeader() map[string]string {
	return h.headers
}

func (h *HttpRequest) SetToken(token string) *HttpRequest {
	return h.SetHeader("Authorization", token)
}

func (h *HttpRequest) GetToken() string {
	return h.headers["Authorization"]
}

func (h *HttpRequest) SetParam(key string, value interface{}) *HttpRequest {
	if h.params == nil {
		h.params = make(map[string]interface{})
	}
	h.params[key] = value
	return h
}

func (h *HttpRequest) GetParam() map[string]interface{} {
	return h.params
}

func (h *HttpRequest) SetEmptyParam() *HttpRequest {
	h.params = make(map[string]interface{})
	return h
}

func (h *HttpRequest) SetContentType(contentType string) *HttpRequest {
	h.SetHeader("Content-Type", contentType)
	return h
}

func (h *HttpRequest) SetFileData(data io.Reader) *HttpRequest {
	h.fileData = data
	return h
}

func (h *HttpRequest) DoWithFile() (int, []byte, error) {
	if h.fileData == nil {
		return 0, nil, errors.New("no file data")
	}
	// 创建请求
	req, err := http.NewRequest(h.method, h.url, h.fileData)
	if err != nil {
		return 0, nil, err
	}
	// 添加Header
	if h.headers != nil {
		for k, v := range h.headers {
			req.Header.Add(k, v)
		}
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, nil, err
	}
	return res.StatusCode, content, nil
}

func (h *HttpRequest) SetBody(body string) *HttpRequest {
	h.body = body
	return h
}

func (h *HttpRequest) DoWithBody() (int, []byte, error) {
	if h.body == "" {
		return 0, nil, errors.New("no body data")
	}
	// 创建请求
	req, err := http.NewRequest(h.method, h.url, strings.NewReader(h.body))
	if err != nil {
		return 0, nil, err
	}
	// 添加Header
	if h.headers != nil {
		for k, v := range h.headers {
			req.Header.Add(k, v)
		}
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, nil, err
	}
	return res.StatusCode, content, nil
}

func (h *HttpRequest) Do() (int, []byte, error) {
	var requestData io.Reader
	// 拼接参数
	if h.params != nil {
		if data, err := jsoniter.Marshal(h.params); err != nil {
			return 0, nil, err
		} else {
			requestData = bytes.NewReader(data)
		}
	}
	// 创建请求
	req, err := http.NewRequest(h.method, h.url, requestData)
	if err != nil {
		return 0, nil, err
	}
	// 添加Header
	if h.headers != nil {
		for k, v := range h.headers {
			req.Header.Add(k, v)
		}
	}
	if h.params != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, nil, err
	}
	return res.StatusCode, content, nil
}
