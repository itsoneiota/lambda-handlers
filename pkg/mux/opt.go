package mux

import (
	"net/http"

	"github.com/itsoneiota/lambda-handlers/pkg/handler"
)

type Setter func(*opt)

type opt struct {
	headers      http.Header
	middlewares  []handler.Middleware
	interceptors []handler.Interceptor
}

func (h *Handler) headers() http.Header {
	result := http.Header{}
	if h.opt != nil {
		result = h.opt.headers
	}

	return result
}

func (h *Handler) middlewares() []handler.Middleware {
	result := []handler.Middleware{}
	if h.opt != nil {
		result = h.opt.middlewares
	}

	return result
}

func (h *Handler) interceptors() []handler.Interceptor {
	result := []handler.Interceptor{}
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

func WithMiddlewares(m ...handler.Middleware) Setter {
	return func(o *opt) {
		o.middlewares = m
	}
}

func WithInterceptors(i ...handler.Interceptor) Setter {
	return func(o *opt) {
		o.interceptors = i
	}
}
