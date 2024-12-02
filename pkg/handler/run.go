package handler

import "net/http"

func (h *Handler) Run() HandlerFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx Contexter, req Requester) *Response {
			resp := h.function(ctx, req)

			headers := http.Header{}
			if resp.Headers != nil {
				headers = resp.Headers
			}

			for k, v := range h.headers {
				for _, val := range v {
					headers.Add(k, val)
				}
			}

			resp.Headers = headers

			return resp
		}
	}(h.function)
}
