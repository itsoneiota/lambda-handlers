package mux

import (
	"net/http"

	"github.com/itsoneiota/lambda-handlers/pkg/handler"
)

func CreateHandler(h handler.HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, NewRequest(r))
	}
}
