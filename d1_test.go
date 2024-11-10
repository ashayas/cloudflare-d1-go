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

func TestListDB(t *testing.T) {
	client := cloudflare_d1_go.NewClient("", "")
	res, err := client.ListDB()
	if err != nil {
		t.Errorf("ListDB failed: %v", err)
	}
	t.Logf("ListDB response: %+v", res)

	if res == nil {
		t.Error("Expected non-nil response from ListDB")
	}
}
