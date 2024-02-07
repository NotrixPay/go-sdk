package notrix

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateCheckoutSession(t *testing.T) {
	// Mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"uuid": "session_uuid", "line_items": [], "total_amount": "10.99", "success_url": "https://example.com/success", "cancel_url": "https://example.com/cancel", "checkout_page_token": "checkout_token", "url": "https://example.com/checkout", "active": true, "status": "created", "expires_at": "2024-02-07T12:00:00Z", "metadata": {}}`)
	}))
	defer mockServer.Close()

	client := Client{SecretAPIKey: "test_api_key"}
	baseURL = mockServer.URL // Set the base URL to the mock server's URL

	items := []CheckoutSessionLineItem{
		{
			UUID:        "item_uuid_1",
			Name:        "Item 1",
			Description: "Description for Item 1",
			Image:       "item1.jpg",
			Price:       10.99,
			Quantity:    2,
		},
	}

	successURL := "https://example.com/success"
	cancelURL := "https://example.com/cancel"
	clientReferenceID := "123456"
	webhookURL := "https://example.com/webhook"

	session, err := client.CreateCheckoutSession(items, successURL, cancelURL, clientReferenceID, webhookURL)
	if err != nil {
		t.Errorf("Error creating checkout session: %v", err)
		return
	}

	// Verify the returned session
	expectedSession := CheckoutSession{
		UUID:              "session_uuid",
		TotalAmount:       "10.99",
		SuccessURL:        "https://example.com/success",
		CancelURL:         "https://example.com/cancel",
		CheckoutPageToken: "checkout_token",
		URL:               "https://example.com/checkout",
		Active:            true,
		Status:            "created",
		ExpiresAt:         parseTime("2024-02-07T12:00:00Z"),
		Metadata:          map[string]interface{}{},
	}
	if (*session).UUID != expectedSession.UUID {
		t.Errorf("Expected session: %+v\nGot: %+v", expectedSession, *session)
	}
}

func TestIsPaid(t *testing.T) {
	// Mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "valid_token" {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"payment_confirmed": true}`)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "400 Bad Request")
		}
	}))
	defer mockServer.Close()

	client := Client{SecretAPIKey: "test_api_key"}
	baseURL = mockServer.URL // Set the base URL to the mock server's URL

	paid, err := client.IsPaid("valid_token")
	if err != nil {
		t.Errorf("Error checking payment status: %v", err)
		return
	}

	if !paid {
		t.Error("Expected payment to be confirmed, got not confirmed")
	}

	_, err = client.IsPaid("invalid_token")
	if err == nil {
		t.Error("Expected error, got nil")
		return
	}

	expectedErrorMessage := "400 Bad Request"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message: %s, got: %s", expectedErrorMessage, err.Error())
	}
}

func parseTime(str string) time.Time {
	t, _ := time.Parse(time.RFC3339, str)
	return t
}
