package handler

import "net/http"

func (h *Handler) Run(f HandlerFunc) HandlerFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx Contexter, req Requester) *Response {
			resp := f(ctx, req)

			headers := http.Header{}
			if resp.Headers != nil {
				headers = resp.Headers
			}

			for k, v := range h.defaultHeaders {
				for _, val := range v {
					headers.Add(k, val)
				}
			}

			resp.Headers = headers

			return resp
		}
	}(f)
}
