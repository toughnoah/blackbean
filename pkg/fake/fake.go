package fake

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type MockEsResponse struct {
	ResponseString string
}

func (t *MockEsResponse) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(strings.NewReader(t.ResponseString))}, nil
}
