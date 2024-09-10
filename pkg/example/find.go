package example

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/itsoneiota/lambda-handlers/v2/pkg/handler"
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
) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("Authorization")
		if err := connector.Authorize(token); err != nil {
			resHander.BuildErrorResponse(w, err)
		}

		query, err := url.ParseQuery(req.URL.RawQuery)
		if err != nil {
			resHander.BuildErrorResponse(w, err)
		}

		var postcode string
		if query.Has("postcode") {
			postcode = query.Get("postcode")
		} else {
			resHander.BuildErrorResponse(w, errors.New("postcode required"))
		}

		addresses, err := connector.Find(postcode)
		if err != nil {
			resHander.BuildErrorResponse(w, err)
		}

		resHander.BuildResponse(w, http.StatusOK, addresses)
	}
}
