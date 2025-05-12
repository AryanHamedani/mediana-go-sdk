package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AryanHamedani/mediana-go-sdk/client"
	"github.com/AryanHamedani/mediana-go-sdk/errors"
	"github.com/AryanHamedani/mediana-go-sdk/models"
)

func main() {
	// Parse command line flags
	phoneNumberFlag := flag.String("phone", "", "Phone number to test (format: 09XXXXXXXXX)")
	sendingNumberFlag := flag.String("sender", "", "Sending number to use")
	patternCodeFlag := flag.String("pattern", "", "Pattern code for pattern SMS")
	otpPatternFlag := flag.String("otp-pattern", "", "Pattern code for OTP")
	flag.Parse()

	// Get API key from environment variable
	apiKey := os.Getenv("MEDIANA_API_KEY")
	if apiKey == "" {
		log.Fatal("MEDIANA_API_KEY environment variable not set")
	}

	// Create a new client
	c := client.New(apiKey)

	// Get phone number from flag or environment variable, with default fallback
	phoneNumber := *phoneNumberFlag
	if phoneNumber == "" {
		// Try to get from environment variable
		phoneNumber = os.Getenv("MEDIANA_TEST_PHONE")
		if phoneNumber == "" {
			// Default phone number as fallback
			phoneNumber = "09199686112"
			fmt.Printf("Using default phone number: %s\n", phoneNumber)
			fmt.Println("To use a different number, set MEDIANA_TEST_PHONE environment variable or use -phone flag")
		} else {
			fmt.Printf("Using phone number from environment: %s\n", phoneNumber)
		}
	} else {
		fmt.Printf("Using phone number from command line: %s\n", phoneNumber)
	}

	// Get sending number
	sendingNumber := *sendingNumberFlag
	if sendingNumber == "" {
		sendingNumber = os.Getenv("MEDIANA_SENDING_NUMBER")
		if sendingNumber == "" {
			sendingNumber = "3000" // Default
			fmt.Printf("Using default sending number: %s\n", sendingNumber)
		} else {
			fmt.Printf("Using sending number from environment: %s\n", sendingNumber)
		}
	} else {
		fmt.Printf("Using sending number from command line: %s\n", sendingNumber)
	}

	// Get pattern code
	patternCode := *patternCodeFlag
	if patternCode == "" {
		patternCode = os.Getenv("MEDIANA_PATTERN_CODE")
		if patternCode == "" {
			patternCode = "your_pattern_code" // You should replace this with a valid pattern code
			fmt.Printf("Using placeholder pattern code: %s\n", patternCode)
			fmt.Println("⚠️ Replace with a valid pattern code for proper testing")
		} else {
			fmt.Printf("Using pattern code from environment: %s\n", patternCode)
		}
	} else {
		fmt.Printf("Using pattern code from command line: %s\n", patternCode)
	}

	// Get OTP pattern code
	otpPattern := *otpPatternFlag
	if otpPattern == "" {
		otpPattern = os.Getenv("MEDIANA_OTP_PATTERN")
		if otpPattern == "" {
			otpPattern = "your_otp_pattern" // You should replace this with a valid OTP pattern code
			fmt.Printf("Using placeholder OTP pattern code: %s\n", otpPattern)
			fmt.Println("⚠️ Replace with a valid OTP pattern code for proper testing")
		} else {
			fmt.Printf("Using OTP pattern code from environment: %s\n", otpPattern)
		}
	} else {
		fmt.Printf("Using OTP pattern code from command line: %s\n", otpPattern)
	}

	// Track request codes for status checking
	var requestCodes []string

	// Example 1: Send regular SMS
	fmt.Println("\n📱 Example 1: Sending regular SMS...")
	requestCode, err := sendRegularSMS(c, phoneNumber, sendingNumber)
	if err == nil && requestCode != "" {
		requestCodes = append(requestCodes, requestCode)
	}

	// Example 2: Send pattern SMS
	fmt.Println("\n📊 Example 2: Sending pattern SMS...")
	patternReqCode, err := sendPatternSMS(c, phoneNumber, patternCode)
	if err == nil && patternReqCode != "" {
		requestCodes = append(requestCodes, patternReqCode)
	}

	// Example 3: Send OTP
	fmt.Println("\n🔑 Example 3: Sending OTP...")
	otpReqCode, err := sendOTP(c, phoneNumber, otpPattern)
	if err == nil && otpReqCode != "" {
		requestCodes = append(requestCodes, otpReqCode)
	}

	// Test: Get Account Balance
	if apiKey != "" {
		fmt.Println("\n💰 Testing GetAccountBalance...")
		testGetAccountBalance(c)
	}

	// Test: Get Sending Lines
	if apiKey != "" {
		fmt.Println("\n📞 Testing GetSendingLines...")
		testGetSendingLines(c)
	}

	// Test: Get Pattern Detail
	if patternCode != "" && patternCode != "your_pattern_code" {
		fmt.Println("\n🔍 Testing GetPatternDetail...")
		testGetPatternDetail(c, patternCode)
	} else {
		fmt.Println("\n⚠️ Skipping GetPatternDetail test: pattern code not set.")
	}

	// Wait a bit for messages to be processed
	if len(requestCodes) > 0 {
		fmt.Println("\n⏳ Waiting 10 seconds for message processing...")
		time.Sleep(10 * time.Second)

		// Example 4: Check delivery status for each request code from previous tests
		fmt.Println("\n📊 Example 4: Checking delivery status...")
		
		for i, code := range requestCodes {
			if code != "" {
				fmt.Printf("\nChecking delivery status for request code from example %d...\n", i+1)
				checkDeliveryStatus(c, code)
			}
		}
	} else {
		fmt.Println("\n⚠️ No request codes available for status check. All send operations failed or returned empty codes.")
	}
}

