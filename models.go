package notrix

import "time"

type CheckoutSessionLineItem struct {
	UUID        string  `json:"uuid"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

type CheckoutSession struct {
	UUID              string                    `json:"uuid"`
	LineItems         []CheckoutSessionLineItem `json:"line_items"`
	TotalAmount       string                    `json:"total_amount"`
	SuccessURL        string                    `json:"success_url"`
	CancelURL         string                    `json:"cancel_url"`
	ClientReferenceID string                    `json:"client_reference_id,omitempty"`
	WebhookURL        string                    `json:"webhook_url,omitempty"`
	CheckoutPageToken string                    `json:"checkout_page_token"`
	URL               string                    `json:"url"`
	Active            bool                      `json:"active"`
	Status            string                    `json:"status"`
	ExpiresAt         time.Time                 `json:"expires_at"`
	Metadata          map[string]interface{}    `json:"metadata"`
}
