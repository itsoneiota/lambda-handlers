package aws

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HandlerSuite struct {
	suite.Suite
	headers http.Header
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
