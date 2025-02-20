package handler

type Middleware func(HandlerFunc) HandlerFunc

func (h *Handler) Middlewares(middlewares ...Middleware) *Handler {
	f := h.Function
	for _, middleware := range middlewares {
		f = middleware(f)
	}

	h.Function = f

	return h
}
