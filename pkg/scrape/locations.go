package scrape

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

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

// tryAndGetTheLocationFromARoom is the hackiest solution to getting the location of a room
// Basically, it scrapes a disgusting ass webpage to get a google maps url and get the lat/long from that
// For this, we fail silently since I'm not sure how long this website will even be up for?
func tryAndGetTheLocationFromARoom(lo *LocationInfo) *db.Loc {
	baseURL, err := url.Parse("https://www.kent.ac.uk/timetabling/rooms/room.html")
	if err != nil {
		return nil
	}
	// Prepare Query Parameters
	params := url.Values{}
	params.Add("room", lo.ID)

	// Add Query Parameters to the URL
	baseURL.RawQuery = params.Encode() // Escape Query Parameters

	// Get the timetable html page
	resp, err := http.Get(baseURL.String())
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil
	}

	// Get the google maps url from this page
	re := regexp.MustCompile(`https:\/\/maps\.google\.co\.uk\/maps[^"]*`)
	result := re.Find(body)
	if result == nil {
		return nil
	}

	// Get the query params from the google maps url, we are looking for "ll"
	u, err := url.Parse(string(result))
	if err != nil {
		return nil
	}

	queries := u.Query()
	ll := queries.Get("ll")
	ls := strings.Split(ll, ",")
	lat, err := strconv.ParseFloat(ls[0], 64)
	if err != nil {
		return nil
	}
	lon, err := strconv.ParseFloat(ls[1], 64)
	if err != nil {
		return nil
	}

	g := &db.Loc{
		Type:   "Point",
		Coords: []float64{lat, lon},
	}

	// g := geom.NewPointFlat(geom.XY, []float64{lat, lon})

	// Finally, if nothing has failed, then return the coordinates
	return g
}

//Locations scrapes the locations from kent api if they dont already exist
func (config *InitialConfig) Locations() error {
	n, countErr := config.DBClient.CountNodesWithFieldUnsafe("location.id")
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
				DType:          []string{"Location"},
			}

			// If it can find the location, then add the damn location
			latlon := tryAndGetTheLocationFromARoom(&loc)
			if latlon != nil {
				tempLoc.Location = *latlon
			}

			_, er1 := config.DBClient.UpsertLocation(tempLoc)
			if er1 != nil {
				return er1
			}
		}
	}

	return nil
}
