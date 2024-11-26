package handler

func (h *Handler) Run(f HandlerFunc) HandlerFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx Contexter, req Requester) (*Response, error) {
			resp, err := f(ctx, req)
			if err != nil {
				return h.BuildErrorResponse(err)
			}

			return resp, nil
		}
	}(f)
}
