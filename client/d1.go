package cloudflared1

import (
	"fmt"
	"strings"

	"github.com/ashayas/cloudflare-d1-go/utils"
)

type Client struct {
	AccountID  string
	APIToken   string
	DatabaseID string
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

func (c *Client) CreateTableWithID(databaseID, createQuery string) (*utils.APIResponse, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/d1/database/%s/raw", c.AccountID, databaseID)
	body := fmt.Sprintf(`{
		"sql": "%s",
		"params": []
	}`, createQuery)
	return utils.DoRequest("POST", url, body, c.APIToken)
}

func (c *Client) RemoveTableWithID(databaseID, tableName string) (*utils.APIResponse, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/d1/database/%s/raw", c.AccountID, databaseID)
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s;", tableName)
	body := fmt.Sprintf(`{
		"sql": "%s",
		"params": []
	}`, query)
	return utils.DoRequest("POST", url, body, c.APIToken)
}

// ConnectDB finds and connects to a database by name, storing its ID for future operations
func (c *Client) ConnectDB(name string) error {
	resp, err := c.ListDB()
	if err != nil {
		return fmt.Errorf("failed to list databases: %w", err)
	}

	// Parse response to find database with matching name
	databases := resp.Result.([]interface{})
	for _, db := range databases {
		dbMap := db.(map[string]interface{})
		if dbMap["name"].(string) == name {
			c.DatabaseID = dbMap["uuid"].(string)
			return nil
		}
	}

	return fmt.Errorf("database with name %s not found", name)
}

// Query runs SQL query on the connected database
func (c *Client) Query(query string, params []string) (*utils.APIResponse, error) {
	if c.DatabaseID == "" {
		return nil, fmt.Errorf("no database connected, call ConnectDB first")
	}
	return c.QueryDB(c.DatabaseID, query, params)
}

// CreateTable creates a table in the connected database
func (c *Client) CreateTable(createQuery string) (*utils.APIResponse, error) {
	if c.DatabaseID == "" {
		return nil, fmt.Errorf("no database connected, call ConnectDB first")
	}
	return c.CreateTableWithID(c.DatabaseID, createQuery)
}

// RemoveTable removes a table from the connected database
func (c *Client) RemoveTable(tableName string) (*utils.APIResponse, error) {
	if c.DatabaseID == "" {
		return nil, fmt.Errorf("no database connected, call ConnectDB first")
	}
	return c.RemoveTableWithID(c.DatabaseID, tableName)
}
