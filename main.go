package main

import (
	"log"
	"net/http"

	muxRouter "github.com/gorilla/mux"
	"github.com/itsoneiota/lambda-handlers/v2/internal/mocks"
	"github.com/itsoneiota/lambda-handlers/v2/pkg/example"
	"github.com/itsoneiota/lambda-handlers/v2/pkg/handler"
	"github.com/itsoneiota/lambda-handlers/v2/pkg/mux"
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
	r.HandleFunc("/test", mux.CreateHandler(example.FindHandler(resHander, c)))

	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
