package aws

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gorilla/mux"
	"github.com/itsoneiota/lambda-handlers/v2/pkg/fakers"
	"github.com/itsoneiota/lambda-handlers/v2/pkg/handler"
	"github.com/itsoneiota/lambda-handlers/v2/pkg/serviceerror"
	"github.com/stretchr/testify/suite"
)

type metasyntactic struct {
	Foo string `json:"foo"`
	Bar string `json:"bar"`
	Baz string `json:"baz"`
}

type HandlerSuite struct {
	suite.Suite
	headers http.Header
	handler http.HandlerFunc
	req     *events.APIGatewayProxyRequest
	resp    *fakers.ReponseWriter
}

func (s *HandlerSuite) SetupTest() {
	s.headers = http.Header{
		"Content-Type": []string{
			"application/json; charset=utf-8",
		},
		"Test": []string{
			"foo",
			"bar",
			"bar",
		},
	}

	s.handler = func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		query, err := url.ParseQuery(r.URL.RawQuery)
		s.NoError(err)

		m := &metasyntactic{
			Foo: params["foo"],
			Bar: query.Get("bar"),
		}

		if v := r.Context().Value("baz"); v != nil {
			m.Baz = v.(string)
		}

		b, err := json.Marshal(m)
		s.NoError(err)

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}

	s.req = &events.APIGatewayProxyRequest{
		Resource: "/api/foo/{foo}",
		Path:     "/api/foo/1",
		PathParameters: map[string]string{
			"foo": "1",
		},
		QueryStringParameters: map[string]string{
			"bar": "2",
		},
	}
}

func (s *HandlerSuite) TestHandle() {
	resp, err := handle(&Handler{
		function: s.handler,
		BaseOpt:  &handler.BaseOpt{},
	})(s.req)
	s.NoError(err)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.JSONEq(`{"foo":"1","bar":"2","baz":""}`, resp.Body)
}

func (s *HandlerSuite) TestHeaders() {
	b := &handler.BaseOpt{}
	b.SetHeaders(http.Header{
		"foo": {
			"bar",
		},
	})

	resp, err := handle(&Handler{
		function: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("bar", "baz")
		},
		BaseOpt: b,
	})(s.req)
	s.NoError(err)

	s.NotEmpty(resp.Headers)
	s.Equal("bar", resp.Headers["foo"])
	s.Equal("baz", resp.Headers["Bar"])
}

func (s *HandlerSuite) TestMiddlewares() {
	b := &handler.BaseOpt{}
	b.SetMiddlewares([]handler.Middleware{
		func(r *http.Request) (*http.Request, error) {
			query, err := url.ParseQuery(r.URL.RawQuery)
			s.NoError(err)

			query.Set("bar", "3")

			r.URL.RawQuery = query.Encode()

			return r, nil
		},
		func(r *http.Request) (*http.Request, error) {
			ctx := context.WithValue(r.Context(), "baz", "4")

			return r.WithContext(ctx), nil
		},
	})

	resp, err := handle(&Handler{
		function: s.handler,
		BaseOpt:  b,
	})(s.req)
	s.NoError(err)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.JSONEq(`{"foo":"1","bar":"3","baz":"4"}`, resp.Body)
}

func (s *HandlerSuite) TestMiddlewareError() {
	b := &handler.BaseOpt{}
	b.SetMiddlewares([]handler.Middleware{
		func(r *http.Request) (*http.Request, error) {
			return r, serviceerror.BadRequest("something bad has happened")
		},
	})

	resp, err := handle(&Handler{
		function: s.handler,
		BaseOpt:  b,
	})(s.req)
	s.NoError(err)

	s.Equal(http.StatusBadRequest, resp.StatusCode)
	s.JSONEq(`{"error":{"id":"BAD_REQUEST","code":"BAD_REQUEST","message":"something bad has happened"}}`, resp.Body)
}

func (s *HandlerSuite) TestInterceptors() {
	b := &handler.BaseOpt{}
	b.SetInterceptors([]handler.Interceptor{
		func(w handler.ResponseWriter) error {
			m := &metasyntactic{}
			err := json.Unmarshal([]byte(w.Body()), m)
			s.NoError(err)

			m.Bar = "4"

			b, err := json.Marshal(m)
			s.NoError(err)

			w.Write(b)

			return nil
		},
	})

	resp, err := handle(&Handler{
		function: s.handler,
		BaseOpt:  b,
	})(s.req)
	s.NoError(err)

	s.Equal(http.StatusOK, resp.StatusCode)
	s.JSONEq(`{"foo":"1","bar":"4","baz":""}`, resp.Body)
}

func (s *HandlerSuite) TestInterceptorError() {
	b := &handler.BaseOpt{}
	b.SetInterceptors([]handler.Interceptor{
		func(_ handler.ResponseWriter) error {
			return serviceerror.BadRequest("something bad has happened")
		},
	})

	resp, err := handle(&Handler{
		function: s.handler,
		BaseOpt:  b,
	})(s.req)
	s.NoError(err)

	s.Equal(http.StatusBadRequest, resp.StatusCode)
	s.JSONEq(`{"error":{"id":"BAD_REQUEST","code":"BAD_REQUEST","message":"something bad has happened"}}`, resp.Body)
}

func (s *HandlerSuite) TestEncodeHeaders() {
	expect := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
		"Test":         "foo; bar",
	}

	s.Equal(expect, encodeHeaders(s.headers))
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
