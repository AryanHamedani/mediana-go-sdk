package models

// SMSRequest represents a request to send a regular SMS
type SMSRequest struct {
	SendingNumber string   `json:"sendingNumber"`
	Recipients    []string `json:"recipients"`
	MessageText   string   `json:"messageText"`
}

// SMSResponse represents the response from sending an SMS
type SMSResponse struct {
	Status     string   `json:"status"`
	MessageIDs []string `json:"messageIds"`
}

// PatternRequest represents a request to send a pattern SMS
type PatternRequest struct {
	Recipients  []string          `json:"recipients"`
	PatternCode string            `json:"patternCode"`
	Parameters  map[string]string `json:"parameters"`
}

// PatternResponse represents the response from sending a pattern SMS
type PatternResponse struct {
	Status     string   `json:"status"`
	MessageIDs []string `json:"messageIds"`
}

// OTPRequest represents a request to send an OTP SMS
type OTPRequest struct {
	PatternCode string `json:"patternCode"`
	Recipient   string `json:"recipient"`
	OTPCode     string `json:"otpCode"`
}

// OTPResponse represents the response from sending an OTP SMS
type OTPResponse struct {
	Status    string `json:"status"`
	MessageID string `json:"messageId"`
}

// DeliveryStatusResponse represents the response from checking message delivery status
type DeliveryStatusResponse struct {
	Meta struct {
		Code         string   `json:"code"`
		ErrorMessage string   `json:"errorMessage,omitempty"`
		Errors       []string `json:"errors,omitempty"`
	} `json:"meta"`
	Data struct {
		Status   string        `json:"status"`
		SmsItems []SmsItemInfo `json:"smsItems"`
	} `json:"data"`
}

// SmsItemInfo represents information about a single SMS item in a delivery status response
type SmsItemInfo struct {
	SmsItemId string `json:"smsItemId"`
	Recipient string `json:"recipient"`
	Status    string `json:"status"`
}
