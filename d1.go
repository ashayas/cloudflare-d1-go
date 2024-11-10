package cloudflared1

import (
	"net/http"
)

type Client struct {
	AccountID  string
	APIToken   string
	HTTPClient *http.Client
}

func NewClient(accountID, apiToken string) *Client {
	// Return nil if either accountID or apiToken is empty
	if accountID == "" || apiToken == "" {
		return nil
	}

	return &Client{
		AccountID:  accountID,
		APIToken:   apiToken,
		HTTPClient: &http.Client{},
	}
}
