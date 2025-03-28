package mux

import (
	"net/http"

	"github.com/itsoneiota/lambda-handlers/v2/pkg/handler"
)

type Interceptor func(*ResponseWriter) error

type Handler struct {
	function http.HandlerFunc
	*Opt
}

func New(
	hf http.HandlerFunc,
	opts ...Setter,
) *Handler {
	result := &Handler{
		function: hf,
		Opt:      &Opt{BaseOpt: &handler.BaseOpt{}},
	}

	for _, o := range opts {
		o(result.Opt)
	}

	return result
}
