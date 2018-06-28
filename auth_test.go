package gremlin

import (
	"net/http"
	"strings"
	"testing"
	"time"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

var AccessTokenBuilder = buildDefaultAccessToken()

const defaultURL = "https://api.gremlin.com/v1/"

// NewClient() tests

func TestNewClientDefaults(t *testing.T) {
	t.Log("Creating client using all defaults")

	// Given
	orgName := "Bob's Burgers, Inc."
	email := "real-email@google.com"
	password := "secure-password"

	// When
	client := NewClient(orgName, email, password)

	// Then
	if got, want := client.Company, orgName; got != want {
		t.Errorf("Expected company name to equal %s, but got %s", want, got)
	}

	if got, want := client.BaseURL.String(), defaultURL; got != want {
		t.Errorf("Expected BaseURL to equal %s, but got %s", want, got)
	}

	if got, want := client.client.Timeout, time.Second*10; got != want {
		t.Error("Inner net client incorrectly initialized")
		t.Errorf("Expected default timeout to equal %s, but got %s", want, got)
	}
}

func TestNewClientWithDifferentURL(t *testing.T) {
	t.Log("Creating client with a non-default URL")

	// Given
	orgName := "New Co"
	myURL := "http://sweetpuppy.io/"
	email := "real-email@google.com"
	password := "secure-password"

	// When
	client := NewClient(orgName, email, password, WithURL(myURL))

	// Then
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
	email := "real-email@google.com"
	password := "secure-password"

	// When
	NewClient("NewCo", email, password, WithURL("://derpa.derp"))
}

func TestNewClientWithCustomInnerHTTPClient(t *testing.T) {
	t.Log("Creating client with non-default HTTPClient")

	// Given
	orgName := "NewCo"
	innerClient := &http.Client{Timeout: 60}
	email := "real-email@google.com"
	password := "secure-password"

	// When
	client := NewClient(orgName, email, password, WithNetClient(innerClient))

	// Then
	if got, want := client.client, innerClient; got != want {
		t.Errorf("Expected HTTP client to be %#v, but got %#v", want, got)
	}
}

// authenticate() tests

func TestAuthenticationSuccess(t *testing.T) {
	// Given
	orgName := "Bob's Burgers, Inc."
	email := "real-email@google.com"
	password := "secure-password"
	client := NewClient(orgName, email, password)

	accessTokenBuilt := AccessTokenBuilder.OrganizationName(orgName).Build()
	mockSucessAuth(defaultURL, []accessToken{accessTokenBuilt})
	defer httpmock.DeactivateAndReset()

	// When
	client.Authenticate()

	// Then
	if got, want := client.Token.Token, "fake-token"; got != want {
		t.Errorf("Expected token to be %#v, but got %#v", want, got)
	}
}

func TestAuthenticationFailsWithWrongURL(t *testing.T) {
	t.Log("Creating client with wrong URL")

	// Given
	orgName := "Bob's Burgers, Inc."
	myURL := "http://not-real.io/"
	email := "real-email@google.com"
	password := "secure-password"
	client := NewClient(orgName, email, password, WithURL(myURL))

	// When
	_, err := client.Authenticate()

	// Then
	if err == nil {
		t.Error("Expected not real domain to result in error.")
	} else {
		if got, want := err.Error(), "Request failed"; !strings.Contains(got, want) {
			t.Errorf("Expected error message to contain %q, but got %q", want, got)
		}
	}
}

func TestAuthenticationFailsWithFailedRequest(t *testing.T) {
	t.Log("Creating client using all defaults and receives 401 auth response")

	// Given
	orgName := "Bob's Burgers, Inc."
	email := "real-email@google.com"
	password := "secure-password"
	client := NewClient(orgName, email, password)

	mockFailAuth(defaultURL, 401)
	defer httpmock.DeactivateAndReset()

	// When
	_, err := client.Authenticate()

	// Then
	if err == nil {
		t.Error("Expected non-200 response on auth to result in error")
	} else {
		if got, want := err.Error(), "status: 401"; !strings.Contains(got, want) {
			t.Errorf("Expected error message to contain %q, but got %q", want, got)
		}
	}
}

func TestAuthenticationFailsWithBadResponseStructure(t *testing.T) {
	t.Log("Creating client using all defaults and fails to marshal auth response")

	// Given
	orgName := "Bob's Burgers, Inc."
	email := "real-email@google.com"
	password := "secure-password"
	client := NewClient(orgName, email, password)

	mockBadResponseStructure(defaultURL)

	// When
	_, err := client.Authenticate()

	// Then
	if err == nil {
		t.Error("Expected malformed response to result in error")
	} else {
		if got, want := err.Error(), "Failed to marshall response:"; !strings.Contains(got, want) {
			t.Errorf("Expected error message to contain %q, but got %q", want, got)
		}
	}
}

func TestAuthenticationFailsWhenTokenForOrganizationNotFound(t *testing.T) {
	t.Log("Creating client using all defaults and fails to find token for organization")

	// Given
	orgName := "Bob's Burgers, Inc."
	email := "real-email@google.com"
	password := "secure-password"
	client := NewClient(orgName, email, password)

	accessTokenBuilt := AccessTokenBuilder.OrganizationName("Different Org").Build()
	mockSucessAuth(defaultURL, []accessToken{accessTokenBuilt})
	defer httpmock.DeactivateAndReset()

	// When
	_, err := client.Authenticate()

	// Then
	if err == nil {
		t.Error("Expected missing org token to result in error")
	} else {
		if got, want := err.Error(), "Unable to find token"; !strings.Contains(got, want) {
			t.Errorf("Expected error message to contain %q, but got %q", want, got)
		}
	}
}
