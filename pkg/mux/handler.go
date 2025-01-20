package mux

import (
	"github.com/itsoneiota/lambda-handlers/pkg/handler"
)

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
