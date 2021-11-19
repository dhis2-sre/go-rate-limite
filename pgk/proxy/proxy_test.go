package proxy

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestTransparentProxyHandler(t *testing.T) {
	expected := http.StatusOK

	backend := "https://httpbin.org"

	backendUrl, err := url.Parse(backend)
	assert.NoError(t, err)

	proxy := ProvideTransparentProxy(backendUrl)

	req, err := http.NewRequest("GET", "/status/200", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(proxy)

	handler.ServeHTTP(recorder, req)

	actual := recorder.Code
	assert.Equal(t, expected, actual)
}
