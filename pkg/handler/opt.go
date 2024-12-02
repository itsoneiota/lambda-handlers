package handler

import "net/http"

type Setter func(*opt)

type opt struct {
	headers http.Header
}

func WithHeaders(headers http.Header) Setter {
	return func(o *opt) {
		o.headers = headers
	}
}
