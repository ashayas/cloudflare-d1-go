package cloudflared1

import (
	"fmt"
	"strings"

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

func (c *Client) CreateDB(name string) (*utils.APIResponse, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/d1/database", c.AccountID)
	body := fmt.Sprintf(`{"name":"%s"}`, name)
	return utils.DoRequest("POST", url, body, c.APIToken)
}

func (c *Client) DeleteDB(databaseID string) (*utils.APIResponse, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/d1/database/%s", c.AccountID, databaseID)
	return utils.DoRequest("DELETE", url, "", c.APIToken)
}

// Runs SQL query on the D1 database with parameters
func (c *Client) QueryDB(databaseID string, query string, params []string) (*utils.APIResponse, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/d1/database/%s/raw", c.AccountID, databaseID)

	// Create the request body with params and sql
	body := fmt.Sprintf(`{
		"sql": "%s",
		"params": %s
	}`, query, formatParams(params))

	return utils.DoRequest("POST", url, body, c.APIToken)
}

// Helper function to format parameters as a JSON array
func formatParams(params []string) string {
	if len(params) == 0 {
		return "[]"
	}

	quoted := make([]string, len(params))
	for i, p := range params {
		quoted[i] = fmt.Sprintf(`"%s"`, p)
	}
	return fmt.Sprintf(`[%s]`, strings.Join(quoted, ","))
}

func (c *Client) CreateTable(databaseID, createQuery string) (*utils.APIResponse, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/d1/database/%s/raw", c.AccountID, databaseID)
	body := fmt.Sprintf(`{
		"sql": "%s",
		"params": []
	}`, createQuery)
	return utils.DoRequest("POST", url, body, c.APIToken)
}

func (c *Client) RemoveTable(databaseID, tableName string) (*utils.APIResponse, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/d1/database/%s/raw", c.AccountID, databaseID)
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s;", tableName)
	body := fmt.Sprintf(`{
		"sql": "%s",
		"params": []
	}`, query)
	return utils.DoRequest("POST", url, body, c.APIToken)
}
