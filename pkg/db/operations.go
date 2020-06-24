package db

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/dgraph-io/dgo/v200/protos/api"
)

// This file should contain methods for interacting with the data easily.
// This includes reading and mutating data for the different data types

// GetScrape should recieve a dgraph client and a scrape struct,
// and return the official scrape struct from the database, complete with Uid for referencing
// if no such struct exists, then it returns an error
func (config *ConfigDB) GetScrape(scrape Scrape) (*Scrape, error) {
	if scrape.UID != "" {
		return config.getScrapeWithID(scrape)
	}
	return config.getScrapeWithoutID(scrape)
}

func (config *ConfigDB) getScrapeWithID(scrape Scrape) (*Scrape, error) {
	txn := config.DBClient.NewReadOnlyTxn()
	ctx := context.Background()
	q :=
		`query FindScrape($uid: string) {
			findScrape(func: uid($uid)) {
				uid
				scrape.id
				scrape.last_scraped
				scrape.found_event {
					uid
					event.id
					event.title
				}
			}
		}
	`
	variables := make(map[string]string)
	variables["$uid"] = scrape.UID

	resp, err := txn.QueryWithVars(ctx, q, variables)
	if err != nil {
		return nil, err
	}
	type Root struct {
		FindScrape []Scrape `json:"findScrape"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, err
	}
	if len(r.FindScrape) == 0 {
		return nil, fmt.Errorf("No Scrape found with id %s", scrape.UID)
	}

	return &r.FindScrape[0], nil
}

func (config *ConfigDB) getScrapeWithoutID(scrape Scrape) (*Scrape, error) {
	txn := config.DBClient.NewReadOnlyTxn()
	ctx := context.Background()
	q :=
		`query FindScrapeNoID($id: int) {
			findScrapeNoID(func: eq(scrape.id, $id)) {
				uid
				scrape.id
				scrape.last_scraped
				scrape.found_event {
					uid
					event.id
					event.title
				}
			}
		}
	`
	variables := make(map[string]string)
	variables["$id"] = strconv.Itoa(scrape.ID)

	resp, err := txn.QueryWithVars(ctx, q, variables)
	if err != nil {
		return nil, err
	}
	type Root struct {
		FindScrapeNoID []Scrape `json:"findScrapeNoID"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, err
	}
	if len(r.FindScrapeNoID) == 0 {
		return nil, nil
	}

	return &r.FindScrapeNoID[0], nil
}

// UpsertScrape upserts the scrape struct into the database
func (config *ConfigDB) UpsertScrape(scrape Scrape) (*api.Response, error) {
	mu := &api.Mutation{
		CommitNow: true,
	}
	ctx := context.Background()
	pb, err := json.Marshal(scrape)
	if err != nil {
		return nil, err
	}

	mu.SetJson = pb
	assigned, err := config.DBClient.NewTxn().Mutate(ctx, mu)
	if err != nil {
		return nil, err
	}
	return assigned, nil
}

//RemoveScrape deletes the specified scrape from the database.
func (config *ConfigDB) RemoveScrape(scrape Scrape) error {
	ctx := context.Background()
	d := map[string]string{"uid": scrape.UID}
	pb, err := json.Marshal(d)
	if err != nil {
		return err
	}

	mu := &api.Mutation{
		CommitNow:  true,
		DeleteJson: pb,
	}

	_, err = config.DBClient.NewTxn().Mutate(ctx, mu)
	if err != nil {
		return err
	}
	return nil
}

// GetEvent should recieve a dgraph client and an event struct,
// and return the official event struct from the database, complete with Uid for referencing
// if no such event exists, then it returns an error
func (config *ConfigDB) GetEvent(event Event) (*Event, error) {
	if event.UID != "" {
		return config.getEventWithUID(event)
	}
	return config.getEventWithoutUID(event)
}

func (config *ConfigDB) getEventWithUID(event Event) (*Event, error) {
	txn := config.DBClient.NewReadOnlyTxn()
	ctx := context.Background()
	q :=
		`query FindEvent($id: string) {
			findEvent(func: uid($id)) {
				uid
				event.id
				event.title
				event.description
				event.start_date
				event.end_date
				event.organiser {
					uid
					person.name
				}
				event.part_of_module {
					uid
					module.code
				}
				event.location {
					uid
					location.id
					location.name
				}
			}
		}
	`
	variables := make(map[string]string)
	variables["$id"] = event.UID

	resp, err := txn.QueryWithVars(ctx, q, variables)
	if err != nil {
		return nil, err
	}
	type Root struct {
		FindEvent []Event `json:"findEvent"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, err
	}

	if len(r.FindEvent) == 0 {
		return nil, fmt.Errorf("No Event found with uid %s", event.UID)
	}

	return &r.FindEvent[0], nil
}

func (config *ConfigDB) getEventWithoutUID(event Event) (*Event, error) {
	txn := config.DBClient.NewReadOnlyTxn()
	ctx := context.Background()
	q :=
		`query FindEventNoUID($id: string) {
			findEvent(func: eq(event.id, $id)) {
				uid
				event.id
				event.title
				event.description
				event.start_date
				event.end_date
				event.organiser {
					uid
					person.name
				}
				event.part_of_module {
					uid
					module.code
				}
				event.location {
					uid
					location.id
					location.name
				}
			}
		}
	`
	variables := make(map[string]string)
	variables["$id"] = event.ID

	resp, err := txn.QueryWithVars(ctx, q, variables)
	if err != nil {
		return nil, err
	}
	type Root struct {
		FindEvent []Event `json:"findEvent"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, err
	}
	if len(r.FindEvent) == 0 {
		return nil, nil
	}

	return &r.FindEvent[0], nil
}

