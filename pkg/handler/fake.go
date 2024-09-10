package handler

import "net/http"

type reponseWriter struct {
	Headers http.Header
	Body    []byte
	Status  int
}

func (r *reponseWriter) Header() http.Header {
	if r.Headers == nil {
		r.Headers = http.Header{}
	}

	return r.Headers
}

func (r *reponseWriter) Write(body []byte) (int, error) {
	r.Body = body

	return len(body), nil
}

func (r *reponseWriter) WriteHeader(status int) {
	r.Status = status
}

type Model struct {
	Success bool `json:"success"`
}
