package aptly

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGet(t *testing.T) {
	aptly := new(Aptly)
	aptly.Url = "http://localhost:8080"
	resp, err := aptly.Get("repos")
	if err != nil {
		t.Error(err)
	}
	assert.IsType(t, &http.Response{}, resp)
}
