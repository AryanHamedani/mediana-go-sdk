package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode int
	Message    string
	Details    interface{}
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (%d): %s", e.StatusCode, e.Message)
}

func ParseError(resp *http.Response) error {
	var errorResponse struct {
		Message string      `json:"message"`
		Details interface{} `json:"details"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    "failed to parse error response",
		}
	}

	return &APIError{
		StatusCode: resp.StatusCode,
		Message:    errorResponse.Message,
		Details:    errorResponse.Details,
	}
}
