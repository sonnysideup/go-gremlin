package gremlin

import (
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	coName := "Bob's Burgers, Inc."
	client := NewClient(coName, nil)

	if got, want := client.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("Expected BaseURL to equal %s, but got %s", want, got)
	}

	if got, want := client.CompanyName, coName; got != want {
		t.Errorf("Expected company name to equal %s, but got %s", want, got)
	}

	if client.httpClient != http.DefaultClient {
		t.Error("HTTP Client incorrectly initialized")
	}
}
