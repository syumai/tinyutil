package httputil

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/syumai/tinyutil/internal/net_http"
)

type Client struct {
	Transport http.RoundTripper
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.Transport.RoundTrip(req)
}

func (c *Client) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}

func (c *Client) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

var DefaultClient = &Client{
	Transport: &net_http.Transport{},
}

func Get(url string) (resp *http.Response, err error) {
	return DefaultClient.Get(url)
}

func Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	return DefaultClient.Post(url, contentType, body)
}

func PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return DefaultClient.PostForm(url, data)
}
