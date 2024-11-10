package cloudflared1

import (
	"fmt"

	"github.com/ashayas/cloudflare-d1-go/utils"
)

type Client struct {
	AccountID string
	APIToken  string
}

func NewClient(accountID, apiToken string) *Client {
	if accountID == "" || apiToken == "" {
		return nil
	}
	return &Client{
		AccountID: accountID,
		APIToken:  apiToken,
	}
}

func (c *Client) ListDB() (*utils.APIResponse, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/d1/database", c.AccountID)
	return utils.DoRequest("GET", url, "", c.APIToken)
}
