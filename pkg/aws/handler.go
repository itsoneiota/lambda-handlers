package aws

import "github.com/itsoneiota/lambda-handlers/pkg/handler"

type Interceptor func(*handler.Response) *handler.Response

type Handler struct {
	handler handler.HandlerFunc
	*Opt
}
