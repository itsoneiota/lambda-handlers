package example

import (
	"net/http"

	"github.com/itsoneiota/lambda-handlers/pkg/handler"
)

type ExampleModel struct {
	Success bool `json:"success"`
}

type Connector interface {
	Authorize(token string) error
	Find(postcode string) (interface{}, error)
}

const findHandlerDefaultCount = 10

// AfterFindHandlerHook is a hook/callback function definition, triggered after the Find connector call on for the FindHandler
type AfterFindHandlerHook func(interface{}) error

// FindHandler returns a handlers.HandlerFunc which is used for the Find endpoint.
// The handler calls the Find method of the connector
func FindHandler(
	connector Connector,
	afterHook AfterFindHandlerHook,
) handler.HandlerFunc {
	return func(ctx handler.Contexter, request handler.Requester) *handler.Response {
		token := request.GetAuthToken()
		if err := connector.Authorize(token); err != nil {
			return handler.NewErrorResponse(err)
		}

		postcode := request.QueryByName("postcode")

		addresses, err := connector.Find(postcode)
		if err != nil {
			return handler.NewErrorResponse(err)
		}

		if afterHook != nil {
			if err := afterHook(addresses); err != nil {
				return handler.NewErrorResponse(err)
			}
		}

		return handler.NewResponse(http.StatusOK, addresses)
	}
}
