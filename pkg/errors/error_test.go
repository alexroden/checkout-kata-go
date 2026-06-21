package errors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ErrorSuite struct {
	suite.Suite
}

func (s *ErrorSuite) Test_New() {
	msg := "required field"
	err := New(CodeBadRequest, msg)
	s.Error(err)

	s.Equal(msg, err.Error())
	s.Equal(http.StatusBadRequest, err.StatusCode())
}

func (s *ErrorSuite) Test_NewFromError() {
	msg := "something went wrong"
	err := NewFromErr(errors.New(msg), "")
	s.Error(err)

	s.Equal(msg, err.Error())
	s.Equal(http.StatusInternalServerError, err.StatusCode())
}

func (s *ErrorSuite) Test_InternalServerError() {
	msg := "something went wrong"
	err := InternalServerError(msg)
	s.Error(err)

	s.Equal(msg, err.Error())
	s.Equal(http.StatusInternalServerError, err.StatusCode())
}

func (s *ErrorSuite) Test_NotImplemented() {
	msg := "type not implemented"
	err := NotImplemented(msg)
	s.Error(err)

	s.Equal(msg, err.Error())
	s.Equal(http.StatusNotImplemented, err.StatusCode())
}

func (s *ErrorSuite) Test_UnprocessableEntity() {
	msg := `email address does not contain "@"`
	err := UnprocessableEntity(msg)
	s.Error(err)

	s.Equal(msg, err.Error())
	s.Equal(http.StatusUnprocessableEntity, err.StatusCode())
}

func (s *ErrorSuite) Test_Conflict() {
	msg := "version already in use"
	err := Conflict(msg)
	s.Error(err)

	s.Equal(msg, err.Error())
	s.Equal(http.StatusConflict, err.StatusCode())
}

func (s *ErrorSuite) Test_RequestTimeout() {
	msg := "response exceeded 30 seconds"
	err := RequestTimeout(msg)
	s.Error(err)

	s.Equal(msg, err.Error())
	s.Equal(http.StatusRequestTimeout, err.StatusCode())
}

func (s *ErrorSuite) Test_NotFound() {
	msg := "user not found"
	err := NotFound(msg)
	s.Error(err)

	s.Equal(msg, err.Error())
	s.Equal(http.StatusNotFound, err.StatusCode())
}

func (s *ErrorSuite) Test_Forbidden() {
	msg := "forbidden action"
	err := Forbidden(msg)
	s.Error(err)

	s.Equal(msg, err.Error())
	s.Equal(http.StatusForbidden, err.StatusCode())
}

func (s *ErrorSuite) Test_Unauthorized() {
	msg := "unauthorized"
	err := Unauthorized(msg)
	s.Error(err)

	s.Equal(msg, err.Error())
	s.Equal(http.StatusUnauthorized, err.StatusCode())
}

func (s *ErrorSuite) Test_BadRequest() {
	msg := "field required"
	err := BadRequest(msg)
	s.Error(err)

	s.Equal(msg, err.Error())
	s.Equal(http.StatusBadRequest, err.StatusCode())
}

func (s *ErrorSuite) Test_Found() {
	msg := "temporary url change"
	err := Found(msg)
	s.Error(err)

	s.Equal(msg, err.Error())
	s.Equal(http.StatusFound, err.StatusCode())
}

func (s *ErrorSuite) Test_MovedPermanently() {
	msg := "url permanently moved"
	err := MovedPermanently(msg)
	s.Error(err)

	s.Equal(msg, err.Error())
	s.Equal(http.StatusMovedPermanently, err.StatusCode())
}

func (s *ErrorSuite) Test_BadGateway() {
	msg := "server overload"
	err := BadGateway(msg)
	s.Error(err)

	s.Equal(msg, err.Error())
	s.Equal(http.StatusBadGateway, err.StatusCode())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestErrorSuite(t *testing.T) {
	suite.Run(t, new(ErrorSuite))
}
