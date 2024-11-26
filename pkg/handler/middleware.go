package handler

type Middleware func(HandlerFunc) HandlerFunc

func (rh *Handler) Middlewares(h HandlerFunc, middlewares ...Middleware) HandlerFunc {
	for _, middleware := range middlewares {
		h = func(next HandlerFunc) HandlerFunc {
			return func(ctx Contexter, req Requester) *Response {
				return middleware(next)(ctx, req)
			}
		}(h)
	}

	return h
}
