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
	resHander *handler.ResponseHandler,
	connector Connector,
	beforeHook handler.BeforeHandlerHook,
	afterHook AfterFindHandlerHook,
) handler.HandlerFunc {
	return func(request handler.Requester) (*handler.Response, error) {
		if beforeHook != nil {
			if err := beforeHook(request); err != nil {
				return resHander.BuildErrorResponse(err, nil)
			}
		}

		token := request.GetAuthToken()
		if err := connector.Authorize(token); err != nil {
			return resHander.BuildErrorResponse(err, nil)
		}

		postcode := request.QueryByName("postcode")

		addresses, err := connector.Find(postcode)
		if err != nil {
			return resHander.BuildErrorResponse(err, nil)
		}

		if afterHook != nil {
			if err := afterHook(addresses); err != nil {
				return resHander.BuildErrorResponse(err, nil)
			}
		}

		return resHander.BuildResponse(http.StatusOK, addresses, nil)
	}
}
