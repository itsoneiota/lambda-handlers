package mux

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/slatermorgan/lambda-handlers/pkg/aws"
	"github.com/slatermorgan/lambda-handlers/pkg/handler"
)

type Request struct {
	request *http.Request
}

func NewRequest(r *http.Request) *Request {
	return &Request{
		request: r,
	}
}

// Body gets request payload
func (r *Request) Body() string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.request.Body)

	return buf.String()
}

// HeaderByName gets a header by its name eg. "content-type"
func (r *Request) Headers() http.Header {
	return r.request.Header
}

// MultipartReader is an iterator over parts in a MIME multipart body
func (r *Request) MultipartReader() (*multipart.Reader, error) {
	return r.request.MultipartReader()
}

// PathByName gets a path parameter by its name eg. "productID"
func (r *Request) PathByName(name string) string {
	vars := mux.Vars(r.request)

	return vars[name]
}

// QueryByName gets a query parameter by its name eg. "locale"
func (r *Request) QueryByName(name string) string {
	v := r.request.URL.Query()

	return v.Get(name)
}

// QueryByName gets a query parameter by its name eg. "locale"
func (r *Request) QueryParams() url.Values {
	return r.request.URL.Query()
}

// SetQueryByName gets a query parameter by its name eg. "locale"
func (r *Request) SetQueryByName(name, set string) {
	v := r.request.URL.Query()
	v.Set(name, set)
}

// PathByName gets a query parameter by its name eg. "locale"
func (r *Request) GetAuthToken() string {
	if r.Headers().Get("Authorization") != "" {
		return r.Headers().Get("Authorization")
	} else {
		return r.Headers().Get("authorization")
	}
}

func (r *Request) Context() handler.Contexter {
	return aws.Context{}
}
