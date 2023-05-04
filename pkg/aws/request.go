package aws

import (
	"bytes"
	"encoding/base64"
	"errors"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

var (
	ErrContentTypeHeaderMissing         = errors.New("content type header missing")
	ErrContentTypeHeaderNotMultipart    = errors.New("content type header not multipart error")
	ErrContentTypeHeaderMissingBoundary = errors.New("content type header missing boundary error")
)

type AWSRequest struct {
	body            string
	pathParams      map[string]string
	queryParams     map[string]string
	headers         http.Header
	isBase64Encoded bool
}

func NewAWSRequest(r *events.APIGatewayProxyRequest) *AWSRequest {
	headers := http.Header{}
	for k, v := range r.Headers {
		headers.Set(k, v)
	}

	return &AWSRequest{
		body:            r.Body,
		pathParams:      r.PathParameters,
		queryParams:     r.QueryStringParameters,
		headers:         headers,
		isBase64Encoded: r.IsBase64Encoded,
	}
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
func (r *AWSRequest) QueryByName(name string) string {
	return r.queryParams[name]
}

// PathByName sets a query parameter by its name eg. "locale"
// This is used to alter requests in middleware functions.
func (r *AWSRequest) SetQueryByName(name, set string) {
	r.queryParams[name] = set
}

// PathByName gets a query parameter by its name eg. "locale"
func (r *AWSRequest) GetAuthToken() string {
	if r.Headers().Get("Authorization") != "" {
		return r.Headers().Get("Authorization")
	} else {
		return r.Headers().Get("authorization")
	}
}
