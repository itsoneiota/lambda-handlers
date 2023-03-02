package mux

import (
	"net/http"

	"github.com/itsoneiota/lambda-handlers/pkg/handler"
)

func CreateHandler(
	h handler.HandlerFunc,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _ := h(NewRequest(r))

		WriteResponse(res, w)
	}
}

func WriteResponse(r *handler.Response, w http.ResponseWriter) {
	for k, v := range r.Headers {
		if len(v) > 0 {
			w.Header().Add(k, v[0])
		}
	}

	w.WriteHeader(r.StatusCode)
	w.Write([]byte(r.Body))
}