func sendRegularSMS(c *client.Client, phoneNumber, sendingNumber string) (string, error) {
	ctx := context.Background()

	// Print request details
	req := models.SMSRequest{
		SendingNumber: sendingNumber,
		Recipients:    []string{phoneNumber},
		MessageText:   "Hello! This is a test message from Mediana Go SDK.",
	}

	reqJSON, _ := json.MarshalIndent(req, "", "  ")
	fmt.Printf("📤 Request:\n%s\n", string(reqJSON))

	resp, err := c.SendSMS(ctx, req)

	if err != nil {
		fmt.Printf("❌ Failed to send SMS: %v\n", err)
		printDetailedError(err)
		return "", err
	}

	// Print full response
	respJSON, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Printf("📥 Response:\n%s\n", string(respJSON))

	fmt.Printf("✅ SMS sent successfully. Succeed: %v, Status: %s\n", resp.Data.Succeed, resp.Data.Status)

	// Return the request code for status checking
	if resp.Data.RequestCode != "" {
		fmt.Printf("📝 Request Code: %s\n", resp.Data.RequestCode)
		
		// Also print SMS Items IDs if available
		if len(resp.Data.SmsItems) > 0 {
			fmt.Println("📱 SMS Items:")
			for _, item := range resp.Data.SmsItems {
				fmt.Printf("  - Item ID: %s, Recipient: %s\n", item.SmsItemId, item.Recipient)
			}
		}
		
		return resp.Data.RequestCode, nil
	}
	return "", nil
}

func sendPatternSMS(c *client.Client, phoneNumber, patternCode string) (string, error) {
	ctx := context.Background()
	
	// Skip if pattern code is clearly a placeholder
	if patternCode == "your_pattern_code" {
		fmt.Println("⚠️ Skipping pattern SMS test: using placeholder pattern code")
		return "", nil
	}

	// Print request details
	req := models.PatternRequest{
		Recipients:  []string{phoneNumber},
		PatternCode: patternCode,
		Parameters: map[string]string{
			"name": "Test User",
			"code": "98765",
			// Add other parameters required by your pattern
		},
	}

	reqJSON, _ := json.MarshalIndent(req, "", "  ")
	fmt.Printf("📤 Request:\n%s\n", string(reqJSON))

	resp, err := c.SendPatternSMS(ctx, req)

	if err != nil {
		fmt.Printf("❌ Failed to send pattern SMS: %v\n", err)
		printDetailedError(err)
		return "", err
	}

	// Print full response
	respJSON, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Printf("📥 Response:\n%s\n", string(respJSON))

	fmt.Printf("✅ Pattern SMS sent successfully. Succeed: %v, Status: %s\n", resp.Data.Succeed, resp.Data.Status)

	// Return the request code for status checking
	if resp.Data.RequestCode != "" {
		fmt.Printf("📝 Request Code: %s\n", resp.Data.RequestCode)
		
		// Also print SMS Items IDs if available
		if len(resp.Data.SmsItems) > 0 {
			fmt.Println("📱 SMS Items:")
			for _, item := range resp.Data.SmsItems {
				fmt.Printf("  - Item ID: %s, Recipient: %s\n", item.SmsItemId, item.Recipient)
			}
		}
		
		return resp.Data.RequestCode, nil
	}
	return "", nil
}

