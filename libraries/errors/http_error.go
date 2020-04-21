package errors

import (
	"encoding/json"
	"net/http"
	"os"
)

var (
	ErrorMode = os.Getenv("HTTP_ERROR")
)

type httpError struct {
	Error      *string `json:"error,omitempty"`
	StackTrace *string `json:"stackTrace,omitempty"`
}

func HandleHTTPError(w http.ResponseWriter, r *http.Request, err Error) {
	msg := err.Text()
	body := &httpError{
		Error: &msg,
	}

	if ErrorMode == "DEBUG" {
		e := err.StackTrace()
		body.StackTrace = &e
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.ErrorCode())
	json.NewEncoder(w).Encode(&body)
}
