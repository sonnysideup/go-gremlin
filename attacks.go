package gremlin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// TODO: Implement all attack endpoints
//
// - GET /attacks/active
// - GET /attacks/completed
// - GET /attacks/{guid}
// - DELETE /attacks/{guid}
// - GET /attacks
// - DELETE /attacks

// CreateAttack will launch a new attack in Gremlin against one of your configured
// clients. If the request succeeds, you will receive a UUID for the newly created
// attack.
func (c *Client) CreateAttack(ac AttackCommand) (*uuid.UUID, error) {
	rurl := c.resourceURL("attacks/new")

	attackJSON, err := json.Marshal(ac)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal attack JSON: %v", err)
	}

	req, err := http.NewRequest("POST", rurl.String(), strings.NewReader(string(attackJSON)))
	if err != nil {
		return nil, fmt.Errorf("Failed to create request object: %v", err)
	}

	req.Header.Set("Authorization", c.token.Header)
	req.Header.Set("Content-Type", "application/json")

	bs, err := c.dispatchRequest(req)
	if err != nil {
		return nil, err
	}

	// TODO: understand why FromBytes() fail
	guid, err := uuid.FromString(string(bs))
	if err != nil {
		return nil, fmt.Errorf("Invalid UUID from server: %v", err)
	}

	return &guid, nil
}
