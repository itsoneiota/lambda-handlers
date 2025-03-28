package aws

import (
	"net/http"

	"github.com/itsoneiota/lambda-handlers/v2/pkg/handler"
)

type AwsOpt interface {
	handler.Opt
	SetInterceptors([]Interceptor)
}

type Setter func(AwsOpt)

type Opt struct {
	*handler.BaseOpt
	interceptors []Interceptor
}

func (o *Opt) SetHeaders(h http.Header) {
	o.BaseOpt.SetHeaders(h)
}

func (o *Opt) SetMiddlewares(m []handler.Middleware) {
	o.BaseOpt.SetMiddlewares(m)
}

func (o *Opt) SetInterceptors(i []Interceptor) {
	o.interceptors = i
}

func (h *Handler) headers() http.Header {
	result := http.Header{}
	if h.Opt != nil && h.Opt.BaseOpt != nil {
		result = h.Opt.BaseOpt.Headers()
	}

	return result
}

func (h *Handler) middlewares() []handler.Middleware {
	result := []handler.Middleware{}
	if h.Opt != nil && h.Opt.BaseOpt != nil {
		result = h.Opt.BaseOpt.Middlewares()
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
	return func(o AwsOpt) {
		o.SetHeaders(h)
	}
}

func WithMiddlewares(m ...handler.Middleware) Setter {
	return func(o AwsOpt) {
		o.SetMiddlewares(m)
	}
}

func WithInterceptors(i ...Interceptor) Setter {
	return func(o AwsOpt) {
		o.SetInterceptors(i)
	}
}
