package gremlin

import (
	"net/http"
	"testing"
	"time"
)

func TestNewClientDefaults(t *testing.T) {
	t.Log("Creating client using all defaults")
	coName := "Bob's Burgers, Inc."
	client := NewClient(coName)

	if got, want := client.Company, coName; got != want {
		t.Errorf("Expected company name to equal %s, but got %s", want, got)
	}

	if got, want := client.BaseURL.String(), "https://api.gremlin.com/v1/"; got != want {
		t.Errorf("Expected BaseURL to equal %s, but got %s", want, got)
	}

	if got, want := client.client.Timeout, time.Second*10; got != want {
		t.Error("Inner net client incorrectly initialized")
		t.Errorf("Expected default timeout to equal %s, but got %s", want, got)
	}
}

func TestNewClientWithDifferentURL(t *testing.T) {
	t.Log("Creating client with a non-default URL")

	myURL := "http://sweetpuppy.io"
	client := NewClient("NewCo", WithURL(myURL))

	if got, want := client.BaseURL.String(), myURL; got != want {
		t.Errorf("Expected configured URL to be %s, but got %s", want, got)
	}
}

func TestNewClientWithInvalidURL(t *testing.T) {
	t.Log("Creating client with a bad URL")
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected bad URL to cause panic")
		}
	}()

	NewClient("NewCo", WithURL("://derpa.derp"))
}

func TestNewClientWithCustomInnerHTTPClient(t *testing.T) {
	t.Log("Creating client with non-default HTTPClient")

	innerClient := &http.Client{Timeout: 60}
	client := NewClient("NewCo", WithNetClient(innerClient))

	if got, want := client.client, innerClient; got != want {
		t.Errorf("Expected HTTP client to be %#v, but got %#v", want, got)
	}
}

// TODO add tests for authenticate method
