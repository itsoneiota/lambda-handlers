package aws

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ResponseWriterSuite struct {
	suite.Suite
	headers http.Header
}

func (s *ResponseWriterSuite) SetupTest() {
	s.headers = http.Header{
		"Content-Type": []string{
			"application/json; charset=utf-8",
		},
	}
}

func (s *ResponseWriterSuite) TestNewResponseWriter() {
	r := NewResponseWriter(s.headers)

	s.IsType(&ResponseWriter{}, r)
	s.NotEmpty(r.Headers)

	r.Write([]byte("foo"))
	s.NotEmpty(r.Body)

	r.WriteHeader(http.StatusOK)
	s.Equal(http.StatusOK, r.StatusCode)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestResponseWriterSuite(t *testing.T) {
	suite.Run(t, new(ResponseWriterSuite))
}
