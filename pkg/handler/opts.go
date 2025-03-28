package handler

import "net/http"

type Setter func(Opt)

type Opt interface {
	SetHeaders(http.Header)
	SetMiddlewares([]Middleware)
}

type BaseOpt struct {
	headers     http.Header
	middlewares []Middleware
}

func (o *BaseOpt) Headers() http.Header {
	return o.headers
}

func (o *BaseOpt) SetHeaders(h http.Header) {
	o.headers = h
}

func (o *BaseOpt) SetMiddlewares(m []Middleware) {
	o.middlewares = m
}

func (o *BaseOpt) Middlewares() []Middleware {
	return o.middlewares
}

func WithHeaders(h http.Header) Setter {
	return func(o Opt) {
		o.SetHeaders(h)
	}
}

func WithMiddlewares(m ...Middleware) Setter {
	return func(o Opt) {
		o.SetMiddlewares(m)
	}
}
