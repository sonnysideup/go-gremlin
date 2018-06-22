package gremlin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	// defaultBaseURL points to the default location of the API
	defaultBaseURL, _ = url.Parse("https://api.gremlin.com/v1/")

	// defaultNetClient sets a default timeout of 10sec
	defaultNetClient = &http.Client{Timeout: time.Second * 10}
)

// Client manages communication with the Gremlin API
type Client struct {
	client *http.Client

	Company string
	BaseURL *url.URL
}

// ConfigOption represents the type interface that can be used to add new
// functional options to the Client constructor.
type ConfigOption func(*Client) error

// WithURL can be used to point the Client at a different API server (e.g.
// enterprise instance running on-premise). This offering currently dne.
func WithURL(urlStr string) ConfigOption {
	return func(c *Client) error {
		customURL, err := url.Parse(urlStr)
		if err != nil {
			return fmt.Errorf("Failed to parse API URL: %v", err)
		}

		c.BaseURL = customURL
		return nil
	}
}

// WithNetClient should be used to override the default http client settings.
func WithNetClient(netClient *http.Client) ConfigOption {
	return func(c *Client) error {
		c.client = netClient
		return nil
	}
}

// Generate a new Gremlin Client, populating the required fields.
func NewClient(company string, options ...ConfigOption) *Client {
	// default client settings
	client := &Client{
		client:  defaultNetClient,
		Company: company,
		BaseURL: defaultBaseURL,
	}

	// apply any functional options
	for _, option := range options {
		if err := option(client); err != nil {
			panic(err)
		}
	}

	return client
}

// Authenticate provides your user credentials to Gremlin and requests an access
// token.
//
// All API requests require an access token so you'll need to provide one to all
// other method invocations.
func (c *Client) Authenticate(email string, password string) (*accessToken, error) {
	rurl := c.resourceURL("users/auth")

	// create request body and object
	form := url.Values{}
	form.Set("email", email)
	form.Set("password", password)
	form.Set("companyName", c.Company)

	req, err := http.NewRequest("POST", rurl.String(), strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatalln("Failed to create new request obj:", err)
	}

	// set required header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// dispatch request and check response status
	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatalln("Auth request failed:", err)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Auth request failed: status: %d body: %s\n", resp.StatusCode, string(body))
	}
	defer resp.Body.Close()

	// marshall JSON response into object
	var tokens []accessToken
	if err = json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		log.Fatalln("Failed to marshall response:", err)
	}

	// search for required company token
	for _, t := range tokens {
		if t.OrganizationName == c.Company {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("Unable to find token for '%s'\nTokens returned: %+v\n",
		c.Company,
		tokens,
	)
}

// resourceURL safely joins a string path (e.g. "my/resource") to an existing URL.
func (c *Client) resourceURL(path string) *url.URL {
	rel := &url.URL{Path: path}
	return c.BaseURL.ResolveReference(rel)
}
