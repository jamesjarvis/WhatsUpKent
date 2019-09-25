package db

import "time"

// Golang schemas

type loc struct {
	Type   string    `json:"type,omitempty"`
	Coords []float64 `json:"coordinates,omitempty"`
}

type Module struct {
	UID     string `json:"uid,omitempty"`
	Code    string `json:"code,omitempty"`
	Name    string `json:"name,omitempty"`
	Subject string `json:"subject,omitempty"`
	// URL     string   `json:"url,omitempty"`
	DType []string `json:"dgraph.type,omitempty"`
}

type Scrape struct {
	UID         string     `json:"uid,omitempty"`
	ID          string     `json:"id,omitempty"`
	LastScraped *time.Time `json:"last_scraped,omitempty"`
	Found       []Event    `json:"found,omitempty"`
	DType       []string   `json:"dgraph.type,omitempty"`
}

type Person struct {
	UID   string   `json:"uid,omitempty"`
	Name  string   `json:"name,omitempty"`
	Email string   `json:"email,omitempty"`
	DType []string `json:"dgraph.type,omitempty"`
}

type Location struct {
	UID            string   `json:"uid,omitempty"`
	Name           string   `json:"name,omitempty"`
	Location       loc      `json:"loc,omitempty"`
	DisabledAccess bool     `json:"disabled_access,omitempty"`
	DType          []string `json:"dgraph.type,omitempty"`
}

type Event struct {
	UID          string     `json:"uid,omitempty"`
	ID           string     `json:"id,omitempty"`
	Title        string     `json:"title,omitempty"`
	Description  string     `json:"description,omitempty"`
	StartDate    *time.Time `json:"start_date,omitempty"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	Organiser    []Person   `json:"organiser,omitempty"`
	PartOfModule []Module   `json:"part_of_module,omitempty"`
	Location     []Location `json:"location,omitempty"`

	DType []string `json:"dgraph.type,omitempty"`
}

// The database schema
var Schema = `
help me 
`
