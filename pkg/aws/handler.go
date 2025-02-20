package aws

import "github.com/itsoneiota/lambda-handlers/pkg/handler"

type Handler struct {
	handler *handler.Handler
	*Opt
}
