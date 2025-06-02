package aws

import (
	"bytes"
	"encoding/base64"
	"errors"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

var (
	ErrContentTypeHeaderMissing         = errors.New("content type header missing")
	ErrContentTypeHeaderNotMultipart    = errors.New("content type header not multipart error")
	ErrContentTypeHeaderMissingBoundary = errors.New("content type header missing boundary error")
)

type AWSRequest struct {
	path            string
	method          string
	body            string
	pathParams      map[string]string
	queryParams     url.Values
	headers         http.Header
	isBase64Encoded bool
}

func NewAWSRequest(r *events.APIGatewayProxyRequest) *AWSRequest {
	headers := http.Header{}
	for k, v := range r.Headers {
		headers.Set(k, v)
	}

	values := url.Values{}
	for k, v := range r.QueryStringParameters {
		values.Set(k, v)
	}

	return &AWSRequest{
		path:            r.Path,
		method:          r.HTTPMethod,
		body:            r.Body,
		pathParams:      r.PathParameters,
		queryParams:     values,
		headers:         headers,
		isBase64Encoded: r.IsBase64Encoded,
	}
}

// Request path
func (r *AWSRequest) Path() string {
	return r.path
}

// Request method
func (r *AWSRequest) Method() string {
	return r.method
}

// Body gets request payload
func (r *AWSRequest) Body() string {
	return r.body
}

// Headers get the request headers
func (r *AWSRequest) Headers() http.Header {
	return r.headers
}

// MultipartReader is an iterator over parts in a MIME multipart body
func (r *AWSRequest) MultipartReader() (*multipart.Reader, error) {
	ct := r.headers.Get("content-type")
	if len(ct) == 0 {
		ct = r.headers.Get("Content-Type")
		if len(ct) == 0 {
			return nil, ErrContentTypeHeaderMissing
		}
	}

	mediatype, params, err := mime.ParseMediaType(ct)
	if err != nil {
		return nil, err
	}

	if strings.Index(strings.ToLower(strings.TrimSpace(mediatype)), "multipart/") != 0 {
		return nil, ErrContentTypeHeaderNotMultipart
	}

	boundary, ok := params["boundary"]
	if !ok {
		return nil, ErrContentTypeHeaderMissingBoundary
	}

	if r.isBase64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(r.body)
		if err != nil {
			return nil, err
		}
		return multipart.NewReader(bytes.NewReader(decoded), boundary), nil
	}

	return multipart.NewReader(strings.NewReader(r.body), boundary), nil
}

// PathByName gets a path parameter by its name eg. "productID"
func (r *AWSRequest) PathByName(name string) string {
	return r.pathParams[name]
}

// QueryByName gets a query parameter by its name eg. "locale"
func (r *AWSRequest) QueryByName(names ...string) string {
	var result string
	for _, name := range names {
		if result != "" {
			break
		}

		result = r.queryParams.Get(name)
	}

	return result
}

// QueryByName gets a query parameter by its name eg. "locale"
func (r *AWSRequest) QueryParams() url.Values {
	return r.queryParams
}

// PathByName sets a query parameter by its name eg. "locale"
// This is used to alter requests in middleware functions.
func (r *AWSRequest) SetQueryByName(name, set string) {
	r.queryParams.Set(name, set)
}

// PathByName gets a query parameter by its name eg. "locale"
func (r *AWSRequest) GetAuthToken() string {
	if r.Headers().Get("Authorization") != "" {
		return r.Headers().Get("Authorization")
	} else {
		return r.Headers().Get("authorization")
	}
}
