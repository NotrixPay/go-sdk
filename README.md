# Notrix API Go Client

This Go package provides a client for interacting with the Notrix API.

## Installation

To install the Notrix API Go client, use the following `go get` command:

```bash
go get github.com/NotrixPay/go-sdk
```

## Usage

Import the notrix package in your Go code:

```go
import (
	"fmt"
	"github.com/NotrixPay/go-sdk"
)
```

Create a Client instance with your secret API key:

```go
client := notrix.NewClient("your_secret_api_key")
```

## Creating a Checkout Session

To create a checkout session, use the CreateCheckoutSession method:

```go
items := []notrix.CheckoutSessionLineItem{
    {
        UUID:        "(empty string)",
        Name:        "Item 1",
        Description: "Description for Item 1",
        Image:       "item1.jpg",
        Price:       10.99,
        Quantity:    2,
    },
    // Add more items if needed
}

successURL := "https://example.com/success"
cancelURL := "https://example.com/cancel"
clientReferenceID := "123456"               // Optional can be left empty
webhookURL := "https://example.com/webhook" // Optional can be left empty

session, err := client.CreateCheckoutSession(items, successURL, cancelURL, clientReferenceID, webhookURL)
if err != nil {
    fmt.Println("Error creating checkout session:", err)
    return
}

fmt.Println("Checkout session created successfully with URL:", session.URL)
```

## Checking Payment Status

To check the payment status of a checkout session, use the IsPaid method:

```go
checkoutPageToken := "your_checkout_page_token"
paid, err := client.IsPaid(checkoutPageToken)
if err != nil {
    fmt.Println("Error checking payment status:", err)
    return
}

if paid {
    fmt.Println("Payment confirmed")
} else {
    fmt.Println("Payment not confirmed yet")
}
```