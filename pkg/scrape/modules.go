package scrape

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dgraph-io/dgo/v2"
	"github.com/jamesjarvis/WhatsUpKent/pkg/db"
)

// This file is going to scrape all of the modules it can find from the kent API
// Then save it to the database

//types

//ModuleInfo is one of the objects the kent API returns
//Most of the fields are empty, and are all string types, because thank you kent
type ModuleInfo struct {
	Code    string `json:"code,omitempty"`
	Running bool   `json:"running,omitempty"`
	SDSCode string `json:"sds_code,omitempty"`
	Title   string `json:"title,omitempty"`
}

type ModuleAPIResult struct {
	Modules       map[string]ModuleInfo `json:"modules,omitempty"`
	Title         string                `json:"title,omitempty"`
	Total         int                   `json:"total,omitempty"`
	TotalFiltered int                   `json:"total_filtered,omitempty"`
}

// //UnmarshalJSON only exists to fix errors with unmarshaling the json from the api
// func (m *ModuleInfo) UnmarshalJSON(d []byte) error {
// 	type T2 LocationInfo // create new type with same structure as T but without its method set!
// 	x := struct {
// 		T2                   // embed
// 		Capacity json.Number `json:"capacity"`
// 	}{T2: T2(*l)} // don't forget this, if you do and 't' already has some fields set you would lose them

// 	if err := json.Unmarshal(d, &x); err != nil {
// 		return err
// 	}
// 	*l = LocationInfo(x.T2)
// 	l.Capacity = x.Capacity.String()
// 	return nil
// }

func downloadAndMarshalModules() (*[]ModuleInfo, error) {
	url := "https://api.kent.ac.uk/api/v1/modules"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	var result ModuleAPIResult
	jsonErr := json.Unmarshal([]byte(body), &result)
	if jsonErr != nil {
		return nil, jsonErr
	}
	modules := make([]ModuleInfo, 0)
	for _, val := range result.Modules {
		modules = append(modules, val)
	}

	return &modules, nil
}

//ScrapeModules scrapes the modules from kent api if they dont already exist
func ScrapeModules(c *dgo.Dgraph) error {
	n, countErr := db.CountNodesWithFieldUnsafe(c, "module.code")
	if countErr != nil {
		return countErr
	}
	if *n == 0 {
		apiModules, apiErr := downloadAndMarshalModules()
		if apiErr != nil {
			return apiErr
		}

		for _, m := range *apiModules {
			tempMod := db.Module{
				Code: m.SDSCode,
				Name: m.Title,
			}

			checkExist, existErr := db.GetModuleFromSDSCode(c, m.SDSCode)
			if existErr != nil {
				return existErr
			}
			if checkExist == nil {
				_, er1 := db.UpsertModule(c, tempMod)
				if er1 != nil {
					return er1
				}
			}
		}
	}

	return nil
}
