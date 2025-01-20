package mux

import "net/http"

type Setter func(*opt)

type opt struct {
	headers      http.Header
	middlewares  []Middleware
	interceptors []Interceptor
}

func (h *Handler) headers() http.Header {
	result := http.Header{}
	if h.opt != nil {
		result = h.opt.headers
	}

	return result
}

func (h *Handler) middlewares() []Middleware {
	result := []Middleware{}
	if h.opt != nil {
		result = h.opt.middlewares
	}

	return result
}

func (h *Handler) interceptors() []Interceptor {
	result := []Interceptor{}
	if h.opt != nil {
		result = h.opt.interceptors
	}

	return result
}

func WithHeaders(h http.Header) Setter {
	return func(o *opt) {
		o.headers = h
	}
}

func WithMiddlewares(m ...Middleware) Setter {
	return func(o *opt) {
		o.middlewares = m
	}
}

func WithInterceptors(i ...Interceptor) Setter {
	return func(o *opt) {
		o.interceptors = i
	}
}
