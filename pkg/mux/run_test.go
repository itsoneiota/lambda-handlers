package mux

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/itsoneiota/lambda-handlers/v2/pkg/handler"
	"github.com/stretchr/testify/suite"
)

type metasyntactic struct {
	Foo string `json:"foo"`
	Bar string `json:"bar"`
	Baz string `json:"baz"`
}

type RunSuite struct {
	suite.Suite
}

func (s *RunSuite) TestRun() {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		m := &metasyntactic{
			Foo: "foo",
		}

		b, err := json.Marshal(m)
		s.NoError(err)

		w.Write(b)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("foo", "bar")
	}

	resp := &fakeResponseWriter{
		headers: http.Header{},
	}
	req := &http.Request{}
	New(testHandler, handler.WithHeaders(http.Header{
		"default": {"header"},
	})).Run()(resp, req)

	s.Equal(http.StatusOK, resp.statusCode)
	s.Equal(`{"foo":"foo","bar":"","baz":""}`, string(resp.body))
	s.NotEmpty(resp.headers)
	s.Equal("header", resp.headers.Get("default"))
	s.Equal("bar", resp.headers.Get("foo"))
}

func (s *RunSuite) TestMiddlewares() {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		query, err := url.ParseQuery(r.URL.RawQuery)
		s.NoError(err)

		m := &metasyntactic{
			Foo: fmt.Sprintf("%d", r.Context().Value("foo").(int)),
			Bar: query.Get("bar"),
		}

		b, err := json.Marshal(m)
		s.NoError(err)

		w.Write(b)
		w.WriteHeader(http.StatusOK)
	}

	resp := &fakeResponseWriter{
		headers: http.Header{},
	}
	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{},
	}
	New(testHandler, handler.WithMiddlewares(
		func(r *http.Request) (*http.Request, error) {
			ctx := context.WithValue(r.Context(), "foo", 1)

			return r.WithContext(ctx), nil
		},
		func(r *http.Request) (*http.Request, error) {
			query, err := url.ParseQuery(r.URL.RawQuery)
			s.NoError(err)

			query.Set("bar", "2")

			r.URL.RawQuery = query.Encode()

			return r, nil
		},
	)).Run()(resp, req)

	s.Equal(http.StatusOK, resp.statusCode)
	s.Equal(`{"foo":"1","bar":"2","baz":""}`, string(resp.body))
}

func (s *RunSuite) TestInterceptors() {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		m := &metasyntactic{
			Foo: "foo",
		}

		b, err := json.Marshal(m)
		s.NoError(err)

		w.Write(b)
		w.WriteHeader(http.StatusOK)
	}

	resp := &fakeResponseWriter{
		headers: http.Header{},
	}
	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{},
	}
	New(testHandler, handler.WithInterceptors(
		func(_ *http.Request, w handler.ResponseWriter) error {
			m := &metasyntactic{}
			err := json.Unmarshal([]byte(w.Body()), m)
			s.NoError(err)

			m.Bar = "4"

			b, err := json.Marshal(m)
			s.NoError(err)

			w.Write(b)

			return nil
		},
	)).Run()(resp, req)

	s.Equal(http.StatusOK, resp.statusCode)
	s.Equal(`{"foo":"foo","bar":"4","baz":""}`, string(resp.body))
}

func (s *RunSuite) TestInterceptorsRequest() {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		m := &metasyntactic{
			Foo: "foo",
		}

		b, err := json.Marshal(m)
		s.NoError(err)

		w.Write(b)
		w.WriteHeader(http.StatusOK)
	}

	resp := &fakeResponseWriter{
		headers: http.Header{},
	}
	req := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			RawQuery: "bar=2",
		},
	}
	New(testHandler, handler.WithInterceptors(
		func(req *http.Request, w handler.ResponseWriter) error {
			m := &metasyntactic{}
			err := json.Unmarshal([]byte(w.Body()), m)
			s.NoError(err)

			query, err := url.ParseQuery(req.URL.RawQuery)
			s.NoError(err)

			m.Bar = query.Get("bar")

			b, err := json.Marshal(m)
			s.NoError(err)

			w.Write(b)

			return nil
		},
	)).Run()(resp, req)

	s.Equal(http.StatusOK, resp.statusCode)
	s.Equal(`{"foo":"foo","bar":"2","baz":""}`, string(resp.body))
}

func (s *RunSuite) TestInterceptorsHeaders() {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		m := &metasyntactic{
			Foo: "foo",
		}

		b, err := json.Marshal(m)
		s.NoError(err)

		w.Write(b)
		w.WriteHeader(http.StatusOK)
	}

	resp := &fakeResponseWriter{
		headers: http.Header{},
	}
	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{},
	}
	New(testHandler, handler.WithInterceptors(
		func(_ *http.Request, w handler.ResponseWriter) error {
			w.Header().Add("foo", "bar")

			return nil
		},
	)).Run()(resp, req)

	s.Equal(http.StatusOK, resp.statusCode)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestRunSuite(t *testing.T) {
	suite.Run(t, new(RunSuite))
}
