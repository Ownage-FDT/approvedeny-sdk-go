package approvedeny

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func NewClient(apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("The Approvedeny SDK requires an API key to be provided at initialization")
	}

	return &Client{
		BaseURL:    "https://api.approvedeny.com",
		APIKey:     apiKey,
		HttpClient: &http.Client{},
	}, nil
}

func (c *Client) makeRequest(method, url string, payload interface{}) (*SuccessResponse, error) {
	var requestBody bytes.Buffer = bytes.Buffer{}

	if payload != nil {
		encodedPayload, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("Error marshalling payload: %s", err.Error())
		}
		requestBody = *bytes.NewBuffer(encodedPayload)
	}

	req, err := http.NewRequest(method, url, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %s", err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "approvedeny-go/1.0.0")

	response, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %s", err.Error())
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %s", err.Error())
	}

	// if status code is not 200, return error
	if response.StatusCode != http.StatusOK {
		// unmarsh the error response
		var errorResponse ErrorResponse
		if err := json.Unmarshal(responseBody, &errorResponse); err != nil {
			return nil, fmt.Errorf("Error unmarshalling response: %s", err.Error())
		}

		return nil, fmt.Errorf(errorResponse.Message)
	}

	// umarsh the success response
	var successResponse SuccessResponse
	if err := json.Unmarshal(responseBody, &successResponse); err != nil {
		return nil, fmt.Errorf("Error unmarshalling response: %s", err.Error())
	}

	return &successResponse, nil
}

func (c *Client) GetCheckRequest(checkRequestID string) (*SuccessResponse, error) {
	url := fmt.Sprintf("%s/v1/requests/%s", c.BaseURL, checkRequestID)

	response, err := c.makeRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) CreateCheckRequest(checkID string, payload CreateCheckRequestPayload) (*SuccessResponse, error) {
	url := fmt.Sprintf("%s/v1/checks/%s", c.BaseURL, checkID)

	response, err := c.makeRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) GetCheckRequestResponse(checkRequestID string) (*SuccessResponse, error) {
	url := fmt.Sprintf("%s/v1/requests/%s/response", c.BaseURL, checkRequestID)

	response, err := c.makeRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) IsValidWebhookSignature(encryptionKey, signature string, payload WebhookPayload) bool {
	h := hmac.New(sha256.New, []byte(encryptionKey))

	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		return false
	}

	h.Write(encodedPayload)

	expectedSignature := fmt.Sprintf("%x", h.Sum(nil))
	return expectedSignature == signature
}
