package errors

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

// errReader always fails, simulating a response body that cannot be read.
type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, io.ErrClosedPipe }

func newResp(status int, body io.Reader) *http.Response {
	return &http.Response{StatusCode: status, Body: io.NopCloser(body)}
}

func TestParseError_JSONEnvelope(t *testing.T) {
	body := `{"meta":{"code":"1041","errorMessage":"validation failed",` +
		`"errors":[{"key":"recipients","errors":["invalid receiver"],"errorCode":1041}]}}`

	err := ParseError(newResp(http.StatusBadRequest, strings.NewReader(body)))

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusBadRequest)
	}
	if apiErr.Code != "1041" {
		t.Errorf("Code = %q, want %q", apiErr.Code, "1041")
	}
	if apiErr.Message != "validation failed" {
		t.Errorf("Message = %q, want %q", apiErr.Message, "validation failed")
	}
	want := "recipients: invalid receiver (code: 1041)"
	if len(apiErr.Errors) != 1 || apiErr.Errors[0] != want {
		t.Errorf("Errors = %v, want [%q]", apiErr.Errors, want)
	}
}

func TestParseError_NonJSONBody(t *testing.T) {
	body := "<html><body>429 Too Many Requests</body></html>"

	err := ParseError(newResp(http.StatusTooManyRequests, strings.NewReader(body)))

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusTooManyRequests {
		t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusTooManyRequests)
	}
	if !strings.HasPrefix(apiErr.Message, "unparseable error response:") {
		t.Errorf("Message = %q, want %q prefix", apiErr.Message, "unparseable error response:")
	}
	if !strings.Contains(apiErr.Message, body) {
		t.Errorf("Message = %q, want it to contain the raw body", apiErr.Message)
	}
}

func TestParseError_EmptyBody(t *testing.T) {
	err := ParseError(newResp(http.StatusBadGateway, strings.NewReader("   ")))

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	want := "unparseable error response: (empty body)"
	if apiErr.Message != want {
		t.Errorf("Message = %q, want %q", apiErr.Message, want)
	}
}

func TestParseError_ReadError(t *testing.T) {
	err := ParseError(newResp(http.StatusInternalServerError, errReader{}))

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if !strings.Contains(apiErr.Message, "failed to read error response body") {
		t.Errorf("Message = %q, want it to report a read failure", apiErr.Message)
	}
}
