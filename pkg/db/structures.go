package db

import "time"

// Golang schemas

type loc struct {
	Type   string    `json:"type,omitempty"`
	Coords []float64 `json:"coordinates,omitempty"`
}

type Module struct {
	UID     string `json:"uid,omitempty"`
	Code    string `json:"module.code,omitempty"`
	Name    string `json:"module.name,omitempty"`
	Subject string `json:"module.subject,omitempty"`
	// URL     string   `json:"module.url,omitempty"`
	DType []string `json:"dgraph.type,omitempty"`
}

type Scrape struct {
	UID         string     `json:"uid,omitempty"`
	ID          int        `json:"scrape.id,omitempty"`
	LastScraped *time.Time `json:"scrape.last_scraped,omitempty"`
	FoundEvent  []Event    `json:"scrape.found_event,omitempty"`
	DType       []string   `json:"dgraph.type,omitempty"`
}

type Person struct {
	UID   string   `json:"uid,omitempty"`
	Name  string   `json:"person.name,omitempty"`
	Email string   `json:"person.email,omitempty"`
	DType []string `json:"dgraph.type,omitempty"`
}

type Location struct {
	UID            string   `json:"uid,omitempty"`
	ID             string   `json:"location.id,omitempty"`
	Name           string   `json:"location.name,omitempty"`
	Location       loc      `json:"location.loc,omitempty"`
	DisabledAccess bool     `json:"location.disabled_access,omitempty"`
	DType          []string `json:"dgraph.type,omitempty"`
}

type Event struct {
	UID          string     `json:"uid,omitempty"`
	ID           string     `json:"event.id,omitempty"`
	Title        string     `json:"event.title,omitempty"`
	Description  string     `json:"event.description,omitempty"`
	StartDate    *time.Time `json:"event.start_date,omitempty"`
	EndDate      *time.Time `json:"event.end_date,omitempty"`
	Organiser    []Person   `json:"event.organiser,omitempty"`
	PartOfModule []Module   `json:"event.part_of_module,omitempty"`
	Location     []Location `json:"event.location,omitempty"`

	DType []string `json:"dgraph.type,omitempty"`
}

// Schema is the database schema
var Schema = `
event.id: string @index(exact) .
event.title: string @index(fulltext) .
event.start_date: datetime @index(day) .
event.end_date: datetime @index(day).
event.organiser: [uid] @reverse .
event.part_of_module: [uid] @reverse .
event.location: [uid] @reverse .

location.id: string @index(exact) .

module.code: string .
module.name: string @index(fulltext) .
module.subject: string @index(fulltext) .

scrape.found_event: [uid] @reverse .
scrape.id: int @index(int) .

type Loc {
  type: string
  coords: float
}

type Location {
	location.id: string
	location.name: string
	location.loc: Loc
	location.disabled_access: bool
}

type Module {
	module.code: string
	module.name: string
	module.subject: string
}

type Person {
	person.name: string
	person.email: string
}

type Scrape {
	scrape.id: int
	scrape.last_scraped: datetime
	scrape.found_event: [Event]
}

type Event {
	event.id: string
	event.title: string
	event.description: string
	event.start_date: datetime
	event.end_date: datetime
	event.organiser: [Person]
	event.part_of_module: [Module]
	event.location: [Location]
}
`
