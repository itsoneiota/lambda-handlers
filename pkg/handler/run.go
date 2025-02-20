package handler

import "net/http"

func (h *Handler) Run() HandlerFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx Contexter, req Requester) *Response {
			resp := h.Function(ctx, req)

			headers := http.Header{}
			if resp.Headers != nil {
				headers = resp.Headers
			}

			for k, v := range h.Headers {
				for _, val := range v {
					headers.Add(k, val)
				}
			}

			resp.Headers = headers

			return resp
		}
	}(h.Function)
}
