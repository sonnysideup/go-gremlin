package gremlin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	client   *http.Client
	Company  string
	BaseURL  *url.URL
	Email    string
	password string
	token    *accessToken
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
func NewClient(company string, email string, password string, options ...ConfigOption) *Client {
	// default client settings
	client := &Client{
		client:   defaultNetClient,
		Company:  company,
		BaseURL:  defaultBaseURL,
		Email:    email,
		password: password,
		token:    &accessToken{},
	}

	// apply any functional options
	for _, option := range options {
		if err := option(client); err != nil {
			panic(err)
		}
	}

	return client
}

// authenticate provides your user credentials to Gremlin and requests an access
// token. If authentication is successful, this token will be stored in the client.
//
// All other API requests require an access token so a token will be required
// prior to invoking other methods.
func (c *Client) authenticate() (*accessToken, error) {
	rurl := c.resourceURL("users/auth")

	// create request body and object
	form := url.Values{}
	form.Set("email", c.Email)
	form.Set("password", c.password)
	form.Set("companyName", c.Company)

	req, err := http.NewRequest("POST", rurl.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("Failed to create new request obj: %s", err.Error())
	}

	// set required header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// dispatch request and check response status
	bs, err := c.dispatchRequest(req, http.StatusOK)
	if err != nil {
		return nil, err
	}

	// marshall JSON response into object
	var tokens []accessToken
	if err := json.Unmarshal(bs, &tokens); err != nil {
		return nil, fmt.Errorf("Failed to marshall response: %s", err.Error())
	}

	// search for required company token
	for _, t := range tokens {
		if t.OrganizationName == c.Company {
			c.token = &t
			return &t, nil
		}
	}

	return nil, fmt.Errorf("Unable to find token for '%s'\nTokens returned: %+v\n", c.Company, tokens)
}

// resourceURL safely joins a string path (e.g. "my/resource") to an existing URL.
func (c *Client) resourceURL(path string) *url.URL {
	rel := &url.URL{Path: path}
	return c.BaseURL.ResolveReference(rel)
}

// dispatchRequest to server and return a byte slice containing the response body.
// An error will be returned instead if the request fails or if the response
// status does not match the expected one.
func (c *Client) dispatchRequest(req *http.Request, status int) ([]byte, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request failed: %v", err)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body) // can't fail because it's reading in memory

	if resp.StatusCode != status {
		return nil, fmt.Errorf("Server failed to process request: status: %d body: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
