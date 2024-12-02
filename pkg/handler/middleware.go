package handler

type Middleware func(HandlerFunc) HandlerFunc

func (h *Handler) Middlewares(middlewares ...Middleware) *Handler {
	f := h.function
	for _, middleware := range middlewares {
		f = middleware(f)
	}

	h.function = f

	return h
}
