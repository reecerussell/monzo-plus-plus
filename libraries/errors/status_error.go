package errors

import (
	"errors"
	"net/http"
)

type statusError struct {
	err        error
	statusCode int
}

func (se *statusError) Text() string {
	return se.err.Error()
}

func (se *statusError) ErrorCode() int {
	return se.statusCode
}

func InternalError(err error) Error {
	return &statusError{err, http.StatusInternalServerError}
}

func BadRequest(err string) Error {
	return &statusError{errors.New(err), http.StatusBadRequest}
}

func NotFound(err string) Error {
	return &statusError{errors.New(err), http.StatusNotFound}
}

func Unauthorised(err string) Error {
	return &statusError{errors.New(err), http.StatusUnauthorized}
}

func Forbidden(err string) Error {
	return &statusError{errors.New(err), http.StatusForbidden}
}
