package aws

type Setter func(*Opt)

type Opt struct {
	interceptors []Interceptor
}

func (h *Handler) interceptors() []Interceptor {
	result := []Interceptor{}
	if h.Opt != nil {
		result = h.Opt.interceptors
	}

	return result
}

func WithInterceptors(i ...Interceptor) Setter {
	return func(o *Opt) {
		o.interceptors = i
	}
}
