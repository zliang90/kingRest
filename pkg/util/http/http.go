package http

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/zliang90/kingRest/pkg/util/slice"
)

// default timeout
var defaultTimeout = 1 * time.Second

// default retry wait time
var defaultRetryWaitTime = 2 * time.Second

// http method
var httpRequestMethods = []string{"GET", "POST", "PATCH", "PUT", "HEAD", "DELETE"}

type Option struct {
	Method      string            `mapstructure:"method" yaml:"method" json:"method"`
	Headers     map[string]string `mapstructure:"headers" yaml:"headers" json:"headers"`
	QueryParams map[string]string `mapstructure:"queryParams" yaml:"queryParams" json:"queryParams"`
	Timeout     time.Duration     `mapstructure:"timeout" yaml:"timeout" json:"timeout"`
	Retry       int               `mapstructure:"retry" json:"retry"`
	FormData    map[string]string `mapstructure:"formData" yaml:"formData" json:"formData"`
	Body        []byte            `mapstructure:"body" yaml:"body" json:"body"`
	Logger      resty.Logger
}

type Client struct {
	c      *resty.Client
	option *Option
}

func NewClient(o *Option) *Client {
	c := resty.New()

	if o != nil {
		if o.Timeout < defaultTimeout {
			o.Timeout = defaultTimeout
		}

		if o.Headers != nil {
			c.SetHeaders(o.Headers)
		}

		c.SetTimeout(o.Timeout).
			SetRetryCount(o.Retry).
			SetRetryWaitTime(defaultRetryWaitTime).
			SetRetryMaxWaitTime(defaultRetryWaitTime * 2)

		if o.Method == "" {
			o.Method = "GET"
		}
	}

	return &Client{c: c, option: o}
}

func (client *Client) SetOption(o *Option) {
	client.option = o
}

func (client *Client) SetTimeout(t time.Duration) {
	client.option.Timeout = t
}

func (client *Client) SetRetry(n int) {
	client.option.Retry = n
}

func (client *Client) SetQueryParams(q map[string]string) {
	client.option.QueryParams = q
}
func (client *Client) SetFormData(d map[string]string) {
	client.option.FormData = d
}

func (client *Client) SetHeaders(headers map[string]string) {
	client.option.Headers = headers
}

func (client *Client) SetBody(body []byte) {
	client.option.Body = body
}

func (client *Client) SetLogger(l resty.Logger) {
	client.c.SetLogger(l)
}

func (client *Client) Request(url string) (*resty.Response, error) {
	// request method type
	var f func(string) (*resty.Response, error)

	if !slice.HasString(httpRequestMethods, client.option.Method) {
		return nil, fmt.Errorf("invalid http method: '%s'", client.option.Method)
	}

	r := client.c.R()

	if client.option.QueryParams != nil {
		r.SetQueryParams(client.option.QueryParams)
	}
	if client.option.Body != nil {
		r.SetBody(client.option.Body)
	}
	if client.option.FormData != nil {
		r.SetFormData(client.option.FormData)
	}

	switch client.option.Method {
	case "GET":
		f = r.Get
	case "POST":
		f = r.Post
	case "HEAD":
		f = r.Head
	case "PATCH":
		f = r.Patch
	case "DELETE":
		f = r.Delete
	}
	resp, err := f(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
