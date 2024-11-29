package handler

type Middleware func(HandlerFunc) HandlerFunc

func (rh *Handler) Middlewares(h HandlerFunc, middlewares ...Middleware) HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}

	return h
}
