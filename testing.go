package gremlin

import (
	"encoding/json"
	"time"

	"github.com/lann/builder"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

// AccessTokenBuilder
type accessTokenBuilder builder.Builder

func (b accessTokenBuilder) OrganizationName(orgName string) accessTokenBuilder {
	return builder.Set(b, "OrganizationName", orgName).(accessTokenBuilder)
}
func (b accessTokenBuilder) Build() accessToken {
	return builder.GetStruct(b).(accessToken)
}
func buildDefaultAccessToken() accessTokenBuilder {
	b := builder.Register(accessTokenBuilder{}, accessToken{}).(accessTokenBuilder)
	b = builder.Set(b, "ID", "fake-id").(accessTokenBuilder)
	b = builder.Set(b, "Header", "fake-header").(accessTokenBuilder)
	b = builder.Set(b, "OrganizationID", "fake-org-id").(accessTokenBuilder)
	b = builder.Set(b, "OrganizationName", "fake-org-name").(accessTokenBuilder)
	b = builder.Set(b, "Token", "fake-token").(accessTokenBuilder)
	b = builder.Set(b, "RenewToken", "fake-renew-token").(accessTokenBuilder)
	b = builder.Set(b, "Role", "fake-role").(accessTokenBuilder)
	b = builder.Set(b, "ExpiresAt", time.Now()).(accessTokenBuilder)
	return b
}

// Mock configuration for auth
func mockSucessAuth(url string, tokenArr []accessToken) {
	httpmock.Activate()

	authResponse, _ := json.Marshal(tokenArr)
	httpmock.RegisterResponder("POST", url+"users/auth", httpmock.NewStringResponder(200, string(authResponse)))
}

func mockFailAuth(url string, httpStatus int) {
	httpmock.Activate()

	httpmock.RegisterResponder("POST", url+"users/auth", httpmock.NewStringResponder(httpStatus, "bad response"))
}

func mockBadResponseStructure(url string) {
	httpmock.Activate()

	httpmock.RegisterResponder("POST", url+"users/auth", httpmock.NewStringResponder(200, string("")))
}
