package scrape

import (
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/apognu/gocal"
	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/y"
	"github.com/jamesjarvis/WhatsUpKent/pkg/db"
)

// The purpose of this section is to deconstruct the cached ical files and to create event objects from that.
// Then it should call the

// Check if database is already populated, if it is then try and get the next id that needs to be updated.
// If the database is unpopulated, then just use the list of 50,000 initial ids
// Send them into the pool of workers
// Only actually download 2 or 3 at once
// Send the actual filename over to another pool of workers
// This pool has no limit, it shuold read the file, deconstruct the ical into individual events and then add it to the database, after that, it should delete the cached version

// ParseCal opens the file and starts the parsing
func ParseCal(c *dgo.Dgraph, fid *FilesIds, mx *sync.Mutex) error {
	f, _ := os.Open(fid.filename)
	defer f.Close()

	// start, end := time.Now(), time.Now().Add(12*30*24*time.Hour)

	parser := gocal.NewParser(f)
	// c.Start, c.End = &start, &end
	parseErr := parser.Parse()
	if parseErr != nil {
		return parseErr
	}

	currentTime := time.Now()

	scrapeEvent := db.Scrape{
		ID:          fid.id,
		LastScraped: &currentTime,
	}

	currentScrape, err := db.GetScrape(c, scrapeEvent)
	if err != nil {
		return err
	}

	events := make([]db.Event, 0)

	for _, e := range parser.Events {
		event, genErr := generateEvent(c, fid, &e, mx) //Currently getting an int error with the thing
		if genErr != nil {
			return genErr
		}
		tempEvent := db.Event{
			UID: event.UID,
		}
		events = append(events, tempEvent)
	}

	if currentScrape != nil {
		scrapeEvent.UID = currentScrape.UID
	}
	scrapeEvent.FoundEvent = events
	_, upsertErr := db.UpsertScrape(c, scrapeEvent)
	if upsertErr != nil {
		return upsertErr
	}

	return nil
}

func generateEvent(c *dgo.Dgraph, fid *FilesIds, scrapedEvent *gocal.Event, mx *sync.Mutex) (*db.Event, error) {
	eventID, idErr := generateEventID(fid, scrapedEvent.Uid)
	if idErr != nil {
		return nil, idErr
	}

	locations := make([]db.Location, 0)
	loc, locErr := db.GetLocationFromKentSlug(c, scrapedEvent.Location)
	if locErr != nil {
		return nil, locErr
	}
	if loc != nil {
		// tempLoc := db.Location{
		// 	UID: loc.UID,
		// }
		// locations = append(locations, tempLoc)
		locations = append(locations, *loc)
	}

	event := db.Event{
		ID:          eventID, //Sort this out
		Title:       scrapedEvent.Summary,
		Description: scrapedEvent.Description,
		StartDate:   scrapedEvent.Start,
		EndDate:     scrapedEvent.End,
		Location:    locations,
	}

	currentEvent, getErr := db.GetEvent(c, event)
	if getErr != nil {
		return nil, getErr
	}
	if currentEvent != nil {
		event.UID = currentEvent.UID
		//Check if the event is basically the same
		//If it is, then dont bother upserting it.
		if event.Equal(*currentEvent) {
			return currentEvent, nil
		}
	}
	upsertErr := UpsertEventRetryable(c, event, false)
	if upsertErr != nil {
		if upsertErr == y.ErrAborted { //If it aborted due to race condition
			//lock
			mx.Lock()
			upsertErrRetried := UpsertEventRetryable(c, event, true)
			mx.Unlock()
			//unlock
			if upsertErrRetried != nil {
				return nil, upsertErrRetried
			}
			//GOOD EXIT
			return db.GetEvent(c, event)
		}
		return nil, upsertErr
	}

	//Exits here if it created a new event, and has then retrieved that event from the database
	return db.GetEvent(c, event)
}

// UpsertEventRetryable provides a method to recursively keep retrying to upsert a specific event
// This is intended to be used within a mutual exclusion lock to prevent race conditions
func UpsertEventRetryable(c *dgo.Dgraph, e db.Event, retry bool) error {
	_, er1 := db.UpsertEvent(c, e)
	if er1 != nil {
		if er1 == y.ErrAborted && retry {
			return UpsertEventRetryable(c, e, retry)
		}
		return er1
	}
	return nil
}

func generateEventID(fid *FilesIds, currentID string) (string, error) {
	id := fid.id
	r1, err1 := regexp.Compile(strconv.Itoa(id) + "_")
	if err1 != nil {
		return "", err1
	}
	r2, err2 := regexp.Compile("@kent.ac.uk")
	if err2 != nil {
		return "", err2
	}
	temp1 := r1.ReplaceAllLiteralString(currentID, "")
	temp2 := r2.ReplaceAllLiteralString(temp1, "")
	return temp2, nil
}
