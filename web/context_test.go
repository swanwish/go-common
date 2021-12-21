package web

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServeContent(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost/test", nil)
	rw := httptest.NewRecorder()
	ctx := CreateHandlerContext(rw, req)
	ctx.ServeContent("test.json", []byte(`{"test": "value"}`))

	resp := rw.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.Equal(t, `{"test": "value"}`, string(body))
}

func TestServeJsonContent(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost/test", nil)
	rw := httptest.NewRecorder()

	data := struct {
		Key string `json:"key"`
	}{Key: "value"}

	ctx := CreateHandlerContext(rw, req)
	ctx.ServeJsonContent("test.json", data)

	resp := rw.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.Equal(t, `{"key":"value"}`, string(body))
}
