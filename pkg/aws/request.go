package aws

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log/slog"
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

func NewHttpRequest(
	ctx context.Context,
	r *events.APIGatewayProxyRequest,
) (*http.Request, error) {
	scheme := "https"
	if v, ok := r.Headers["X-Forwarded-Proto"]; ok {
		scheme = v
	}

	host := "example.com"
	if v, ok := r.Headers["Host"]; ok {
		host = v
	}

	parsedUrl, err := url.Parse(fmt.Sprintf("%s://%s%s", scheme, host, r.Path))
	if err != nil {
		return nil, err
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

	parsedUrl.RawQuery = q.Encode()

	var body io.ReadCloser
	if r.IsBase64Encoded {
		decodedBody, err := base64.StdEncoding.DecodeString(r.Body)
		if err != nil {
			return nil, err
		}
		body = io.NopCloser(bytes.NewReader(decodedBody))
	} else {
		body = io.NopCloser(strings.NewReader(r.Body))
	}

	req, err := http.NewRequest(r.HTTPMethod, parsedUrl.String(), body)
	if err != nil {
		return nil, err
	}

	for key, values := range r.MultiValueHeaders {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	if r.RequestContext.Identity.SourceIP != "" {
		req.RemoteAddr = r.RequestContext.Identity.SourceIP
		if v, ok := r.Headers["X-Forwarded-Port"]; ok {
			req.RemoteAddr = fmt.Sprintf("%s:%s", req.RemoteAddr, v)
		}
	}

	if userAgent := r.RequestContext.Identity.UserAgent; userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}

	contentType := req.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		mediaType, params, err := mime.ParseMediaType(contentType)
		if err != nil {
			return nil, err
		}

		if strings.HasPrefix(mediaType, "multipart/") {
			mr := multipart.NewReader(body, params["boundary"])

			multipartForm, err := mr.ReadForm(10 << 20) // 10MB max memory for the form
			if err != nil {
				return nil, err
			}

			req.MultipartForm = multipartForm
		}
	}

	req.WithContext(ctx)

	req.Context()

	slog.Debug(fmt.Sprintf("%v", req.Context()))

	return req, nil
}
