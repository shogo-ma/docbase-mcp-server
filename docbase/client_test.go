package docbase

import (
	"testing"
)

func TestNewDocBaseClient(t *testing.T) {
	domain := "example"
	apiToken := "test-token"
	client := NewDocBaseClient(domain, apiToken)

	expectedBaseURL := "https://api.docbase.io/teams/example"
	if client.BaseURL != expectedBaseURL {
		t.Errorf("Expected base URL to be %q, but got %q", expectedBaseURL, client.BaseURL)
	}

	if client.Domain != domain {
		t.Errorf("Expected domain to be %q, but got %q", domain, client.Domain)
	}

	if client.APIToken != apiToken {
		t.Errorf("Expected API token to be %q, but got %q", apiToken, client.APIToken)
	}

	if client.Client == nil {
		t.Error("Expected http.Client to be initialized, but it's nil")
	}
}
