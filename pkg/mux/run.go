package mux

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/itsoneiota/lambda-handlers/pkg/handler"
	"github.com/itsoneiota/lambda-handlers/pkg/serviceerror"
)

func (h *Handler) Run() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := NewResponseWriter(h.headers())

		ctx := NewContext(r)
		req := NewRequest(r)

		for _, middleware := range h.middlewares() {
			var err error
			req, err = middleware(req)
			if err != nil {
				errorResponse(resp, serviceerror.NewFromErr(err, ""))

				writeResponse(resp, w)

				return
			}
		}

		handlerResponse(h.function(ctx, req), resp)

		for _, interceptor := range h.interceptors() {
			if err := interceptor(resp); err != nil {
				errorResponse(resp, serviceerror.NewFromErr(err, ""))

				writeResponse(resp, w)

				return
			}
		}

		writeResponse(resp, w)

		return
	}
}

func writeResponse(resp *ResponseWriter, w http.ResponseWriter) {
	for key, values := range resp.Headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.Status)
	if _, err := w.Write(resp.Body); err != nil {
		errorResponse(w, serviceerror.NewFromErr(err, ""))

		return
	}
}

func handlerResponse(r *handler.Response, w http.ResponseWriter) {
	for k, v := range r.Headers {
		if len(v) > 0 {
			w.Header().Add(k, v[0])
		}
	}

	w.WriteHeader(r.StatusCode)
	w.Write([]byte(r.Body))
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
