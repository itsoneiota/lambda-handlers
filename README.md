# Lambda Handlers

Lambda Handlers is a go module allowing Serverless handler functions to be ran as a local [Gorilla Mux](https://github.com/gorilla/mux) server or within any cloud server provider event.

Currently supported:
 - AWS Lambda.
 - Standard library HTTP.

## Usage

The lambda handler uses the default `http.Request` core package to handle any request, there conversion of other request packages must be done in order to use this package. There is a built in aws handler, which converter `events.APIGatewayProxyRequest` to `http.Request`. The handler also uses response writer, therefore fulfilling the contract of using a `mux` router.

```go
package example

import (
	"net/http"

	"github.com/itsoneiota/lambda-handlers/v2/pkg/handler"
)

type ExampleModel struct {
	Success bool `json:"success"`
}

type Connector interface {
	Authorize(token string) error
	Find(query string) (interface{}, error)
}

const findHandlerDefaultCount = 10

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

```

In the case where you want to run this handler in a Mux router, call the `CreateHandler` method, pass in the generic handler defined above and pass it into the HandleFunc method on the router.

```go
r := muxRouter.NewRouter()
r.HandleFunc("/test", mux.CreateHandler(handler))

log.Fatal(http.ListenAndServe("localhost:8080", r))
```

In the case where you want to run this handler in AWS Lambda, simply pass the handler into the `Start` method found within the `aws` package of this module.
```go

aws.Start(
	handler,
	nil,
	nil,
	http.Headers{},
)
```

When implemeting the lambda `Start` method you can also define before hooks (which means you can manipluate a request within you code base), or after hooks (for maniplate the response object of a handler). Any default headers that you wish to be added to your response can be defined as the parameter of the `Start` method.

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.


## License
[MIT](https://choosealicense.com/licenses/mit/)
