package approvedeny

import (
	"net/http"
)

type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type CreateCheckRequestPayload struct {
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type CheckRequestResponse struct {
	ID             string                 `json:"id"`
	Status         string                 `json:"status"`
	Metadata       map[string]interface{} `json:"metadata"`
	CheckRequestID string                 `json:"checkRequestId"`
	CreatedAt      string                 `json:"createdAt"`
	UpdatedAt      string                 `json:"updatedAt"`
}

type WebhookPayload struct {
	Event string                 `json:"event"`
	Data  map[string]interface{} `json:"data"`
}

type CheckRequest struct {
	ID          string                 `json:"id"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
	CheckID     string                 `json:"checkId"`
	CreatedAt   string                 `json:"createdAt"`
	UpdatedAt   string                 `json:"updatedAt"`
	Response    *CheckRequestResponse  `json:"response,omitempty"`
}

type Client struct {
	BaseURL    string
	APIKey     string
	HttpClient *http.Client
}
