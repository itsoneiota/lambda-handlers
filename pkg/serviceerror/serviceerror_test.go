package serviceerror

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceErrorSuite struct {
	suite.Suite
}

func (s *ServiceErrorSuite) TestCanGetErrorMessage() {
	tests := []struct {
		errorCode string
		message   string
	}{
		{"NOT_FOUND", "not found message"},
		{"UNKNOWN", "unknown error"},
		{"SERVER_ERROR", "an unknown server has occurred"},
	}
	for _, test := range tests {
		e := NewServiceError(test.errorCode, test.errorCode, test.message)
		s.Equal(test.message, e.Error(), "Message returned did not match expected")
		s.Equal(test.errorCode, e.Code(), "Code returned does did not match expected")
	}
}

func (s *ServiceErrorSuite) TestGetStatusCode() {
	tests := []struct {
		errorCode  string
		statusCode int
	}{
		{CodeInternalServerError, http.StatusInternalServerError},
		{CodeNotFound, http.StatusNotFound},
		{CodeBadRequest, http.StatusBadRequest},
		{CodeNotImplemented, http.StatusNotImplemented},
		{CodeUnauthorized, http.StatusUnauthorized},
		{CodeForbidden, http.StatusForbidden},
		{CodeRequestTimeout, http.StatusRequestTimeout},
		{CodeConflict, http.StatusConflict},
	}
	err := NewServiceError("", "", "")
	for _, test := range tests {
		err.Err.Code = test.errorCode
		s.Equal(test.statusCode, err.StatusCode(), "Status code does not match expected")
	}
}

func (s *ServiceErrorSuite) TestCanCreateNewFromExistingError() {
	// New from regular error
	err := errors.New("an error occurred")
	sut := NewFromErr(err, "error")
	expected := "error: an error occurred"
	s.Equal(expected, sut.Error(), "Error message did not match expected")
	s.Equal(sut.Code(), CodeInternalServerError, "Error code did not match expected")
	// New from service error
	err = NewServiceError(CodeNotFound, CodeNotFound, "not found")
	sut = NewFromErr(err, "error")
	expected = "error: not found"
	s.Equal(expected, sut.Error(), "Error message did not match expected")
	s.Equal(sut.Code(), CodeNotFound, "Error code did not match expected")
}

func (s *ServiceErrorSuite) TestConvenienceMethods() {
	tests := []struct {
		method       func(string) *ServiceError
		expectedType string
	}{
		{InternalServerError, CodeInternalServerError},
		{NotImplemented, CodeNotImplemented},
		{UnprocessableEntity, CodeUnprocessableEntity},
		{Conflict, CodeConflict},
		{RequestTimeout, CodeRequestTimeout},
		{NotFound, CodeNotFound},
		{Forbidden, CodeForbidden},
		{Unauthorized, CodeUnauthorized},
		{BadRequest, CodeBadRequest},
		{Found, CodeFound},
		{MovedPermanently, CodeMovedPermanently},
	}

	for _, test := range tests {
		e := test.method("borken")
		s.Equal(e.Code(), test.expectedType)
		s.Equal(e.Error(), "borken")
	}
}

func (s *ServiceErrorSuite) TestGetDefaultErrorMessage() {
	tests := []struct {
		method       func(string) *ServiceError
		expectedType string
	}{
		{InternalServerError, CodeInternalServerError},
		{NotImplemented, CodeNotImplemented},
		{UnprocessableEntity, CodeUnprocessableEntity},
		{Conflict, CodeConflict},
		{RequestTimeout, CodeRequestTimeout},
		{NotFound, CodeNotFound},
		{Forbidden, CodeForbidden},
		{Unauthorized, CodeUnauthorized},
		{BadRequest, CodeBadRequest},
		{Found, CodeFound},
		{MovedPermanently, CodeMovedPermanently},
	}

	for _, test := range tests {
		e := test.method("borken")
		s.Equal(e.Code(), test.expectedType)
		s.Equal(e.Error(), "borken")
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestServiceErrorSuite(t *testing.T) {
	suite.Run(t, new(ServiceErrorSuite))
}
