package gremlin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// defaultBaseURL is the default API location
const defaultBaseURL = "https://api.gremlin.com/v1/"

// Client manages communication with the Gremlin API
type Client struct {
	httpClient *http.Client

	BaseURL     *url.URL
	CompanyName string
}

// Generate a new Gremlin Client, populating the required fields.
func NewClient(company string, httpClient *http.Client) *Client {
	url, err := url.Parse(defaultBaseURL)
	if err != nil {
		log.Fatalln("Failed to parse API URL:", err)
	}

	return &Client{
		BaseURL:     url,
		CompanyName: company,

		// NOTE: add any HTTP-specific configs here, like timeout, transport and redirect settings
		httpClient: &http.Client{},
	}
}

// Authenticate provides your user credentials to Gremlin and requests an access
// token.
//
// All API requests require an access token so you'll need to provide one to all
// other method invocations.
func (c *Client) Authenticate(email string, password string) (*accessToken, error) {
	rurl := resourceURL(c.BaseURL, "users/auth")

	// create request body and object
	form := url.Values{}
	form.Set("email", email)
	form.Set("password", password)
	form.Set("companyName", c.CompanyName)

	req, err := http.NewRequest("POST", rurl.String(), strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatalln("Failed to create new request obj:", err)
	}

	// set required header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// dispatch request and check response status
	resp, err := c.httpClient.Do(req)
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
		if t.OrganizationName == c.CompanyName {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("Unable to find token for '%s'\nTokens returned: %+v\n",
		c.CompanyName, tokens)
}

// func (c *Client) CreateKillProcess() (???, error) {}

// resourceURL safely joins a string path (e.g. "my/resource") to an existing URL.
func resourceURL(base *url.URL, path string) *url.URL {
	rel := &url.URL{Path: path}
	return base.ResolveReference(rel)
}
