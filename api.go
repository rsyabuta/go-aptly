package aptly

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Aptly struct {
	Url string
}

func (aptly *Aptly) Get(endpoint string) (*http.Response, error) {
	request, _ := http.NewRequest("GET", aptly.requestUrl(endpoint), nil)
	return aptly.makeRequest(request)
}

func (aptly *Aptly) Post(endpoint string, contentType string, params map[string]string, body io.Reader) (*http.Response, error) {
	query, err := aptly.buildQueryString(endpoint, params)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", query, body)
	request.Header.Set("Content-Type", contentType)
	return aptly.makeRequest(request)
}

func (aptly *Aptly) Delete(endpoint string, params map[string]string) (*http.Response, error) {
	query, err := aptly.buildQueryString(endpoint, params)
	if err != nil {
		return nil, err
	}

	request, _ := http.NewRequest("DELETE", query, nil)
	return aptly.makeRequest(request)
}

func (aptly *Aptly) Put(endpoint string, params map[string]string, body io.Reader) (*http.Response, error) {
	query, err := aptly.buildQueryString(endpoint, params)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("PUT", query, body)
	request.Header.Set("Content-Type", "application/json")
	return aptly.makeRequest(request)
}

func (aptly *Aptly) makeRequest(request *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(request)
}

func (aptly *Aptly) requestUrl(endpoint string) string {
	return fmt.Sprintf("%s/api/%s", aptly.Url, endpoint)
}

func (aptly *Aptly) buildQueryString(endpoint string, params map[string]string) (string, error) {
	u, err := url.Parse(aptly.requestUrl(endpoint))
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
