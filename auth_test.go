package gremlin

import (
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

const defaultURL = "https://api.gremlin.com/v1/"

var AccessTokenBuilder = buildDefaultAccessToken()

var _ = Describe("Auth", func() {
	Describe("NewClient", func() {
		Context("Creating client using all defaults", func() {
			// Given
			orgName := "Bob's Burgers, Inc."
			email := "real-email@google.com"
			password := "secure-password"

			// When
			client := NewClient(orgName, email, password)

			It("should have a company", func() {
				Expect(client.Company).To(Equal(orgName))
			})

			It("should have a BaseURL set to the default", func() {
				Expect(client.BaseURL.String()).To(Equal(defaultURL))
			})

			It("should have a default timeout of 10", func() {
				Expect(client.client.Timeout).To(Equal(time.Second * 10))
			})
		})

		Context("Creating client using non-default URL", func() {
			// Given
			orgName := "New Co"
			myURL := "http://sweetpuppy.io/"
			email := "real-email@google.com"
			password := "secure-password"

			// When
			client := NewClient(orgName, email, password, WithURL(myURL))

			It("should have a BaseURL set to myURL", func() {
				Expect(client.BaseURL.String()).To(Equal(myURL))
			})
		})

		Context("Creating client with a bad URL", func() {
			// Given
			orgName := "New Co"
			badURL := "://derpa.derp"
			email := "real-email@google.com"
			password := "secure-password"

			It("should cause a panic", func() {
				defer func() {
					Expect(recover()).NotTo(BeNil())
				}()
				// When
				NewClient(orgName, email, password, WithURL(badURL))
			})
		})

		Context("Creating client with non-default HTTPClient", func() {
			// Given
			orgName := "NewCo"
			innerClient := &http.Client{Timeout: 60}
			email := "real-email@google.com"
			password := "secure-password"

			// When
			client := NewClient(orgName, email, password, WithNetClient(innerClient))

			It("should have a new innerClient", func() {
				Expect(client.client).To(Equal(innerClient))
			})
		})
	})
	Describe("authenticate function", func() {
		Context("Test successful authentication", func() {
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

			It("should get a token", func() {
				Expect(client.Token.Token).To(Equal("fake-token"))
			})
		})

		Context("Creating client with wrong URL", func() {
			// Given
			orgName := "Bob's Burgers, Inc."
			myURL := "http://not-real.io/"
			email := "real-email@google.com"
			password := "secure-password"
			client := NewClient(orgName, email, password, WithURL(myURL))

			// When
			_, err := client.Authenticate()

			It("should cause an error to be thrown", func() {
				Expect(err.Error()).To(ContainSubstring("Request failed"))
			})
		})

		Context("Creating client using all defaults and receives 401 auth response", func() {
			// Given
			orgName := "Bob's Burgers, Inc."
			email := "real-email@google.com"
			password := "secure-password"
			client := NewClient(orgName, email, password)

			mockFailAuth(defaultURL, 401)
			defer httpmock.DeactivateAndReset()

			// When
			_, err := client.Authenticate()

			It("should cause an error to be thrown", func() {
				Expect(err.Error()).To(ContainSubstring("status: 401"))
			})
		})

		Context("Creating client using all defaults and fails to marshal auth response", func() {
			// Given
			orgName := "Bob's Burgers, Inc."
			email := "real-email@google.com"
			password := "secure-password"
			client := NewClient(orgName, email, password)

			mockBadResponseStructure(defaultURL)

			// When
			_, err := client.Authenticate()

			It("should cause an error to be thrown", func() {
				Expect(err.Error()).To(ContainSubstring("Failed to marshall response:"))
			})
		})

		Context("Creating client using all defaults and fails to find token for organization", func() {
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

			It("should cause an error to be thrown", func() {
				Expect(err.Error()).To(ContainSubstring("Unable to find token"))
			})
		})
	})
})
