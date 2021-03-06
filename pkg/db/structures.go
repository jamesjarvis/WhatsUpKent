package db

import (
	"time"
)

// Golang schemas

type Loc struct {
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
	Location       Loc      `json:"location.loc,omitempty"`
	DisabledAccess bool     `json:"location.disabled_access"`
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

//Equal checks if the two events are equal
//Does not check UID, as the contents could change
//Does not check the contents of Location, as these are decided at the start
func (e Event) Equal(e2 Event) bool {
	if len(e.Location) != len(e2.Location) {
		return false
	}
	if len(e.Organiser) != len(e2.Organiser) {
		return false
	}
	if len(e.PartOfModule) != len(e2.PartOfModule) {
		return false
	}

	locEqual := true
	for _, loc := range e.Location {
		locEqualTemp := false
		for _, loc2 := range e2.Location {
			locEqualTemp = locEqualTemp || loc.Equal(loc2)
		}
		locEqual = locEqual && locEqualTemp
	}

	modEqual := true
	for _, mod := range e.PartOfModule {
		modEqualTemp := false
		for _, mod2 := range e2.PartOfModule {
			modEqualTemp = modEqualTemp || mod.Equal(mod2)
		}
		modEqual = modEqual && modEqualTemp
	}

	return (e.ID == e2.ID &&
		e.Title == e2.Title &&
		// e.Description == e2.Description &&
		e.StartDate.Equal(*e2.StartDate) &&
		e.EndDate.Equal(*e2.EndDate) &&
		locEqual &&
		modEqual)
}

// Equal returns whether or not the two locations are equivalent
func (l Location) Equal(l2 Location) bool {
	return (l.UID == l2.UID || l.ID == l2.ID || l.Name == l2.Name)
}

// Equal returns whether or not the two modules are equivalent
func (m Module) Equal(m2 Module) bool {
	return (m.UID == m2.UID || m.Code == m2.Code || m.Name == m2.Name)
}

// Schema is the database schema
var Schema = `
location.id: string @index(exact) .
location.name: string .
location.loc: geo .
location.disabled_access: bool .

module.code: string @index(exact) .
module.name: string @index(fulltext) .
module.subject: string @index(fulltext, exact) .

person.name: string .
person.email: string .

scrape.id: int @index(int) .
scrape.last_scraped: datetime .
scrape.found_event: [uid] @reverse .

event.id: string @index(hash) .
event.title: string @index(fulltext, term) .
event.description: string .
event.start_date: datetime @index(hour) .
event.end_date: datetime @index(hour).
event.organiser: [uid] @reverse .
event.part_of_module: [uid] @reverse .
event.location: [uid] @reverse .


type Location {
	location.id: string
	location.name: string
	location.loc: geo
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
