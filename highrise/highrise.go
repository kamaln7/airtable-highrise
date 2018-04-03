package highrise

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	httpClient *http.Client
	url, key   string
}

func New(key, team string) *Client {
	client := &Client{
		httpClient: &http.Client{},
		key:        key,
		url:        fmt.Sprintf("https://%s.highrisehq.com", team),
	}

	return client
}

func (c *Client) request(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", c.url, strings.TrimLeft(url, "/")), body)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.key, "X")
	req.Header.Set("User-Agent", "airtable-highrise/0.1")
	req.Header.Set("Content-Type", "text/xml")
	resp, err := c.httpClient.Do(req)

	return resp, err
}
