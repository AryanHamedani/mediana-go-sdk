package main

import (
	"context"
	"log"
	"os"

	"github.com/mediana-ir/gosdk/client"
	"github.com/mediana-ir/gosdk/models"
)

func main() {
	apiKey := os.Getenv("MEDIANA_API_KEY")
	if apiKey == "" {
		log.Fatal("MEDIANA_API_KEY environment variable not set")
	}

	c := client.New(apiKey)

	// Example 1: Send regular SMS
	sendRegularSMS(c)

	// Example 2: Send pattern SMS
	sendPatternSMS(c)

	// Example 3: Send OTP
	sendOTP(c)
}

func sendRegularSMS(c *client.Client) {
	ctx := context.Background()

	resp, err := c.SendSMS(ctx, models.SMSRequest{
		SendingNumber: "3000", // Your sending number
		Recipients:    []string{"09100000029"},
		MessageText:   "hi",
	})

	if err != nil {
		log.Printf("Failed to send SMS: %v", err)
		return
	}

	log.Printf("SMS sent successfully: %+v", resp)
}

func sendPatternSMS(c *client.Client) {
	ctx := context.Background()

	resp, err := c.SendPatternSMS(ctx, models.PatternRequest{
		Recipients:  []string{"09123456789"},
		PatternCode: "welcome_pattern",
		Parameters: map[string]string{
			"name": "John",
			"code": "1234",
		},
	})

	if err != nil {
		log.Printf("Failed to send pattern SMS: %v", err)
		return
	}

	log.Printf("Pattern SMS sent successfully: %+v", resp)
}

func sendOTP(c *client.Client) {
	ctx := context.Background()

	resp, err := c.SendOTP(ctx, models.OTPRequest{
		PatternCode: "otp_pattern",
		Recipient:   "09123456789",
		OTPCode:     "12345",
	})

	if err != nil {
		log.Printf("Failed to send OTP: %v", err)
		return
	}

	log.Printf("OTP sent successfully: %+v", resp)
}
