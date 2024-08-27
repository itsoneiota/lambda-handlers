package main

import (
	"log"
	"net/http"

	muxRouter "github.com/gorilla/mux"
	"github.com/itsoneiota/lambda-handlers/internal/mocks"
	"github.com/itsoneiota/lambda-handlers/pkg/example"
	"github.com/itsoneiota/lambda-handlers/pkg/handler"
	"github.com/itsoneiota/lambda-handlers/pkg/mux"
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

	resHander := handler.NewResponseHandler()

	r := muxRouter.NewRouter()
	r.HandleFunc("/test", mux.CreateHandler(example.FindHandler(resHander, c, nil, nil)))

	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
