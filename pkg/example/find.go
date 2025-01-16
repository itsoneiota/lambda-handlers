package example

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/itsoneiota/lambda-handlers/v2/pkg/serviceerror"
)

type ExampleModel struct {
	Success bool `json:"success"`
}

type Connector interface {
	Authorize(token string) error
	Find(postcode string) (interface{}, error)
}

const findHandlerDefaultCount = 10

// FindHandler returns a handlers.HandlerFunc which is used for the Find endpoint.
// The handler calls the Find method of the connector
func FindHandler(connector Connector) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("Authorization")
		if err := connector.Authorize(token); err != nil {
			e := serviceerror.NewFromErr(err, "")

			w.WriteHeader(e.StatusCode())
			w.Write(e.Bytes())

			return
		}

		query, err := url.ParseQuery(req.URL.RawQuery)
		if err != nil {
			e := serviceerror.NewFromErr(err, "")

			w.WriteHeader(e.StatusCode())
			w.Write(e.Bytes())

			return
		}

		var postcode string
		if query.Has("postcode") {
			postcode = query.Get("postcode")
		} else {
			e := serviceerror.BadRequest("postcode required")

			w.WriteHeader(e.StatusCode())
			w.Write(e.Bytes())

			return
		}

		addresses, err := connector.Find(postcode)
		if err != nil {
			e := serviceerror.NewFromErr(err, "")

			w.WriteHeader(e.StatusCode())
			w.Write(e.Bytes())

			return
		}

		b, err := json.Marshal(addresses)
		if err != nil {
			e := serviceerror.NewFromErr(err, "")

			w.WriteHeader(e.StatusCode())
			w.Write(e.Bytes())

			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}
