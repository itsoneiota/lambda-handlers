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
		ctx := NewContext(r)
		req := NewRequest(r)

		f := h.function
		for i := len(h.middlewares()) - 1; i >= 0; i-- {
			f = h.middlewares()[i](f)
		}

		resp := f(ctx, req)

		handlerResponse(resp, w)

		for _, interceptor := range h.interceptors() {
			resp = interceptor(resp)
		}

		for k, v := range h.headers() {
			resp.Headers.Add(k, v[0])
		}

		writeResponse(resp, w)

		return
	}
}

func writeResponse(resp *handler.Response, w http.ResponseWriter) {
	for key, values := range resp.Headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	if _, err := w.Write([]byte(resp.Body)); err != nil {
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
