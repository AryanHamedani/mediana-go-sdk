package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// maxErrorBodyBytes caps how much of an error response body we retain.
const maxErrorBodyBytes = 2048

type APIError struct {
	StatusCode int
	Code       string
	Message    string
	Errors     []string
	Details    map[string]interface{}
}

func (e *APIError) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf("API error (%d): %s - %s", e.StatusCode, e.Message, strings.Join(e.Errors, ", "))
	}
	return fmt.Sprintf("API error (%d): %s", e.StatusCode, e.Message)
}

func ParseError(resp *http.Response) error {
	body, readErr := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodyBytes))
	if readErr != nil {
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("failed to read error response body: %v", readErr),
		}
	}

	var errorResponse struct {
		Meta struct {
			Code         string `json:"code"`
			ErrorMessage string `json:"errorMessage"`
			Errors       []struct {
				Key       string   `json:"key"`
				Errors    []string `json:"errors"`
				ErrorCode int      `json:"errorCode"`
			} `json:"errors"`
		} `json:"meta"`
		Data map[string]interface{} `json:"data"`
	}

	if err := json.Unmarshal(body, &errorResponse); err != nil {
		// Body is not the expected JSON envelope (e.g. an HTML/plain-text 429
		// page from a gateway). Surface the raw body so the failure stays
		// diagnosable instead of the opaque "failed to parse" message.
		raw := strings.TrimSpace(string(body))
		if raw == "" {
			raw = "(empty body)"
		}
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("unparseable error response: %s", raw),
		}
	}

	apiError := &APIError{
		StatusCode: resp.StatusCode,
		Code:       errorResponse.Meta.Code,
		Message:    errorResponse.Meta.ErrorMessage,
		Details:    errorResponse.Data,
		Errors:     []string{},
	}

	// Extract detailed error messages
	for _, errDetail := range errorResponse.Meta.Errors {
		for _, errMsg := range errDetail.Errors {
			apiError.Errors = append(apiError.Errors, fmt.Sprintf("%s: %s (code: %d)", errDetail.Key, errMsg, errDetail.ErrorCode))
		}
	}

	return apiError
}
