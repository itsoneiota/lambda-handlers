package aws

import "github.com/itsoneiota/lambda-handlers/pkg/handler"

type Setter func(*Opt)

type Opt struct {
	interceptors []handler.Interceptor
}

func (h *Handler) interceptors() []handler.Interceptor {
	result := []handler.Interceptor{}
	if h.Opt != nil {
		result = h.Opt.interceptors
	}

	return result
}

func WithInterceptors(i ...handler.Interceptor) Setter {
	return func(o *Opt) {
		o.interceptors = i
	}
}
