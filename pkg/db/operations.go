package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
)

// This file should contain methods for interacting with the data easily.
// This includes reading and mutating data for the different data types

// GetScrape should recieve a dgraph client and a scrape struct,
// and return the official scrape struct from the database, complete with Uid for referencing
// if no such struct exists, then it returns an error
func GetScrape(c *dgo.Dgraph, scrape Scrape) (*Scrape, error) {
	if scrape.UID != "" {
		return getScrapeWithID(c, scrape)
	}
	return getScrapeWithoutID(c, scrape)
}

func getScrapeWithID(c *dgo.Dgraph, scrape Scrape) (*Scrape, error) {
	txn := c.NewReadOnlyTxn()
	ctx := context.Background()
	q :=
		`query FindScrape($uid: string) {
			findScrape(func: uid($uid)) {
				uid
				id
				last_scraped
				found_event {
					uid
					id
					title
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

func getScrapeWithoutID(c *dgo.Dgraph, scrape Scrape) (*Scrape, error) {
	txn := c.NewReadOnlyTxn()
	ctx := context.Background()
	q :=
		`query FindScrapeNoID($id: string) {
			findScrapeNoID(func: eq(id, $id)) {
				uid
				id
				last_scraped
				found_event {
					uid
					id
					title
				}
			}
		}
	`
	variables := make(map[string]string)
	variables["$id"] = string(scrape.ID)

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
		return nil, fmt.Errorf("No Scrape found with id %d", scrape.ID)
	}

	return &r.FindScrapeNoID[0], nil
}

// UpsertScrape upserts the scrape struct into the database
func UpsertScrape(c *dgo.Dgraph, scrape Scrape) (string, error) {
	mu := &api.Mutation{
		CommitNow: true,
	}
	ctx := context.Background()
	pb, err := json.Marshal(scrape)
	if err != nil {
		return "", err
	}

	mu.SetJson = pb
	assigned, err := c.NewTxn().Mutate(ctx, mu)
	if err != nil {
		return "", err
	}
	fmt.Print(assigned.Uids)
	return assigned.Uids["blank-0"], nil
}
