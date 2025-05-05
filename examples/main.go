package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mediana-ir/gosdk/client"
	"github.com/mediana-ir/gosdk/errors"
	"github.com/mediana-ir/gosdk/models"
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

	// Track request IDs for status checking
	var requestIDs []string

	// Example 1: Send regular SMS
	fmt.Println("\n📱 Example 1: Sending regular SMS...")
	smsID, err := sendRegularSMS(c, phoneNumber, sendingNumber)
	if err == nil && smsID != "" {
		requestIDs = append(requestIDs, smsID)
	}

	// Example 2: Send pattern SMS
	fmt.Println("\n📊 Example 2: Sending pattern SMS...")
	patternID, err := sendPatternSMS(c, phoneNumber, patternCode)
	if err == nil && patternID != "" {
		requestIDs = append(requestIDs, patternID)
	}

	// Example 3: Send OTP
	fmt.Println("\n🔑 Example 3: Sending OTP...")
	otpID, err := sendOTP(c, phoneNumber, otpPattern)
	if err == nil && otpID != "" {
		requestIDs = append(requestIDs, otpID)
	}
	
	// Wait a bit for messages to be processed
	if len(requestIDs) > 0 {
		fmt.Println("\n⏳ Waiting 10 seconds for message processing...")
		time.Sleep(10 * time.Second)
		
		// Example 4: Check delivery status
		fmt.Println("\n📊 Example 4: Checking delivery status...")
		// You'll need to use an actual request ID here for real testing
		checkDeliveryStatus(c, 123456) // Default test ID
		
		// If we have actual request IDs from this run, check those as well
		for i, id := range requestIDs {
			if id != "" {
				fmt.Printf("\nChecking delivery status for request ID from example %d...\n", i+1)
				// You may need to convert the string ID to an integer based on your API response
				// For now, this is a placeholder
				fmt.Printf("⚠️ Request ID: %s - Unable to check (format conversion needed)\n", id)
			}
		}
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
	
	fmt.Printf("✅ SMS sent successfully: Status=%s\n", resp.Status)
	
	// Return the first message ID if available
	if len(resp.MessageIDs) > 0 {
		fmt.Printf("📝 Message ID: %s\n", resp.MessageIDs[0])
		return resp.MessageIDs[0], nil
	}
	return "", nil
}

func sendPatternSMS(c *client.Client, phoneNumber, patternCode string) (string, error) {
	ctx := context.Background()

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
	
	fmt.Printf("✅ Pattern SMS sent successfully: Status=%s\n", resp.Status)
	
	// Return the first message ID if available
	if len(resp.MessageIDs) > 0 {
		fmt.Printf("📝 Message ID: %s\n", resp.MessageIDs[0])
		return resp.MessageIDs[0], nil
	}
	return "", nil
}

func sendOTP(c *client.Client, phoneNumber, otpPattern string) (string, error) {
	ctx := context.Background()

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
	
	fmt.Printf("✅ OTP sent successfully: Status=%s, OTP Code=%s\n", resp.Status, otpCode)
	
	// Return the message ID if available
	if resp.MessageID != "" {
		fmt.Printf("📝 Message ID: %s\n", resp.MessageID)
		return resp.MessageID, nil
	}
	return "", nil
}

func checkDeliveryStatus(c *client.Client, requestID int) {
	ctx := context.Background()

	// Print request details
	fmt.Printf("📤 Checking delivery status for request ID: %d\n", requestID)

	resp, err := c.GetDeliveryStatus(ctx, requestID)
	if err != nil {
		fmt.Printf("❌ Failed to check delivery status: %v\n", err)
		printDetailedError(err)
		return
	}

	// Print full response
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
		
		// If there's a raw response body, print it
		if apiErr.RawBody != "" {
			fmt.Printf("  - Response Body: %s\n", apiErr.RawBody)
		}
		
		// If there are specific API error details, print them
		if apiErr.Message != "" {
			fmt.Printf("  - Message: %s\n", apiErr.Message)
		}
		if apiErr.Code != "" {
			fmt.Printf("  - Code: %s\n", apiErr.Code)
		}
		return
	}
	
	// For other error types, just print the error
	fmt.Printf("🔍 Error Details: %v\n", err)
}
