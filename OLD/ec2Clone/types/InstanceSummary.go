package types

type InstanceSummary struct {
	Name      string            `json:"name"`
	Status    string            `json:"status"`
	Type      string            `json:"type"`
	IPv4      []string          `json:"ipv4,omitempty"`
	Profiles  []string          `json:"profiles"`
	CreatedAt string            `json:"created_at"`
	Config    map[string]string `json:"config,omitempty"`
}
