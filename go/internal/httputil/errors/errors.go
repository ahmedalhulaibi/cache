package errors

import (
	"encoding/json"
	"net/http"
)

type httpError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func WriteError(w http.ResponseWriter, err httpError) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(err.Status)
	return json.NewEncoder(w).Encode(err)
}

func NotImplementedError() httpError {
	return httpError{http.StatusNotImplemented, "Not implemented"}
}

func InternalServerError(msg string) httpError {
	return httpError{http.StatusInternalServerError, "Internal server error" + msg}
}

func BadRequest(msg string) httpError {
	return httpError{http.StatusBadRequest, "Bad request - " + msg}
}

func Unauthorized() httpError {
	return httpError{http.StatusUnauthorized, "Unauthorized"}
}
