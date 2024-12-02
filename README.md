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

type AfterFindHandlerHook func(interface{}) error
func FindHandler(
	connector Connector,
	beforeHook handler.BeforeHandlerHook,
	afterHook AfterFindHandlerHook,
) handler.HandlerFunc {
	return func(ctx handler.Contexter, request handler.Requester) *handler.Response {
		if beforeHook != nil {
			if err := beforeHook(request); err != nil {
				return handler.NewErrorResponse(err)
			}
		}

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

## Middleware

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

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.


## License
[MIT](https://choosealicense.com/licenses/mit/)