func sendOTP(c *client.Client, phoneNumber, otpPattern string) (string, error) {
	ctx := context.Background()
	
	// Skip if OTP pattern is clearly a placeholder
	if otpPattern == "your_otp_pattern" {
		fmt.Println("⚠️ Skipping OTP test: using placeholder OTP pattern code")
		return "", nil
	}

	// Generate a random OTP code (in real scenarios, you would generate this securely)
	otpCode := "123456"

	// Print request details
	req := models.OTPRequest{
		PatternCode: otpPattern,
		Recipient:   phoneNumber,
		OTPCode:     otpCode,
	}

	reqJSON, _ := json.MarshalIndent(req, "", "  ")
	fmt.Printf("📤 Request:\n%s\n", string(reqJSON))

	resp, err := c.SendOTP(ctx, req)

	if err != nil {
		fmt.Printf("❌ Failed to send OTP: %v\n", err)
		printDetailedError(err)
		return "", err
	}

	// Print full response
	respJSON, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Printf("📥 Response:\n%s\n", string(respJSON))

	fmt.Printf("✅ OTP sent successfully. Succeed: %v, Status: %s, OTP Code: %s\n", 
		resp.Data.Succeed, resp.Data.Status, otpCode)

	// Return the request code for status checking
	if resp.Data.RequestCode != "" {
		fmt.Printf("📝 Request Code: %s\n", resp.Data.RequestCode)
		
		// Also print SMS Items IDs if available
		if len(resp.Data.SmsItems) > 0 {
			fmt.Println("📱 SMS Items:")
			for _, item := range resp.Data.SmsItems {
				fmt.Printf("  - Item ID: %s, Recipient: %s\n", item.SmsItemId, item.Recipient)
			}
		}
		
		return resp.Data.RequestCode, nil
	}
	return "", nil
}

func testGetAccountBalance(c *client.Client) {
	ctx := context.Background()
	resp, err := c.GetAccountBalance(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to get account balance: %v\n", err)
		printDetailedError(err)
		return
	}
	respJSON, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Printf("📥 Response:\n%s\n", string(respJSON))
	fmt.Printf("✅ Account Balance: %d\n", resp.Data.Balance)
}

func testGetSendingLines(c *client.Client) {
	ctx := context.Background()
	resp, err := c.GetSendingLines(ctx)
	if err != nil {
		fmt.Printf("❌ Failed to get sending lines: %v\n", err)
		printDetailedError(err)
		return
	}
	respJSON, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Printf("📥 Response:\n%s\n", string(respJSON))
	fmt.Println("✅ Sending Line:")
	fmt.Printf("  - Number: %s, Description: %s, Dedicated: %v, Advertisement: %v, Service: %v, UsableUntil: %s\n",
		resp.Data.Number, resp.Data.Description, resp.Data.IsDedicated, resp.Data.IsAdvertisement, 
		resp.Data.IsService, resp.Data.UsableUntil)
}

func testGetPatternDetail(c *client.Client, patternCode string) {
	ctx := context.Background()
	fmt.Printf("🔍 Requesting pattern detail for code: %s\n", patternCode)
	resp, err := c.GetPatternDetail(ctx, patternCode)
	if err != nil {
		fmt.Printf("❌ Failed to get pattern detail: %v\n", err)
		printDetailedError(err)
		return
	}
	respJSON, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Printf("📥 Response:\n%s\n", string(respJSON))
	fmt.Printf("✅ Pattern Title: %s, Usable: %v\n", resp.Data.Title, resp.Data.IsUsable)
}

func checkDeliveryStatus(c *client.Client, requestIDStr string) {
	ctx := context.Background()
	var requestID int
	_, err := fmt.Sscanf(requestIDStr, "%d", &requestID)
	if err != nil {
		fmt.Printf("⚠️ Unable to parse request ID '%s' as int: %v\n", requestIDStr, err)
		return
	}
	fmt.Printf("📤 Checking delivery status for request ID: %d\n", requestID)
	resp, err := c.GetDeliveryStatus(ctx, requestID)
	if err != nil {
		fmt.Printf("❌ Failed to check delivery status: %v\n", err)
		printDetailedError(err)
		return
	}
	respJSON, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Printf("📥 Response:\n%s\n", string(respJSON))
	fmt.Printf("✅ Delivery status: %s\n", resp.Data.Status)
	fmt.Println("📱 SMS Items:")
	for _, item := range resp.Data.SmsItems {
		fmt.Printf("  - ID: %s, Recipient: %s, Status: %s\n",
			item.SmsItemId, item.Recipient, item.Status)
	}
}

func printDetailedError(err error) {
	// Check if the error is an API error which contains more details
	if apiErr, ok := err.(*errors.APIError); ok {
		fmt.Printf("🔍 API Error Details:\n")
		fmt.Printf("  - Status Code: %d\n", apiErr.StatusCode)
		fmt.Printf("  - Error: %s\n", apiErr.Error())

		// If there are specific API error details, print them
		if apiErr.Message != "" {
			fmt.Printf("  - Message: %s\n", apiErr.Message)
		}
		if apiErr.Code != "" {
			fmt.Printf("  - Code: %s\n", apiErr.Code)
		}
		if len(apiErr.Errors) > 0 {
			fmt.Printf("  - Detailed Errors:\n")
			for _, e := range apiErr.Errors {
				fmt.Printf("    - %s\n", e)
			}
		}
		return
	}

	// For other error types, just print the error
	fmt.Printf("🔍 Error Details: %v\n", err)
}
