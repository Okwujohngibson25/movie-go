package response

import (
	"encoding/json"
	"net/http"
)

type ErrorDetails struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

type APIResponse[T any] struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Data    T              `json:"data,omitempty"`
	Details []ErrorDetails `json:"details,omitempty"`
}

func WriteJSON(w http.ResponseWriter, statusCode int, response any) error {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(response)
}
