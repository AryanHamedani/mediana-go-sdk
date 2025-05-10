package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AryanHamedani/gosdk/models"
)

func (c *Client) SendSMS(ctx context.Context, req models.SMSRequest) (*models.SMSResponse, error) {
	resp, err := c.doRequest(ctx, "POST", "send/sms", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response models.SMSResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

func (c *Client) SendPatternSMS(ctx context.Context, req models.PatternRequest) (*models.PatternResponse, error) {
	resp, err := c.doRequest(ctx, "POST", "send/pattern", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response models.PatternResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

func (c *Client) SendOTP(ctx context.Context, req models.OTPRequest) (*models.OTPResponse, error) {
	resp, err := c.doRequest(ctx, "POST", "send/otp", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response models.OTPResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// GetDeliveryStatus retrieves the delivery status of a message by request ID
func (c *Client) GetDeliveryStatus(ctx context.Context, requestID int) (*models.DeliveryStatusResponse, error) {
	endpoint := fmt.Sprintf("send-requests/status/%d", requestID)
	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response models.DeliveryStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
