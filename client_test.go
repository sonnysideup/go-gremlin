package gremlin

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

var AccessTokenBuilder = buildDefaultAccessToken()

const defaultURL = "https://api.gremlin.com/v1/"

func TestNewClientDefaultsConfirmsAuthenticationSucceeds(t *testing.T) {
	t.Log("Creating client using all defaults")
	// Given
	orgName := "Bob's Burgers, Inc."
	email := "real-email@google.com"
	password := "secure-password"

	accessTokenBuilt := AccessTokenBuilder.OrganizationName(orgName).Build()
	mockSucessAuth(defaultURL, []accessToken{accessTokenBuilt})
	defer httpmock.DeactivateAndReset()
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

	if got, want := client.token.Token, "fake-token"; got != want {
		t.Errorf("Expected token to be %#v, but got %#v", want, got)
	}
}

func TestNewClientDefaultsAuthenticationFailsWithWrongURL(t *testing.T) {
	t.Log("Creating client with wrong URL")
	// Given
	orgName := "Bob's Burgers, Inc."
	myURL := "http://not-real.io/"
	email := "real-email@google.com"
	password := "secure-password"

	defer func() {
		r := recover()
		if r != nil {
			panicMessage := fmt.Sprintf("%v", r)
			if got, want := panicMessage, "Auth request failed:"; !strings.Contains(got, want) {
				t.Errorf("Expected panic message to contain %s, but got %s", want, got)
			}
		} else {
			t.Error("Expected not real domain to throw panic.")
		}
	}()
	// When
	NewClient(orgName, email, password, WithURL(myURL))
}

func TestNewClientDefaultsAuthenticationFailsWithFailedRequest(t *testing.T) {
	t.Log("Creating client using all defaults and receives 401 auth response")
	// Given
	orgName := "Bob's Burgers, Inc."
	email := "real-email@google.com"
	password := "secure-password"

	mockFailAuth(defaultURL, 401)
	defer httpmock.DeactivateAndReset()

	defer func() {
		r := recover()
		if r != nil {
			panicMessage := fmt.Sprintf("%v", r)
			if got, want := panicMessage, "status: 401"; !strings.Contains(got, want) {
				t.Errorf("Expected panic message to contain %s, but got %s", want, got)
			}
		} else {
			t.Error("Expected non 200 response on auth response to cause panic.")
		}
	}()
	// When
	NewClient(orgName, email, password)
}

func TestNewClientDefaultsAuthenticationFailsWithBadResponseStructure(t *testing.T) {
	t.Log("Creating client using all defaults and fails to marshal auth response")
	// Given
	orgName := "Bob's Burgers, Inc."
	email := "real-email@google.com"
	password := "secure-password"

	mockBadResponseStructure(defaultURL)

	defer func() {
		r := recover()
		if r != nil {
			panicMessage := fmt.Sprintf("%v", r)
			if got, want := panicMessage, "Failed to marshall response:"; !strings.Contains(got, want) {
				t.Errorf("Expected panic message to contain %s, but got %s", want, got)
			}
		} else {
			t.Error("Expected response which doesn't match authToken struct to cause panic.")
		}
	}()
	// When
	NewClient(orgName, email, password)
}

func TestNewClientDefaultsAuthenticationFailsWhenTokenForOrganizationNotFound(t *testing.T) {
	t.Log("Creating client using all defaults and fails to find token for organization")
	// Given
	orgName := "Bob's Burgers, Inc."
	email := "real-email@google.com"
	password := "secure-password"

	accessTokenBuilt := AccessTokenBuilder.OrganizationName("Different Org").Build()
	mockSucessAuth(defaultURL, []accessToken{accessTokenBuilt})
	defer httpmock.DeactivateAndReset()

	defer func() {
		r := recover()
		if r != nil {
			panicMessage := fmt.Sprintf("%v", r)
			if got, want := panicMessage, "Unable to find token"; !strings.Contains(got, want) {
				t.Errorf("Expected panic message to contain %s, but got %s", want, got)
			}
		} else {
			t.Error("Expected inability to find organization token to cause panic.")
		}
	}()
	// When
	NewClient(orgName, email, password)
}

func TestNewClientWithDifferentURL(t *testing.T) {
	t.Log("Creating client with a non-default URL")
	// Given
	orgName := "New Co"
	myURL := "http://sweetpuppy.io/"
	email := "real-email@google.com"
	password := "secure-password"

	accessTokenBuilt := AccessTokenBuilder.OrganizationName(orgName).Build()
	mockSucessAuth(myURL, []accessToken{accessTokenBuilt})
	defer httpmock.DeactivateAndReset()
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

	accessTokenBuilt := AccessTokenBuilder.OrganizationName(orgName).Build()
	mockSucessAuth(defaultURL, []accessToken{accessTokenBuilt})
	defer httpmock.DeactivateAndReset()
	// When
	client := NewClient(orgName, email, password, WithNetClient(innerClient))

	// Then
	if got, want := client.client, innerClient; got != want {
		t.Errorf("Expected HTTP client to be %#v, but got %#v", want, got)
	}
}
