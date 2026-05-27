package mux

import (
	"bytes"
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func NewResponseWriter(w http.ResponseWriter, headers http.Header) *ResponseWriter {
	for k, v := range headers {
		if len(v) > 0 {
			w.Header().Set(k, v[0])
		}
	}

	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (w *ResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	w.body.Reset()

	return w.body.Write(b)
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *ResponseWriter) Body() string {
	return w.body.String()
}

func (w *ResponseWriter) SetBody(b []byte) {
	w.body.Reset()
	w.body.Write(b)
}

func (w *ResponseWriter) StatusCode() int {
	return w.statusCode
}

func (w *ResponseWriter) SetHttpResponseWriter() error {
	w.ResponseWriter.WriteHeader(w.statusCode)

	_, err := w.ResponseWriter.Write(w.body.Bytes())
	return err
}
