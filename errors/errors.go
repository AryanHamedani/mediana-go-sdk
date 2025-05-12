package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

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
	var errorResponse struct {
		Meta struct {
			Code         string `json:"code"`
			ErrorMessage string `json:"errorMessage"`
			Errors       []struct {
				Key      string   `json:"key"`
				Errors   []string `json:"errors"`
				ErrorCode int      `json:"errorCode"`
			} `json:"errors"`
		} `json:"meta"`
		Data map[string]interface{} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    "failed to parse error response",
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
