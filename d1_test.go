package cloudflared1_test

import (
	"testing"

	cloudflare_d1_go "github.com/ashayas/cloudflare-d1-go/client"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name      string
		accountID string
		apiToken  string
		wantErr   bool
	}{
		{
			name:      "valid credentials",
			accountID: "1234567890",
			apiToken:  "1234567890",
			wantErr:   false,
		},
		{
			name:      "empty account ID",
			accountID: "",
			apiToken:  "1234567890",
			wantErr:   true,
		},
		{
			name:      "empty API token",
			accountID: "1234567890",
			apiToken:  "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := cloudflare_d1_go.NewClient(tt.accountID, tt.apiToken)

			if tt.wantErr {
				if client != nil {
					t.Errorf("NewClient() = %v, want nil for invalid inputs", client)
				}
				return
			}

			if client == nil {
				t.Fatal("NewClient() returned nil for valid inputs")
			}

			if client.AccountID != tt.accountID {
				t.Errorf("NewClient().AccountID = %v, want %v", client.AccountID, tt.accountID)
			}

			if client.APIToken != tt.apiToken {
				t.Errorf("NewClient().APIToken = %v, want %v", client.APIToken, tt.apiToken)
			}
		})
	}
}

// TestListDB lists the databases
func TestListDB(t *testing.T) {
	client := cloudflare_d1_go.NewClient("account_id", "api_token")
	res, err := client.ListDB()
	if err != nil {
		t.Errorf("ListDB failed: %v", err)
	}
	t.Logf("ListDB response: %+v", res)

	if res == nil {
		t.Error("Expected non-nil response from ListDB")
	}
}

// TestCreateAndDeleteDB creates a database, then deletes it
func TestCreateAndDeleteDB(t *testing.T) {
	client := cloudflare_d1_go.NewClient("account_id", "api_token")
	res, err := client.CreateDB("test-db-2")
	if err != nil {
		t.Errorf("CreateDB failed: %v", err)
	}
	t.Logf("CreateDB response: %+v", res)

	if res == nil {
		t.Error("Expected non-nil response from CreateDB")
	}

	// Only do this if the database was created successfully
	if res != nil && res.Success {
		res, err = client.DeleteDB(res.Result.(map[string]interface{})["uuid"].(string))
		if err != nil {
			t.Errorf("DeleteDB failed: %v", err)
		}
		t.Logf("DeleteDB response: %+v", res)

		if res == nil {
			t.Error("Expected non-nil response from DeleteDB")
		}
	}

}

// TestCreateAndRemoveTable creates a table, then removes it
func TestCreateAndRemoveTable(t *testing.T) {
	client := cloudflare_d1_go.NewClient("account_id", "api_token")

	// Create a test database
	res, err := client.CreateDB("test_db_3")
	if err != nil {
		t.Errorf("CreateDB failed: %v", err)
		return
	}

	if !res.Success {
		t.Errorf("CreateDB was not successful: %v", res.Errors)
		return
	}

	dbID, ok := res.Result.(map[string]interface{})["uuid"].(string)
	if !ok {
		t.Error("Failed to get database UUID from response")
		return
	}

	// Create a test table
	createQuery := "CREATE TABLE IF NOT EXISTS test_table (id INTEGER PRIMARY KEY, name TEXT);"
	res, err = client.CreateTable(dbID, createQuery)
	if err != nil {
		t.Errorf("CreateTable failed: %v", err)
		return
	}

	if !res.Success {
		t.Errorf("CreateTable was not successful: %v", res.Errors)
		return
	}

	t.Logf("CreateTable response: %+v", res)

	// Only attempt to remove if table was created successfully
	res, err = client.RemoveTable(dbID, "test_table")
	if err != nil {
		t.Errorf("RemoveTable failed: %v", err)
	}
	t.Logf("RemoveTable response: %+v", res)

	if !res.Success {
		t.Errorf("RemoveTable was not successful: %v", res.Errors)
	}
}

// TestQueryDB creates a table, inserts a row, then selects it and deletes the table and database
func TestQueryDB(t *testing.T) {
	client := cloudflare_d1_go.NewClient("account_id", "api_token")

	// Create a test database
	res, err := client.CreateDB("test_db_6")
	if err != nil {
		t.Errorf("CreateDB failed: %v", err)
		return
	}

	if !res.Success {
		t.Errorf("CreateDB was not successful: %v", res.Errors)
		return
	}

	dbID, ok := res.Result.(map[string]interface{})["uuid"].(string)
	if !ok {
		t.Error("Failed to get database UUID from response")
		return
	}

	// Create a test table
	createQuery := "CREATE TABLE IF NOT EXISTS test_table (id INTEGER PRIMARY KEY, name TEXT);"
	res, err = client.CreateTable(dbID, createQuery)
	if err != nil || !res.Success {
		t.Errorf("CreateTable failed: %v, errors: %v", err, res.Errors)
		return
	}

	t.Logf("CreateTable response: %+v", res)

	// Insert test data
	insertQuery := "INSERT INTO test_table (name) VALUES (?);"
	params := []string{"test_name"}
	res, err = client.QueryDB(dbID, insertQuery, params)
	if err != nil || !res.Success {
		t.Errorf("Insert query failed: %v, errors: %v", err, res.Errors)
		return
	}

	// Select the data
	selectQuery := "SELECT * FROM test_table WHERE name = ?;"
	res, err = client.QueryDB(dbID, selectQuery, params)
	if err != nil {
		t.Errorf("Select query failed: %v", err)
		return
	}

	if !res.Success {
		t.Errorf("Select query was not successful: %v", res.Errors)
		return
	}

	// Parse the response correctly
	results, ok := res.Result.([]interface{})
	if !ok {
		t.Error("Failed to parse Result as array")
		return
	}

	if len(results) == 0 {
		t.Error("No results returned from query")
		return
	}

	// First element contains the query results
	queryResult, ok := results[0].(map[string]interface{})
	if !ok {
		t.Error("Failed to parse query result")
		return
	}

	// Access the actual rows from the results
	resultsMap, ok := queryResult["results"].(map[string]interface{})
	if !ok {
		t.Error("Failed to parse results map")
		return
	}

	rows, ok := resultsMap["rows"].([]interface{})
	if !ok {
		t.Error("Failed to parse rows")
		return
	}

	// Delete the table
	res, err = client.RemoveTable(dbID, "test_table")
	if err != nil || !res.Success {
		t.Errorf("RemoveTable failed: %v, errors: %v", err, res.Errors)
		return
	}

	// Delete the database
	res, err = client.DeleteDB(dbID)
	if err != nil || !res.Success {
		t.Errorf("DeleteDB failed: %v, errors: %v", err, res.Errors)
		return
	}

	t.Logf("Query returned %d rows", len(rows))
	t.Logf("Full response: %+v", res)
}