// UpsertEvent upserts the event struct into the database
func (config *ConfigDB) UpsertEvent(event Event) (*api.Response, error) {
	mu := &api.Mutation{
		CommitNow: true,
	}
	ctx := context.Background()
	pb, jsonErr := json.Marshal(event)
	if jsonErr != nil {
		return nil, jsonErr
	}

	mu.SetJson = pb
	assigned, upsertErr := config.DBClient.NewTxn().Mutate(ctx, mu)
	if upsertErr != nil {
		// if upsertErr == y.ErrAborted {
		// }
		return nil, upsertErr
	}
	return assigned, nil
}

//GetLocationFromKentSlug returns a matching location from the slug kent uses internally
func (config *ConfigDB) GetLocationFromKentSlug(slug string) (*Location, error) {
	txn := config.DBClient.NewReadOnlyTxn()
	ctx := context.Background()
	q :=
		`query FindLocationFromSlug($id: string) {
			findLocation(func: eq(location.id, $id)) {
				uid
				location.id
				location.name
				location.disabled_access
			}
		}
	`
	variables := make(map[string]string)
	variables["$id"] = slug

	resp, err := txn.QueryWithVars(ctx, q, variables)
	if err != nil {
		return nil, err
	}
	type Root struct {
		FindLocation []Location `json:"findLocation"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, err
	}
	if len(r.FindLocation) == 0 {
		return nil, nil
	}

	return &r.FindLocation[0], nil
}

// UpsertLocation upserts the location struct into the database
func (config *ConfigDB) UpsertLocation(loc Location) (*api.Response, error) {
	mu := &api.Mutation{
		CommitNow: true,
	}
	ctx := context.Background()
	pb, err := json.Marshal(loc)
	if err != nil {
		return nil, err
	}

	mu.SetJson = pb
	assigned, err := config.DBClient.NewTxn().Mutate(ctx, mu)
	if err != nil {
		return nil, err
	}
	return assigned, nil
}

//GetModuleFromSDSCode returns a matching module from the slug kent uses internally, or nil if it doesnt exist
func (config *ConfigDB) GetModuleFromSDSCode(slug string) (*Module, error) {
	txn := config.DBClient.NewReadOnlyTxn()
	ctx := context.Background()
	q :=
		`query FindModuleFromCode($id: string) {
			findModule(func: eq(module.code, $id)) {
				uid
				module.code
				module.name
				module.subject
			}
		}
	`
	variables := make(map[string]string)
	variables["$id"] = slug

	resp, err := txn.QueryWithVars(ctx, q, variables)
	if err != nil {
		return nil, err
	}
	type Root struct {
		FindModule []Module `json:"findModule"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, err
	}
	if len(r.FindModule) == 0 {
		return nil, nil
	}

	return &r.FindModule[0], nil
}

// UpsertModule upserts the location struct into the database
func (config *ConfigDB) UpsertModule(m Module) (*api.Response, error) {
	mu := &api.Mutation{
		CommitNow: true,
	}
	ctx := context.Background()
	pb, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	mu.SetJson = pb
	assigned, err := config.DBClient.NewTxn().Mutate(ctx, mu)
	if err != nil {
		return nil, err
	}
	return assigned, nil
}

// CountNodesWithFieldUnsafe returns the number of nodes which contain the specified field
// this is a good indicator of the number of nodes of a certain type
// this is unsafe, there is no input sanitation and is open to injection attacks
func (config *ConfigDB) CountNodesWithFieldUnsafe(f string) (*int, error) {
	txn := config.DBClient.NewReadOnlyTxn()
	ctx := context.Background()

	q := fmt.Sprintf(
		`query Count {
			nodeCount(func: has(%s)) {
				nodeCount: count(uid)
			}
		}
		`, f)

	resp, err := txn.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	type Root struct {
		NodeCount []struct {
			NodeCount int `json:"nodeCount"`
		} `json:"nodeCount"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, err
	}

	return &r.NodeCount[0].NodeCount, nil
}

//GetOldestScrape retrieves the oldest scrape from the database
func (config *ConfigDB) GetOldestScrape() (*Scrape, error) {
	txn := config.DBClient.NewReadOnlyTxn()
	ctx := context.Background()

	//First, check if there even is anything in the database
	tot, totErr := config.CountNodesWithFieldUnsafe("scrape.id")
	if totErr != nil {
		return nil, totErr
	}
	if *tot == 0 {
		nilTime := time.Unix(0, 0)
		return &Scrape{
			LastScraped: &nilTime,
		}, nil
	}

	q := `{
		oldestScrape(func: type(Scrape), orderasc: scrape.last_scraped, first: 1) {
			uid
			scrape.id
			scrape.last_scraped
		}
	}`

	resp, err := txn.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	type Root struct {
		OldestScrape []Scrape `json:"oldestScrape"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, err
	}
	if len(r.OldestScrape) == 0 {
		return nil, nil
	}

	return &r.OldestScrape[0], nil
}

//ReadOnly is a read only transaction on the database - this is assumed to be ok
func (config *ConfigDB) ReadOnly(q string) ([]byte, error) {
	txn := config.DBClient.NewReadOnlyTxn()
	txn.BestEffort()
	ctx := context.Background()

	resp, err := txn.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	return resp.Json, nil
}
