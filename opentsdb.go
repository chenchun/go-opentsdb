package opentsdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Options struct {
	// Host value for the opentsdb server
	// Default: 127.0.0.1
	Host string

	// Port for the opentsdb server
	// Default: 4242
	Port int

	// Timeout for http client
	// Default: no timeout
	Timeout time.Duration
}

type Client struct {
	url        *url.URL
	httpClient *http.Client
	tr         *http.Transport
}

func NewClient(opt Options) (*Client, error) {
	if opt.Host == "" {
		opt.Host = "127.0.0.1"
	}
	if opt.Port == 0 {
		opt.Port = 4242
	}

	u, err := url.Parse(fmt.Sprintf("http://%s:%d", opt.Host, opt.Port))
	if err != nil {
		return nil, err
	}

	tr := &http.Transport{}

	return &Client{
		url: u,
		httpClient: &http.Client{
			Timeout:   opt.Timeout,
			Transport: tr,
		},
		tr: tr,
	}, nil
}

func (c *Client) Close() error {
	c.tr.CloseIdleConnections()
	return nil
}

func (c *Client) Aggregators() error {
	return nil
}

func (c *Client) Annotation() error {
	return nil
}

func (c *Client) Config() error {
	return nil
}

func (c *Client) Dropcaches() error {
	return nil
}

func (c *Client) Put(bp *BatchPoints, params string) (int, []byte, error) {
	data, err := bp.ToJson()
	if err != nil {
		return 0, nil, err
	}

	u := c.url
	u.Path = "api/put"
	u.RawQuery = params

	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(data))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, body, err
}

func (c *Client) Query(q *QueryParams) (int, []byte, error) {
	data, err := json.Marshal(q)
	if err != nil {
		return 0, nil, err
	}

	u := c.url
	u.Path = "api/query"

	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(data))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, body, nil
}

func (c *Client) Search() error {
	return nil
}

func (c *Client) Serializers() error {
	return nil
}

func (c *Client) Stats() error {
	return nil
}

func (c *Client) Suggest() error {
	return nil
}

func (c *Client) Tree() error {
	return nil
}

func (c *Client) Uid() error {
	return nil
}

func (c *Client) Version() error {
	return nil
}
