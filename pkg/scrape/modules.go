package scrape

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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

//Modules scrapes the modules from kent api if they dont already exist
func (config *InitialConfig) Modules() error {
	n, countErr := config.DBClient.CountNodesWithFieldUnsafe("module.code")
	if countErr != nil {
		return countErr
	}
	if *n == 0 {
		apiModules, apiErr := downloadAndMarshalModules()
		if apiErr != nil {
			return apiErr
		}

		for _, m := range *apiModules {
			subject, subjectErr := getSubjectFromModuleCode(m.SDSCode)
			if subjectErr != nil {
				return subjectErr
			}
			tempMod := db.Module{
				Code:    m.SDSCode,
				Name:    m.Title,
				Subject: subject,
				DType:   []string{"Module"},
			}

			checkExist, existErr := config.DBClient.GetModuleFromSDSCode(m.SDSCode)
			if existErr != nil {
				return existErr
			}
			if checkExist == nil {
				_, er1 := config.DBClient.UpsertModule(tempMod)
				if er1 != nil {
					return er1
				}
			}
		}
	}

	return nil
}
