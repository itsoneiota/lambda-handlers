package mux

import "net/http"

type fakeResponseWriter struct {
	headers    http.Header
	statusCode int
	body       []byte
}

func (w *fakeResponseWriter) Header() http.Header {
	return w.headers
}

func (w *fakeResponseWriter) Write(body []byte) (int, error) {
	w.body = body

	return len(body), nil
}

func (w *fakeResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}
