# Lambda Handlers

Lambda Handlers is a go module allowing Serverless handler functions to be ran as a local [Gorilla Mux](https://github.com/gorilla/mux) server or within any cloud server provider event.

Currently supported:
 - AWS Lambda.
 - Standard library HTTP.

## Usage

The first step is to swap our your CSP specific event request and response objects with the generic `Requester` and `Responder` interfaces defined in the handler package of this module.

```go
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
	Find(query string) (interface{}, error)
}

const findHandlerDefaultCount = 10

func FindHandler(
	connector Connector,
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

		return handler.NewResponse(http.StatusOK, addresses)
	}
}
```

In the case where you want to run this handler in a Mux router, call the `CreateHandler` method, pass in the generic handler defined above and pass it into the HandleFunc method on the router.

```go
r.HandleFunc("/test", mux.CreateHandler(handler.New(example.FindHandler(c, nil, nil)).Run()))

log.Fatal(http.ListenAndServe("localhost:8080", r))
```

In the case where you want to run this handler in AWS Lambda, simply pass the handler into the `Start` method found within the `aws` package of this module.
```go

aws.Start(New(testHandler).Run())
```

### Middleware

Middleware can be abled by using the `Middlewares` method on the handler:
```go
aws.Start(New(testHandler).Middleware(
	testMiddleware,
).Run())

```

A middleware must fulfill the `Middleware` type contract:
```go
type Middleware func(HandlerFunc) HandlerFunc
```

An example usage of this is below:
```go
testMiddleware := func(next HandlerFunc) HandlerFunc {
	return func(ctx Contexter, req Requester) *Response {
		ctx.SetValue("foo", 1)
		return next(ctx, req)
	}
}
```

In the middleware you can manipulate both the request and the context, which then get passed through to your handler.

You can also chain middleware, with them running in chronological order:
```go
aws.Start(New(testHandler).Middleware(
	testMiddleware,
	testMiddlewareTwo,
).Run())
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
type Interceptor func(*handler.Response) *handler.Response
```

In order to manipulate the handler response you can take the `handler.Response` `Body` to change the response, and then add it back to the `handler.Response`, e.g.:

```go
func(r *handler.Response) *handler.Response {
	m := &metasyntactic{}
	err := json.Unmarshal([]byte(r.Body), m)
	assert.NoError(t, err)

	m.Bar = "interceptor 1"

	b, err := json.Marshal(m)
	assert.NoError(t, err)

	r.Body = string(b)

	return r
},
```


## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.


## License
[MIT](https://choosealicense.com/licenses/mit/)
