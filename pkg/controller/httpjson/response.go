package httpjson

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error   string      `json:"error"`
	Code    string      `json:"code"`
	Details interface{} `json:"details,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, code, message string, details interface{}) {
	WriteJSON(w, status, ErrorResponse{Error: message, Code: code, Details: details})
}

func BadRequest(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusBadRequest, "BAD_REQUEST", message, nil)
}

func Unauthorized(w http.ResponseWriter, code, message string) {
	WriteError(w, http.StatusUnauthorized, code, message, nil)
}

func Forbidden(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusForbidden, "FORBIDDEN", message, nil)
}

func Internal(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", message, nil)
}
