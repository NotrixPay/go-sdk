package notrix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var (
	baseURL = "https://api.notrix.io"
)

type Client struct {
	SecretAPIKey string
}

func (c *Client) authHeaders() map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Token %s", c.SecretAPIKey),
	}
}

func (c *Client) makeRequest(method, path string, body interface{}, params map[string]string) (*http.Response, error) {
	baseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	baseURL.Path = path

	query := baseURL.Query()
	for key, value := range params {
		query.Add(key, value)
	}
	baseURL.RawQuery = query.Encode()

	url := baseURL.String()

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	for headerName, headerValue := range c.authHeaders() {
		request.Header.Add(headerName, headerValue)
	}
	client := &http.Client{}
	return client.Do(request)
}

func (c *Client) CreateCheckoutSession(items []CheckoutSessionLineItem, successURL, cancelURL, clientReferenceID, webhookURL string) (*CheckoutSession, error) {
	body := map[string]interface{}{
		"success_url": successURL,
		"cancel_url":  cancelURL,
		"line_items":  items,
	}
	if clientReferenceID != "" {
		body["client_reference_id"] = clientReferenceID
	}
	if webhookURL != "" {
		body["webhook_url"] = webhookURL
	}

	response, err := c.makeRequest("POST", "console/checkout-sessions/", body, nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var session CheckoutSession
	if err := json.NewDecoder(response.Body).Decode(&session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (c *Client) IsPaid(checkoutPageToken string) (bool, error) {
	params := map[string]string{
		"token": checkoutPageToken,
	}
	response, err := c.makeRequest("GET", "console/check-payment-status/", nil, params)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errString, err := io.ReadAll(response.Body)
		if err == nil {
			err = fmt.Errorf(string(errString))
		}
		return false, err
	}

	var paymentStatus struct {
		PaymentConfirmed bool `json:"payment_confirmed"`
	}
	if err := json.NewDecoder(response.Body).Decode(&paymentStatus); err != nil {
		return false, err
	}
	return paymentStatus.PaymentConfirmed, nil
}
