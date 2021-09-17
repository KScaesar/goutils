package xTest

import (
	"io"
	"net/http"
	"net/http/httptest"
)

func HttpClientDoJson(router http.Handler, httpMethod string, path string, body io.Reader) (respBody string, status int) {
	wRecorder := httptest.NewRecorder()
	req := httptest.NewRequest(httpMethod, path, body)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wRecorder, req)

	resp := wRecorder.Result()
	defer resp.Body.Close()
	return wRecorder.Body.String(), resp.StatusCode
}
