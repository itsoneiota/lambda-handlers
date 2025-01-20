package mux

import (
	"github.com/itsoneiota/lambda-handlers/pkg/handler"
)

type Middleware func(*Request) (*Request, error)
type Interceptor func(*ResponseWriter) error

type Handler struct {
	function handler.HandlerFunc
	*opt
}

func New(
	function handler.HandlerFunc,
	opts ...Setter,
) *Handler {
	result := &Handler{
		function: function,
		opt:      &opt{},
	}

	for _, o := range opts {
		o(result.opt)
	}

	return result
}
