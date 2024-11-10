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
	return &Client{
		AccountID:  accountID,
		APIToken:   apiToken,
		HTTPClient: &http.Client{},
	}
}
