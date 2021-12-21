package web

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeContent(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://localhost/test", nil)
	rw := httptest.NewRecorder()
	ctx := CreateHandlerContext(rw, req)
	ctx.ServeContent("test.json", []byte(`{"test": "value"}`))

	resp := rw.Result()
	body, _ := io.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))
}
