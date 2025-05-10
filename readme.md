# Mediana.ir SMS Go SDK

Go client library for interacting with the Mediana.ir SMS API.

## Installation

```bash
go get github.com/AryanHamedani/gosdk
```

## Usage

### Initialization

```go
import "github.com/AryanHamedani/gosdk/client"

apiKey := "your-api-key"
c := client.New(apiKey)

// With custom options
c := client.New(apiKey,
    client.WithBaseURL("https://custom.api.url"),
    client.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
)
```

### Sending SMS

```go
resp, err := c.SendSMS(context.Background(), models.SMSRequest{
    SendingNumber: "3000",
    Recipients:    []string{"09100000029"},
    MessageText:   "تست پیامک",
})
```

### Sending Pattern SMS

```go
resp, err := c.SendPatternSMS(context.Background(), models.PatternRequest{
    Recipients:  []string{"09123456789"},
    PatternCode: "welcome_pattern",
    Parameters: map[string]string{
        "name": "John",
        "code": "1234",
    },
})
```

### Sending OTP

```go
resp, err := c.SendOTP(context.Background(), models.OTPRequest{
    PatternCode: "otp_pattern",
    Recipient:   "09123456789",
    OTPCode:     "12345",
})
```

## Error Handling

All API errors are returned as `*errors.APIError` which includes:

- StatusCode: HTTP status code
- Message: Error message from API
- Details: Additional error details

```go
if err != nil {
    if apiErr, ok := err.(*errors.APIError); ok {
        log.Printf("API Error (%d): %s", apiErr.StatusCode, apiErr.Message)
    } else {
        log.Printf("Other error: %v", err)
    }
}
```

## Examples

See the [examples](examples/) directory for complete usage examples.
