package ovopark

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type myClient struct {
	cli *http.Client
}

type Client interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

func NewClient(timeout time.Duration) Client {
	return &myClient{
		cli: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *myClient) Do(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
	req = req.WithContext(ctx)
	resp, err = c.cli.Do(req)
	return
}

func doHttpFormReq(ctx context.Context, urlStr string, params interface{}, options ...interface{}) (reqResp *http.Response, bodyByte []byte, err error) {
	if len(urlStr) == 0 {
		err = errors.New("missing urlStr")
		return
	}
	var (
		method = "GET"
	)
	if len(options) >= 1 {
		if methodVal, ok := options[0].(string); ok {
			method = methodVal
		}
	}
	var req *http.Request
	c := NewClient(time.Duration(30) * time.Second)
	if strings.ToUpper(method) == "GET" {
		if val, ok := params.(url.Values); ok {
			urlStr += "?" + val.Encode()
		} else if val, ok := params.(string); ok {
			urlStr += "?" + val
		} else if val, ok := params.([]byte); ok {
			urlStr += "?" + string(val)
		}
		req, err = http.NewRequest("GET", urlStr, nil)
	} else {
		var ioReader *bytes.Reader
		if val, ok := params.(url.Values); ok {
			ioReader = bytes.NewReader([]byte(val.Encode()))
		} else if val, ok := params.(string); ok {
			ioReader = bytes.NewReader([]byte(val))
		} else if val, ok := params.([]byte); ok {
			ioReader = bytes.NewReader(val)
		}
		req, err = http.NewRequest("POST", urlStr, ioReader)
	}
	if err != nil {
		return
	}
	if len(options) >= 2 {
		if headers, ok := options[0].(map[string]string); ok {
			if len(headers) > 0 {
				for k, v := range headers {
					req.Header.Add(k, v)
				}
			}
		}
		if headers, ok := options[1].(map[string]string); ok {
			if len(headers) > 0 {
				for k, v := range headers {
					req.Header.Set(k, v)
				}
			}
		}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "GoOVOPark/1.0.0 Go/1.17.6")
	reqResp, err = c.Do(ctx, req)
	if err != nil {
		return
	}
	if reqResp.Body != nil {
		defer reqResp.Body.Close()
		bodyByte, err = ioutil.ReadAll(reqResp.Body)
	}
	return
}
