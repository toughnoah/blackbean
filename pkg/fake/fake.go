package fake

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type Mock interface {
	RoundTrip(*http.Request) (*http.Response, error)
}

type MockEsResponse struct {
	ResponseString string
}

func (t *MockEsResponse) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(strings.NewReader(t.ResponseString))}, nil
}

type MockErrorEsResponse struct{}

func (t *MockErrorEsResponse) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(strings.NewReader("send request failed"))}, errors.New("mock failed response")
}
