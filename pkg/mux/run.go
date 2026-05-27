package mux

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/itsoneiota/lambda-handlers/v2/pkg/helpers"
	"github.com/itsoneiota/lambda-handlers/v2/pkg/serviceerror"
)

func (h *Handler) Run() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := NewResponseWriter(w, h.Headers())

		for _, middleware := range h.Middlewares() {
			var err error
			r, err = middleware(r)
			if err != nil {
				errorResponse(w, serviceerror.NewFromErr(err, ""))
				return
			}
		}

		h.function(resp, r)

		if !helpers.IsOkRange(resp.StatusCode()) {
			resp.SetHttpResponseWriter()

			return
		}

		for _, interceptor := range h.Interceptors() {
			if err := interceptor(r, resp); err != nil {
				errorResponse(w, serviceerror.NewFromErr(err, ""))
				return
			}
		}

		resp.SetHttpResponseWriter()

		return
	}
}

func errorResponse(w http.ResponseWriter, srvErr *serviceerror.ServiceError) {
	w.WriteHeader(srvErr.StatusCode())

	b, err := json.Marshal(srvErr)
	if err != nil {
		slog.Error(err.Error())
		b, _ = json.Marshal(serviceerror.InternalServerError(err.Error()))
	}

	w.Write(b)
}
