package gremlin

import "time"

// accessToken represents the object returned by the Gremlin API when making an
// authentication request.
type accessToken struct {
	ID               string    `json:"identifier"`
	Header           string    `json:"header"`
	OrganizationID   string    `json:"org_id"`
	OrganizationName string    `json:"org_name"`
	Token            string    `json:"token"`
	RenewToken       string    `json:"renew_token"`
	Role             string    `json:"role"`
	ExpiresAt        time.Time `json:"expires_at"`
}
