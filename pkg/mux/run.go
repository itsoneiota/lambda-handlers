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

		for _, interceptor := range h.interceptors() {
			resp = interceptor(ctx, req, resp)
		}

		h.writeResponse(resp, w)

		return
	}
}

func (h *Handler) writeResponse(resp *handler.Response, w http.ResponseWriter) {
	if resp.Headers == nil {
		resp.Headers = http.Header{}
	}

	for key, values := range h.headers() {
		for _, value := range values {
			resp.Headers.Add(key, value)
		}
	}

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

func errorResponse(w http.ResponseWriter, srvErr *serviceerror.ServiceError) {
	w.WriteHeader(srvErr.StatusCode())

	b, err := json.Marshal(srvErr)
	if err != nil {
		slog.Error(err.Error())
		b, _ = json.Marshal(serviceerror.InternalServerError(err.Error()))
	}

	w.Write(b)
}
