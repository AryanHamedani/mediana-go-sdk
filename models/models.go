package models

// Meta represents common metadata in responses
type Meta struct {
	Code         string   `json:"code"`
	ErrorMessage string   `json:"errorMessage,omitempty"`
	Errors       []string `json:"errors,omitempty"`
}

// SmsItemInfo represents information about a single SMS item in responses
type SmsItemInfo struct {
	SmsItemId string `json:"smsItemId"`
	Recipient string `json:"recipient"`
	Status    string `json:"status,omitempty"`
}

// SMSRequest represents a request to send a regular SMS
type SMSRequest struct {
	Type          string   `json:"type,omitempty"` // Informational, PromotionalToCustomers, PromotionalAll
	SendingNumber string   `json:"sendingNumber,omitempty"`
	Recipients    []string `json:"recipients"`
	MessageText   string   `json:"messageText"`
}

// SMSResponse represents the response from sending an SMS
type SMSResponse struct {
	Meta Meta `json:"meta"`
	Data struct {
		Succeed     bool         `json:"succeed"`
		RequestCode string       `json:"requestCode"`
		Message     string       `json:"message"`
		Status      string       `json:"status"`
		SmsItems    []SmsItemInfo `json:"smsItems"`
	} `json:"data"`
}

// PatternRequest represents a request to send a pattern SMS
type PatternRequest struct {
	Recipients  []string          `json:"recipients"`
	PatternCode string            `json:"patternCode"`
	Parameters  map[string]string `json:"parameters"`
}

// PatternResponse represents the response from sending a pattern SMS
type PatternResponse struct {
	Meta Meta `json:"meta"`
	Data struct {
		Succeed     bool         `json:"succeed"`
		RequestCode string       `json:"requestCode"`
		Message     string       `json:"message"`
		Status      string       `json:"status"`
		SmsItems    []SmsItemInfo `json:"smsItems"`
	} `json:"data"`
}

// OTPRequest represents a request to send an OTP SMS
type OTPRequest struct {
	PatternCode string `json:"patternCode"`
	Recipient   string `json:"recipient"`
	OTPCode     string `json:"otpCode"`
}

// OTPResponse represents the response from sending an OTP SMS
type OTPResponse struct {
	Meta Meta `json:"meta"`
	Data struct {
		Succeed     bool         `json:"succeed"`
		RequestCode string       `json:"requestCode"`
		Message     string       `json:"message"`
		Status      string       `json:"status"`
		SmsItems    []SmsItemInfo `json:"smsItems"`
	} `json:"data"`
}

// DeliveryStatusResponse represents the response from checking message delivery status
type DeliveryStatusResponse struct {
	Meta Meta `json:"meta"`
	Data struct {
		Status   string        `json:"status"`
		SmsItems []SmsItemInfo `json:"smsItems"`
	} `json:"data"`
}

// BalanceResponse represents the response for account balance inquiry
type BalanceResponse struct {
	Meta Meta `json:"meta"`
	Data struct {
		Balance int `json:"Balance"`
	} `json:"data"`
}

// LineInfo represents a single sending line information
type LineInfo struct {
	Number         string `json:"Number"`
	Description    string `json:"Description"`
	IsDedicated    bool   `json:"IsDedicated"`
	IsAdvertisement bool  `json:"IsAdvertisement"`
	IsService      bool   `json:"IsService"`
	UsableUntil    string `json:"UsableUntil"`
}

// LinesResponse represents the response for account lines query
type LinesResponse struct {
	Meta Meta      `json:"meta"`
	Data LineInfo  `json:"data"`
}

// PatternDetailResponse represents the response for a pattern detail query
type PatternDetailResponse struct {
	Meta Meta `json:"meta"`
	Data struct {
		MessagePatternId int    `json:"MessagePatternId"`
		Title            string `json:"Title"`
		Type             string `json:"type"`
		IsUsable         bool   `json:"IsUsable"`
		Code             string `json:"Code"`
		Description      string `json:"Description"`
		ThePattern       struct {
			Pattern               string `json:"Pattern"`
			Status               string `json:"Status"`
			SendingNumber        string `json:"SendingNumber"`
			IsLockedBySendingNumber bool  `json:"IsLockedBySendingNumber"`
			ApprovalDescription  string `json:"ApprovalDescription"`
			Fields               []struct {
				FieldTitle     string `json:"FieldTitle"`
				FieldKey       string `json:"FieldKey"`
				MaxCharacters  int    `json:"MaxCharacters"`
				FieldType      string `json:"FieldType"`
			} `json:"GetMessagePatternsByIdResponseField"`
		} `json:"ThePattern"`
		SettingInfo struct {
			Website            string `json:"Website"`
			AverageSendingCount int    `json:"AverageSendingCount"`
		} `json:"SettingInfo"`
		Patterns    []struct {
			Pattern string `json:"Pattern"`
			Status  string `json:"Status"`
		} `json:"Patterns"`
		CreateDate string `json:"CreateDate"`
	} `json:"data"`
}
