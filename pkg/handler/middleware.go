package handler

type Middleware func(HandlerFunc) HandlerFunc

func (rh *Handler) Middlewares(h HandlerFunc, middlewares ...Middleware) HandlerFunc {
	for _, middleware := range middlewares {
		h = middleware(h)
	}

	return h
}
