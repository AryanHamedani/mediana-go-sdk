# Mediana Go SDK

Go SDK for the Mediana.ir SMS API, providing services for sending SMS messages, including OTP, standard SMS, and pattern-based SMS.

## Installation

```bash
go get github.com/AryanHamedani/mediana-go-sdk
```

## Usage

### Initialization

```go
import "github.com/AryanHamedani/mediana-go-sdk/client"

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
    Recipients:    []string{"09123456789"},
    MessageText:   "Your message here",
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

### Check Delivery Status

```go
resp, err := c.GetDeliveryStatus(context.Background(), requestID)
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

## Testing the SDK

The SDK includes example code that demonstrates all available functionality. You can configure the example program using environment variables or command-line flags.

### Running the examples

1. Set your API key (required):

```bash
export MEDIANA_API_KEY="your_api_key_here"
```

2. Configure the examples using any of these methods:

   a. **Command-line flags**:

   ```bash
   cd examples
   go run main.go -phone 09XXXXXXXXX -sender 3000 -pattern your_pattern_code -otp-pattern your_otp_pattern
   ```

   b. **Environment variables**:

   ```bash
   export MEDIANA_TEST_PHONE="09XXXXXXXXX"
   export MEDIANA_SENDING_NUMBER="3000"
   export MEDIANA_PATTERN_CODE="your_pattern_code"
   export MEDIANA_OTP_PATTERN="your_otp_pattern"
   cd examples
   go run main.go
   ```

   c. **Default values**:
   If you don't specify values, the example will use default values where possible and display notices about required configurations.

### Available Configuration Options

| Parameter      | Flag           | Environment Variable     | Description                             |
| -------------- | -------------- | ------------------------ | --------------------------------------- |
| Phone Number   | `-phone`       | `MEDIANA_TEST_PHONE`     | Target phone number to send messages to |
| Sending Number | `-sender`      | `MEDIANA_SENDING_NUMBER` | Your Mediana sending number             |
| Pattern Code   | `-pattern`     | `MEDIANA_PATTERN_CODE`   | Pattern code for pattern-based SMS      |
| OTP Pattern    | `-otp-pattern` | `MEDIANA_OTP_PATTERN`    | Pattern code for OTP messages           |

## Docker Support

The SDK includes Docker support for easy building and testing.

### Using Docker

1. First, copy the example environment file and customize it:

```bash
cp .env.example .env
# Edit .env with your specific values
```

2. Build the Docker image:

```bash
docker build -t mediana-sdk .
```

3. Run the container with your environment variables:

```bash
docker run --env-file .env mediana-sdk
```

Alternatively, you can provide environment variables directly:

```bash
docker run -e MEDIANA_API_KEY=your_key -e MEDIANA_TEST_PHONE=09XXXXXXXXX mediana-sdk
```

## Examples

See the [examples](examples/) directory for complete usage examples.
