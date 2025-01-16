package aws

import (
	"net/http"
)

type Setter func(*Opt)

type Opt struct {
	headers      http.Header
	middlewares  []Middleware
	interceptors []Interceptor
}

func (h *Handler) headers() http.Header {
	result := http.Header{}
	if h.Opt != nil {
		result = h.Opt.headers
	}

	return result
}

func (h *Handler) middlewares() []Middleware {
	result := []Middleware{}
	if h.Opt != nil {
		result = h.Opt.middlewares
	}

	return result
}

func (h *Handler) interceptors() []Interceptor {
	result := []Interceptor{}
	if h.Opt != nil {
		result = h.Opt.interceptors
	}

	return result
}

func WithHeaders(h http.Header) Setter {
	return func(o *Opt) {
		o.headers = h
	}
}

func WithMiddlewares(m ...Middleware) Setter {
	return func(o *Opt) {
		o.middlewares = m
	}
}

func WithInterceptors(i ...Interceptor) Setter {
	return func(o *Opt) {
		o.interceptors = i
	}
}
