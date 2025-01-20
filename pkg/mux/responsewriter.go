package mux

import "net/http"

type ResponseWriter struct {
	Headers http.Header
	Body    []byte
	Status  int
}

func NewResponseWriter(headers http.Header) *ResponseWriter {
	return &ResponseWriter{
		Headers: headers,
	}
}

func (r *ResponseWriter) Header() http.Header {
	if r.Headers == nil {
		r.Headers = http.Header{}
	}

	return r.Headers
}

func (r *ResponseWriter) Write(body []byte) (int, error) {
	r.Body = body

	return len(body), nil
}

func (r *ResponseWriter) WriteHeader(status int) {
	r.Status = status
}
