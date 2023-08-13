package approvedeny_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ownage-FDT/approvedeny-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	t.Run("create client with a valid api key", func(t *testing.T) {
		apiKey := "test_api_key"
		client, err := approvedeny.NewClient("test_api_key")

		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, "https://api.approvedeny.com", client.BaseURL)
		assert.Equal(t, apiKey, client.APIKey)
		assert.NotNil(t, client.HttpClient)
	})

	t.Run("create client with invalid api key", func(t *testing.T) {
		client, err := approvedeny.NewClient("")

		assert.Nil(t, client)
		assert.NotNil(t, err)
	})
}

func TestClient_GetCheckRequest(t *testing.T) {
	t.Run("get check request with valid check request id", func(t *testing.T) {
		checkRequestID := "test_check_request_id"
		expectedResponse := &approvedeny.SuccessResponse{
			Status: "success",
			Data: map[string]interface{}{
				"id": checkRequestID,
			},
		}

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/v1/requests/"+checkRequestID, r.URL.Path)

			responseBody, err := json.Marshal(expectedResponse)
			assert.NoError(t, err)

			w.WriteHeader(http.StatusOK)
			w.Write(responseBody)
		})

		server := httptest.NewServer(handler)
		defer server.Close()

		client := &approvedeny.Client{
			BaseURL:    server.URL,
			APIKey:     "test_api_key",
			HttpClient: server.Client(),
		}

		response, err := client.GetCheckRequest(checkRequestID)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
	})

	t.Run("get check request with invalid check request id", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message": "Check request not found"}`))
		})

		server := httptest.NewServer(handler)
		defer server.Close()

		client := &approvedeny.Client{
			BaseURL:    server.URL,
			APIKey:     "test_api_key",
			HttpClient: server.Client(),
		}

		response, err := client.GetCheckRequest("invalid_check_request_id")
		assert.Nil(t, response)
		assert.NotNil(t, err)
	})
}

func TestClient_CreateCheckRequest(t *testing.T) {
	t.Run("create check request with valid check id and payload", func(t *testing.T) {
		checkRequestID := "test_check_request_id"
		expectedResponse := &approvedeny.SuccessResponse{
			Status: "success",
			Data: map[string]interface{}{
				"id": checkRequestID,
			},
		}

		payload := approvedeny.CreateCheckRequestPayload{
			Description: "A test check request",
			Metadata: map[string]interface{}{
				"foo": "bar",
			},
		}

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/v1/checks/"+checkRequestID, r.URL.Path)

			responseBody, err := json.Marshal(expectedResponse)
			assert.NoError(t, err)

			w.WriteHeader(http.StatusOK)
			w.Write(responseBody)
		})

		server := httptest.NewServer(handler)
		defer server.Close()

		client := &approvedeny.Client{
			BaseURL:    server.URL,
			APIKey:     "test_api_key",
			HttpClient: server.Client(),
		}

		response, err := client.CreateCheckRequest(checkRequestID, payload)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
	})

	t.Run("create check request with invalid check id", func(t *testing.T) {
		payload := approvedeny.CreateCheckRequestPayload{
			Description: "A test check request",
			Metadata: map[string]interface{}{
				"foo": "bar",
			},
		}

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message": "Check request not found"}`))
		})

		server := httptest.NewServer(handler)
		defer server.Close()

		client := &approvedeny.Client{
			BaseURL:    server.URL,
			APIKey:     "test_api_key",
			HttpClient: server.Client(),
		}

		response, err := client.CreateCheckRequest("invalid_check_id", payload)
		assert.Nil(t, response)
		assert.NotNil(t, err)
	})
}

func TestClient_GetCheckRequestResponse(t *testing.T) {
	t.Run("get check request response with valid check request id", func(t *testing.T) {
		checkRequestID := "test_check_request_id"
		expectedResponse := &approvedeny.SuccessResponse{
			Status: "success",
			Data: map[string]interface{}{
				"id": checkRequestID,
			},
		}

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/v1/requests/"+checkRequestID+"/response", r.URL.Path)
			responseBody, err := json.Marshal(expectedResponse)
			assert.NoError(t, err)

			w.WriteHeader(http.StatusOK)
			w.Write(responseBody)
		})

		server := httptest.NewServer(handler)
		defer server.Close()

		client := &approvedeny.Client{
			BaseURL:    server.URL,
			APIKey:     "test_api_key",
			HttpClient: server.Client(),
		}

		response, err := client.GetCheckRequestResponse(checkRequestID)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
	})

	t.Run("get check request response with invalid check request id", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message": "Check request not found"}`))
		})

		server := httptest.NewServer(handler)
		defer server.Close()

		client := &approvedeny.Client{
			BaseURL:    server.URL,
			APIKey:     "test_api_key",
			HttpClient: server.Client(),
		}

		response, err := client.GetCheckRequestResponse("invalid_check_request_id")
		assert.Nil(t, response)
		assert.EqualError(t, err, "Check request not found")
	})
}

func TestClient_IsValidWebhookSignature(t *testing.T) {
	t.Run("valid webhook signature", func(t *testing.T) {
		encryptionKey := "test_encryption_key"
		payload := approvedeny.WebhookPayload{
			Event: "response.created",
			Data: map[string]interface{}{
				"ID": "test_id",
			},
		}

		h := hmac.New(sha256.New, []byte(encryptionKey))
		encodedPayload, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("Failed to marshal payload: %v", err)
		}
		h.Write(encodedPayload)
		signature := fmt.Sprintf("%x", h.Sum(nil))

		client := &approvedeny.Client{}

		isValid := client.IsValidWebhookSignature(encryptionKey, signature, payload)
		assert.True(t, isValid)
	})

	t.Run("invalid webhook signature", func(t *testing.T) {
		encryptionKey := "test_encryption_key"
		signature := "invalid_signature"
		payload := approvedeny.WebhookPayload{
			Event: "response.created",
			Data: map[string]interface{}{
				"id": "test_id",
			},
		}

		client := &approvedeny.Client{}

		isValid := client.IsValidWebhookSignature(encryptionKey, signature, payload)
		assert.False(t, isValid)
	})
}
