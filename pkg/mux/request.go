package mux

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/itsoneiota/lambda-handlers/pkg/handler"
)

type Request struct {
	request *http.Request
}

func NewRequest(r *http.Request) *Request {
	return &Request{
		request: r,
	}
}

// Add cookie
func (r *Request) AddCookie(c *http.Cookie) {
	cookies := r.Headers().Get("Cookie")
	if cookies == "" {
		cookies = r.Headers().Get("cookie")
	}

	if cookies != "" {
		cookies = fmt.Sprintf("%s;", cookies)
	}

	cookies = fmt.Sprintf("%s %s=%s", cookies, c.Name, c.Value)

	r.Headers().Set("Cookie", cookies)
}

// Body gets request payload
func (r *Request) Body() string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.request.Body)

	return buf.String()
}

// Get context
func (r *Request) Context() handler.Contexter {
	return &Context{
		Request: r.request,
	}
}

// Get cookie
func (r *Request) Cookie(name string) (*http.Cookie, error) {
	var result *http.Cookie
	for _, cookie := range r.Cookies() {
		if cookie.Name == name {
			result = cookie
			break
		}
	}

	return result, nil
}

// Get cookies
func (r *Request) Cookies() []*http.Cookie {
	var result []*http.Cookie
	cookies := r.Headers().Get("Cookie")
	if cookies == "" {
		cookies = r.Headers().Get("cookie")
	}

	for _, cookie := range strings.Split(";", cookies) {
		if s := strings.Split("=", cookie); len(s) > 1 {
			result = append(result, &http.Cookie{
				Name:  s[0],
				Value: s[1],
			})
		}
	}

	return result
}

// Get auth token
func (r *Request) GetAuthToken() string {
	if v := r.Headers().Get("Authorization"); v != "" {
		return v
	}

	return r.Headers().Get("authorization")
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

// Get referer
func (r *Request) Referer() string {
	if v := r.Headers().Get("Referer"); v != "" {
		return v
	}

	return r.Headers().Get("referer")
}

// SetQueryByName gets a query parameter by its name eg. "locale"
func (r *Request) SetQueryByName(name, set string) {
	v := r.request.URL.Query()
	v.Set(name, set)
}

// Get user agent
func (r *Request) UserAgent() string {
	if v := r.Headers().Get("User-Agent"); v != "" {
		return v
	}

	return r.Headers().Get("user-agent")
}
