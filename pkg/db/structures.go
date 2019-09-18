package db

type Module struct {
	ID      string   `json:"id,omitempty"`
	Name    string   `json:"name,omitempty"`
	Subject string   `json:"subject,omitempty"`
	URL     string   `json:"url,omitempty"`
	DType   []string `json:"dgraph.type,omitempty"`
}

// type
