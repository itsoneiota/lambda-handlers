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
	http.Headers{},
)
```

When implemeting the lambda `Start` method you can also define before hooks (which means you can manipluate a request within you code base), or after hooks (for maniplate the response object of a handler). Any default headers that you wish to be added to your response can be defined as the parameter of the `Start` method.

### ResponseWriter
The `aws.ResponseWriter` inherits the `events.APIGatewayProxyResponse` to fulfill the `http.ResponseWriter` contract. Allow you to access the response `Header`, `Write` data and `WriteHeader` status code to the response.

With this you can also access methods that are available on the `events.APIGatewayProxyResponse`. such as `StatusCode`, `Headers`, `MultiValueHeaders`, `Body` and `IsBase64Encoded`.

### Middlewares
Middlewares can be used on a handler, which will be run before the actual handler function.

These can be be added using the `WithMiddlewares` setter method on the `Start` function.

```go
aws.Start(
	handler,
	aws.WithMiddleware(middleware)
)
```

Mutliple middlewares can be passed through in the method, which will be chained in the given order that they are passed through.

Any middleware that is being used on this must fulfill the `Middleware` contract, which is shown below:

```go
type Middleware func(*http.Request) (*http.Request, error)
```

Within the middleware you can manipulate the `http.Request`, including the query parameters:

```go
func(r *http.Request) (*http.Request, error) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return nil, err
	}

	query.Set("bar", "3")

	r.URL.RawQuery = query.Encode()

	return r, nil
}
```

and request context:

```go
func(r *http.Request) (*http.Request, error) {
	ctx := context.WithValue(r.Context(), "baz", "4")

	return r.WithContext(ctx), nil
},
```

### Interceptors
Interceptors can be used to manipulate the handler response before it is commuincated back in the request.

These can be be added using the `WithInterceptors` setter method on the `Start` function.

```go
aws.Start(
	handler,
	aws.WithInterceptors(interceptor)
)
```

Any interceptor that is being used on this must fulfill the `Interceptor` contract, which is shown below:

```go
type Interceptor func(*aws.ResponseWriter) error
```

In order to manipulate the handler response you can take the `aws.ResponseWriter` `Body` to change the response, and then add it back to the `aws.ResponseWriter`, e.g.:

```go
func(w *ResponseWriter) error {
	m := &metasyntactic{}
	err := json.Unmarshal([]byte(w.Body), m)
	s.NoError(err)

	m.Bar = "4"

	b, err := json.Marshal(m)
	s.NoError(err)

	w.Write(b)

	return nil
},
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.


## License
[MIT](https://choosealicense.com/licenses/mit/)
