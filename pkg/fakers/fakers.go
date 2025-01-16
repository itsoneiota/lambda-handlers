package fakers

import "net/http"

type ReponseWriter struct {
	Headers http.Header
	Body    []byte
	Status  int
}

func (r *ReponseWriter) Header() http.Header {
	if r.Headers == nil {
		r.Headers = http.Header{}
	}

	return r.Headers
}

func (r *ReponseWriter) Write(body []byte) (int, error) {
	r.Body = body

	return len(body), nil
}

func (r *ReponseWriter) WriteHeader(status int) {
	r.Status = status
}
