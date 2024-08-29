package aws

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
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

func NewHttpRequest(r *events.APIGatewayProxyRequest) (*http.Request, error) {
	form, err := getMultipartForm(r, 32<<20) // 32 MB max memory
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: r.HTTPMethod,
		URL: &url.URL{
			Scheme:   "https",
			Path:     r.Path,
			RawPath:  getRawPath(r, true),
			RawQuery: getRawPath(r, false),
		},
		Header:        getHeaders(r),
		Body:          getBody(r),
		MultipartForm: form,
	}

	ctx := context.Background()
	req.WithContext(ctx)

	return req, nil
}

func getBody(r *events.APIGatewayProxyRequest) io.ReadCloser {
	result := io.NopCloser(strings.NewReader(r.Body))
	result.Close()

	return result
}

func getHeaders(r *events.APIGatewayProxyRequest) http.Header {
	result := http.Header{}
	for k, v := range r.Headers {
		result.Add(k, v)
	}

	return result
}

func getMultipartForm(r *events.APIGatewayProxyRequest, maxMemory int64) (*multipart.Form, error) {
	ct := r.Headers["content-type"]
	if len(ct) == 0 {
		ct = r.Headers["Content-Type"]
		if len(ct) == 0 {
			return nil, ErrContentTypeHeaderMissing
		}
	}

	if !strings.HasPrefix(ct, "multipart/") {
		return nil, nil
	}

	mediatype, params, err := mime.ParseMediaType(ct)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(mediatype)), "multipart/") {
		return nil, ErrContentTypeHeaderNotMultipart
	}

	boundary, ok := params["boundary"]
	if !ok {
		return nil, ErrContentTypeHeaderMissingBoundary
	}

	var reader io.Reader
	if r.IsBase64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(r.Body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(decoded)
	} else {
		reader = strings.NewReader(r.Body)
	}

	multipartReader := multipart.NewReader(reader, boundary)

	form, err := multipartReader.ReadForm(maxMemory)
	if err != nil {
		return nil, err
	}

	return form, nil
}

func getRawPath(r *events.APIGatewayProxyRequest, hasPath bool) string {
	var result string
	if hasPath {
		result = r.Path
	}

	q := url.Values{}
	for k, v := range r.QueryStringParameters {
		q.Add(k, v)
	}

	for k, vals := range r.MultiValueQueryStringParameters {
		for _, v := range vals {
			q.Add(k, v)
		}
	}

	if len(q) > 0 {
		if hasPath {
			result = fmt.Sprintf("%s?", result)
		}
		result = fmt.Sprintf("%s%s", result, q.Encode())
	}

	return result
}

// type AWSRequest struct {
// 	body            string
// 	context         handler.Contexter
// 	cookies         []*http.Cookie
// 	headers         http.Header
// 	isBase64Encoded bool
// 	pathParams      map[string]string
// 	queryParams     url.Values
// }

// func NewAWSRequest(r *events.APIGatewayProxyRequest) *AWSRequest {
// 	headers := http.Header{}
// 	for k, v := range r.Headers {
// 		headers.Set(k, v)
// 	}

// 	values := url.Values{}
// 	for k, v := range r.QueryStringParameters {
// 		values.Set(k, v)
// 	}

// 	cookies := []*http.Cookie{}
// 	c := headers.Get("Cookie")
// 	if c == "" {
// 		c = headers.Get("cookie")
// 	}

// 	if c != "" {
// 		for _, cookie := range strings.Split(";", c) {
// 			if s := strings.Split("=", cookie); len(s) > 1 {
// 				cookies = append(cookies, &http.Cookie{
// 					Name:  s[0],
// 					Value: s[1],
// 				})
// 			}
// 		}
// 	}

// 	return &AWSRequest{
// 		body: r.Body,
// 		context: &Context{
// 			APIGatewayProxyRequestContext: r.RequestContext,
// 		},
// 		cookies:         cookies,
// 		isBase64Encoded: r.IsBase64Encoded,
// 		headers:         headers,
// 		pathParams:      r.PathParameters,
// 		queryParams:     values,
// 	}
// }

// // Add cookie
// func (r *AWSRequest) AddCookie(c *http.Cookie) {
// 	r.cookies = append(r.cookies, c)
// }

// // Body gets request payload
// func (r *AWSRequest) Body() string {
// 	return r.body
// }

// // Get context
// func (r *AWSRequest) Context() handler.Contexter {
// 	return r.context
// }

// // Get cookie
// func (r *AWSRequest) Cookie(name string) (*http.Cookie, error) {
// 	var result *http.Cookie
// 	for _, c := range r.cookies {
// 		if c.Name == name {
// 			result = c
// 			break
// 		}
// 	}

// 	return result, nil
// }

// // Get cookies
// func (r *AWSRequest) Cookies() []*http.Cookie {
// 	return r.cookies
// }

// // Get auth token from headers.
// func (r *AWSRequest) GetAuthToken() string {
// 	if v := r.Headers().Get("Authorization"); v != "" {
// 		return v
// 	}

// 	return r.Headers().Get("authorization")
// }

// // Headers get the request headers
// func (r *AWSRequest) Headers() http.Header {
// 	return r.headers
// }

// // MultipartReader is an iterator over parts in a MIME multipart body
// func (r *AWSRequest) MultipartReader() (*multipart.Reader, error) {
// 	ct := r.headers.Get("content-type")
// 	if len(ct) == 0 {
// 		ct = r.headers.Get("Content-Type")
// 		if len(ct) == 0 {
// 			return nil, ErrContentTypeHeaderMissing
// 		}
// 	}

// 	mediatype, params, err := mime.ParseMediaType(ct)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if strings.Index(strings.ToLower(strings.TrimSpace(mediatype)), "multipart/") != 0 {
// 		return nil, ErrContentTypeHeaderNotMultipart
// 	}

// 	boundary, ok := params["boundary"]
// 	if !ok {
// 		return nil, ErrContentTypeHeaderMissingBoundary
// 	}

// 	if r.isBase64Encoded {
// 		decoded, err := base64.StdEncoding.DecodeString(r.body)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return multipart.NewReader(bytes.NewReader(decoded), boundary), nil
// 	}

// 	return multipart.NewReader(strings.NewReader(r.body), boundary), nil
// }

// // PathByName gets a path parameter by its name eg. "productID"
// func (r *AWSRequest) PathByName(name string) string {
// 	return r.pathParams[name]
// }

// // QueryByName gets a query parameter by its name eg. "locale"
// func (r *AWSRequest) QueryByName(name string) string {
// 	return r.queryParams.Get(name)
// }

// // QueryByName gets all query parameters
// func (r *AWSRequest) QueryParams() url.Values {
// 	return r.queryParams
// }

// // Get referer
// func (r *AWSRequest) Referer() string {
// 	if v := r.Headers().Get("Referer"); v != "" {
// 		return v
// 	}

// 	return r.Headers().Get("referer")
// }

// // Sets a query parameter against the request.
// func (r *AWSRequest) SetQueryByName(name, set string) {
// 	r.queryParams.Set(name, set)
// }

// // Get user agent
// func (r *AWSRequest) UserAgent() string {
// 	if v := r.Headers().Get("User-Agent"); v != "" {
// 		return v
// 	}

// 	return r.Headers().Get("user-agent")
// }
