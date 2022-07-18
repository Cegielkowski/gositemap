package mocks

import "net/http"

type MockClient struct {
    GetFunc func(url string) (resp *http.Response, err error)
}

var (
    GetDoFunc func(url string) (resp *http.Response, err error)
)

func (m *MockClient) Get(url string) (resp *http.Response, err error) {
    return GetDoFunc(url)
}
