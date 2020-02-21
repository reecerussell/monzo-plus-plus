package errors

import (
	"encoding/json"
	"net/http"
)

type httpError struct {
	Error *string `json:"error,omitempty"`
}

func HandleHTTPError(w http.ResponseWriter, r *http.Request, err Error) {
	msg := err.Text()
	body := &httpError{
		Error: &msg,
	}

	w.WriteHeader(err.ErrorCode())
	json.NewEncoder(w).Encode(&body)
}
