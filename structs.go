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

// Attack command details
type Command struct {
	// Type should be one of the following: blackhole, cpu, io, latency, memory,
	// package_loss, shutdown, dns, time_travel, disk, process_killer
	Type string `json:"type"`

	// Args supplied to the command should be identical to those passed to the CLI
	// and UI. Please note that some commands have no required args.
	Args []string `json:"args,omitempty"`
}

// Attack target details
type Target struct {
	// Type should either be "Random" or "Exact"
	Type string `json:"type"`

	// Exact list of hosts to target
	Exact []string `json:"exact,omitempty"`

	// Tags restrict an attack only to hosts with the corresponding kv tags.
	Tags map[string]string `json:"tags,omitempty"`
}

// Encapsulates the details required to launch an attack
type AttackCommand struct {
	Command `json:"command"`
	Target  `json:"target"`

	// Labels are used to target Docker containers running on target hosts
	Labels map[string]string `json:"labels,omitempty"`
}
