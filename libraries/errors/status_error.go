package errors

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
)

type statusError struct {
	Traceable
	err        error
	statusCode int
}

func (se *statusError) Text() string {
	return se.err.Error()
}

func (se *statusError) ErrorCode() int {
	return se.statusCode
}

func new(e interface{}, code int) *statusError {
	var err error

	switch v := e.(type) {
	case error:
		err = v
		break
	default:
		err = fmt.Errorf("%v", e)
	}

	stack := make([]uintptr, 50)
	length := runtime.Callers(2, stack[:])

	se := &statusError{
		err:        err,
		statusCode: code,
	}
	se.Stack = stack[:length]

	return se
}

func InternalError(err error) Error {
	return new(err, http.StatusInternalServerError)
}

func BadRequest(err string) Error {
	return new(err, http.StatusBadRequest)
}

func NotFound(err string) Error {
	return new(err, http.StatusNotFound)
}

func Unauthorised(err string) Error {
	return new(err, http.StatusUnauthorized)
}

func Forbidden() Error {
	return new(errors.New("insufficient permissions"), http.StatusForbidden)
}
