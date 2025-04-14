package handler

import "net/http"

type Setter func(*BaseOpt)

type BaseOpt struct {
	headers      http.Header
	middlewares  []Middleware
	interceptors []Interceptor
}

func (o *BaseOpt) Headers() http.Header {
	return o.headers
}

func (o *BaseOpt) SetHeaders(h http.Header) {
	o.headers = h
}

func (o *BaseOpt) Middlewares() []Middleware {
	return o.middlewares
}

func (o *BaseOpt) SetMiddlewares(m []Middleware) {
	o.middlewares = m
}

func (o *BaseOpt) Interceptors() []Interceptor {
	return o.interceptors
}

func (o *BaseOpt) SetInterceptors(i []Interceptor) {
	o.interceptors = i
}

func WithHeaders(h http.Header) Setter {
	return func(o *BaseOpt) {
		o.headers = h
	}
}

func WithMiddlewares(m ...Middleware) Setter {
	return func(o *BaseOpt) {
		o.middlewares = m
	}
}

func WithInterceptors(i ...Interceptor) Setter {
	return func(o *BaseOpt) {
		o.interceptors = i
	}
}
