package mux

import (
	"net/http"

	"github.com/itsoneiota/lambda-handlers/v2/pkg/handler"
)

type Handler struct {
	function http.HandlerFunc
	*handler.BaseOpt
}

func New(
	hf http.HandlerFunc,
	opts ...handler.Setter,
) *Handler {
	result := &Handler{
		function: hf,
		BaseOpt:  &handler.BaseOpt{},
	}

	for _, o := range opts {
		o(result.BaseOpt)
	}

	return result
}
