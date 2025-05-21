package main

import (
	"net/http"
)

// This file contains alternative API utility functions.
// NOTE: The main application uses bodyUnmarshal and bodyMarshal from utils_http.go.
// These functions are provided for compatibility or future use.

// unmarshalBody parses the request body into the provided struct
// This is an alternative implementation to bodyUnmarshal in utils_http.go
func unmarshalBody(r *http.Request, x interface{}) error {
	return bodyUnmarshal(r, x)
}

// writeAPIResponse writes an API response with the given error message and data
// This is an alternative implementation to bodyMarshal in utils_http.go
func writeAPIResponse(w http.ResponseWriter, success bool, message string, data interface{}) {
	resp := map[string]interface{}{
		"success": success,
		"message": message,
	}
	
	if data != nil {
		resp["data"] = data
	}
	
	bodyMarshal(w, resp)
}
