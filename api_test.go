package aptly

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testClient() *Client {
	c, _ := NewClient("http://localhost:8080")
	return c
}

func TestGet(t *testing.T) {
	client := testClient()
	resp, err := client.Get("repos")
	if err != nil {
		t.Error(err)
	}
	assert.IsType(t, &http.Response{}, resp)
}
