package aws

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceErrorSuite struct {
	suite.Suite
}

func (s *ServiceErrorSuite) TestNotFoundError() {
	code := CodeNotFound
	message := "not found message"

	e := NewServiceError(code, code, message)
	s.Equal(message, e.Error())
	s.Equal(code, e.Code())
}

func (s *ServiceErrorSuite) TestUnknownError() {
	code := CodeUnknown
	message := "unknown error"

	e := NewServiceError(code, code, message)
	s.Equal(message, e.Error())
	s.Equal(code, e.Code())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestServiceErrorSuite(t *testing.T) {
	suite.Run(t, new(ServiceErrorSuite))
}
