package main

import (
	"log"
	"log/slog"
	"net/http"

	muxRouter "github.com/gorilla/mux"
	"github.com/itsoneiota/lambda-handlers/internal/mocks"
	"github.com/itsoneiota/lambda-handlers/pkg/example"
	"github.com/itsoneiota/lambda-handlers/pkg/handler"
	"github.com/itsoneiota/lambda-handlers/pkg/mux"
	"github.com/itsoneiota/lambda-handlers/pkg/test"
)

func main() {
	expectToken := "authToken"
	model := example.ExampleModel{
		Success: true,
	}

	expectQuery := "M36FJ"

	c := new(mocks.Connector)
	c.On("Authorize",
		expectToken,
	).Return(
		nil,
	).Times(1)

	c.On("Find",
		expectQuery,
	).Return(
		model,
		nil,
	).Times(1)

	l := test.NewNullLogger()
	slog.SetDefault(l)

	resHander := handler.NewResponseHandler(http.Header{})

	r := muxRouter.NewRouter()
	r.HandleFunc("/test", mux.CreateHandler(example.FindHandler(resHander, c, nil, nil)))

	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
