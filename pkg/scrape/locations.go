package scrape

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dgraph-io/dgo/v2"
	"github.com/jamesjarvis/WhatsUpKent/pkg/db"
)

// This file is going to scrape all of the locations it can find from the kent API
// Then save it to the database

//types

//LocationInfo is one of the objects the kent API returns
//Most of the fields are empty, and are all string types, because thank you kent
type LocationInfo struct {
	BookBy         string `json:"book_by,omitempty"`
	Campus         string `json:"campus,omitempty"`
	CampusID       string `json:"campus_id,omitempty"`
	Capacity       string `json:"capacity,omitempty"`
	Classification string `json:"classification,omitempty"`
	Department     string `json:"department,omitempty"`
	Directions     string `json:"directions,omitempty"`
	DisabledAccess string `json:"disabled_access,omitempty"`
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Photo          string `json:"photo,omitempty"`
	Site           string `json:"site,omitempty"`
	SiteID         string `json:"site_id,omitempty"`
	Type           string `json:"type,omitempty"`
	UFName         string `json:"uf_name,omitempty"`
}

//UnmarshalJSON only exists to fix errors with unmarshaling the json from the api
func (l *LocationInfo) UnmarshalJSON(d []byte) error {
	type T2 LocationInfo // create new type with same structure as T but without its method set!
	x := struct {
		T2                   // embed
		Capacity json.Number `json:"capacity"`
	}{T2: T2(*l)} // don't forget this, if you do and 't' already has some fields set you would lose them

	if err := json.Unmarshal(d, &x); err != nil {
		return err
	}
	*l = LocationInfo(x.T2)
	l.Capacity = x.Capacity.String()
	return nil
}

func downloadAndMarshal() (*[]LocationInfo, error) {
	url := "https://api.kent.ac.uk/api/v1/rooms"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	var result map[string]LocationInfo
	jsonErr := json.Unmarshal([]byte(body), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}
	locations := make([]LocationInfo, 0)
	for _, val := range result {
		locations = append(locations, val)
	}

	return &locations, nil
}

func yesNoToBool(s string) bool {
	return s == "Yes" || s == "yes"
}

//Locations scrapes the locations from kent api if they dont already exist
func Locations(c *dgo.Dgraph) error {
	n, countErr := db.CountNodesWithFieldUnsafe(c, "location.id")
	if countErr != nil {
		return countErr
	}
	if *n == 0 {
		apiLocations, apiErr := downloadAndMarshal()
		if apiErr != nil {
			return apiErr
		}

		for _, loc := range *apiLocations {
			tempLoc := db.Location{
				ID:             loc.ID,
				Name:           loc.UFName,
				DisabledAccess: yesNoToBool(loc.DisabledAccess),
			}

			_, er1 := db.UpsertLocation(c, tempLoc)
			if er1 != nil {
				return er1
			}
		}
	}

	return nil
}
