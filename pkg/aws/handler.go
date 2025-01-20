package aws

import "github.com/itsoneiota/lambda-handlers/pkg/handler"

type Handler struct {
	handler handler.HandlerFunc
	*Opt
}
