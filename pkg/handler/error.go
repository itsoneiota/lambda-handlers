package handler

// ServiceError - represents the service error
type ServiceError struct {
	Err Error `json:"error"`
}

// Error holds the error contents of the service error
type Error struct {
	ID      string `json:"id"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewServiceError(id, code, message string) *ServiceError {
	return &ServiceError{
		Err: Error{
			ID:      id,
			Code:    code,
			Message: message,
		},
	}
}

// Error returns the error message
func (se *ServiceError) Error() string {
	return se.Err.Message
}
