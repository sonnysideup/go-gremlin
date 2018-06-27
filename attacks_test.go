package gremlin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

// helper methods

func setup() (mux *http.ServeMux, client *Client, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)

	client = NewClient("Test Org", "user@domain.com", "secret", WithURL(server.URL))
	client.token = &accessToken{Header: "Bearer fake-token"}

	return mux, client, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %s, want %s", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}

func buildAttack() AttackCommand {
	return AttackCommand{
		Command: Command{
			Type: "cpu",
			Args: []string{"-c", "1", "--length", "5"},
		},
		Target: Target{
			Type:  "Exact",
			Exact: []string{"some-client"},
		},
	}
}

// tests
// - test failure with bad input

func TestCreateAttackSuccess(t *testing.T) {
	mux, client, teardown := setup()
	defer teardown()

	attack := buildAttack()
	someGUID := "123e4567-e89b-12d3-a456-426655440000"

	mux.HandleFunc("/attacks/new", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)

		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/json")
		testHeader(t, r, "Authorization", "Bearer fake-token")

		v := new(AttackCommand)
		json.NewDecoder(r.Body).Decode(v)

		if !reflect.DeepEqual(v, &attack) {
			t.Errorf("Request body = %+v, want %+v", v, attack)
		}

		fmt.Fprint(w, someGUID)
	})

	guid, err := client.CreateAttack(attack)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if guid.String() != someGUID {
		t.Errorf("Expected guid to be %s, but got %s", someGUID, guid)
	}
}

func TestCreateAttackWithServiceUnavailable(t *testing.T) {
	mux, client, teardown := setup()
	defer teardown()

	attack := buildAttack()

	mux.HandleFunc("/attacks/new", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	})

	_, err := client.CreateAttack(attack)

	if err == nil {
		t.Errorf("Expected service unavailable to result in error")
	} else {
		if got, want := err.Error(), "status: 503"; !strings.Contains(got, want) {
			t.Errorf("Expected error to match %q, but got %q", want, got)
		}
	}
}

func TestCreateAttackWithBadResponse(t *testing.T) {
	mux, client, teardown := setup()
	defer teardown()

	attack := buildAttack()

	mux.HandleFunc("/attacks/new", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "abc-its-easy-as-1-2-3")
	})

	_, err := client.CreateAttack(attack)

	if err == nil {
		t.Errorf("Expected invalid response to result in error")
	} else {
		if got, want := err.Error(), "Invalid UUID from server:"; !strings.Contains(got, want) {
			t.Errorf("Expected error to match %q, but got %q", want, got)
		}
	}
}
