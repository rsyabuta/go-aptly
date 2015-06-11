package aptly

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	Url string

	File          *FileService
	LocalRepo     *LocalRepoService
	Snapshot      *SnapshotService
	PublishedRepo *PublishedRepoService
}

func NewClient(url string) (*Client, error) {
	c := &Client{
		Url: url,
	}

	c.File = &FileService{Client: c}
	c.LocalRepo = &LocalRepoService{Client: c}
	c.Snapshot = &SnapshotService{Client: c}
	c.PublishedRepo = &PublishedRepoService{Client: c}
	return c, nil
}

func (client *Client) Get(endpoint string) (*http.Response, error) {
	request, _ := http.NewRequest("GET", client.requestUrl(endpoint), nil)
	return client.makeRequest(request)
}

func (client *Client) Post(endpoint string, contentType string, params map[string]string, body io.Reader) (*http.Response, error) {
	query, err := client.buildQueryString(endpoint, params)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", query, body)
	request.Header.Set("Content-Type", contentType)
	return client.makeRequest(request)
}

func (client *Client) Delete(endpoint string, params map[string]string) (*http.Response, error) {
	query, err := client.buildQueryString(endpoint, params)
	if err != nil {
		return nil, err
	}

	request, _ := http.NewRequest("DELETE", query, nil)
	return client.makeRequest(request)
}

func (client *Client) Put(endpoint string, params map[string]string, body io.Reader) (*http.Response, error) {
	query, err := client.buildQueryString(endpoint, params)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("PUT", query, body)
	request.Header.Set("Content-Type", "application/json")
	return client.makeRequest(request)
}

func (client *Client) makeRequest(request *http.Request) (*http.Response, error) {
	c := &http.Client{}
	return c.Do(request)
}

func (client *Client) requestUrl(endpoint string) string {
	return fmt.Sprintf("%s/api/%s", client.Url, endpoint)
}

func (client *Client) buildQueryString(endpoint string, params map[string]string) (string, error) {
	u, err := url.Parse(client.requestUrl(endpoint))
	if err != nil {
		return "", err
	}

	query := u.Query()
	for k, v := range params {
		query.Set(k, v)
	}
	u.RawQuery = query.Encode()
	if err != nil {
		return "", err
	}

	return u.String(), nil
}
