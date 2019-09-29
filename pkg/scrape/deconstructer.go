package scrape

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/apognu/gocal"
	"github.com/dgraph-io/dgo/v2"
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
func ParseCal(c *dgo.Dgraph, fid *FilesIds) error {
	f, _ := os.Open(fid.filename)
	defer f.Close()

	// start, end := time.Now(), time.Now().Add(12*30*24*time.Hour)

	parser := gocal.NewParser(f)
	// c.Start, c.End = &start, &end
	parser.Parse()

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
		event, err := generateEvent(c, fid, &e) //Currently getting an int error with the thing
		if err != nil {
			log.Fatal(err)
		}
		events = append(events, *event)
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

func generateEvent(c *dgo.Dgraph, fid *FilesIds, scrapedEvent *gocal.Event) (*db.Event, error) {
	eventID, err := generateEventID(fid, scrapedEvent.Uid)
	if err != nil {
		return nil, err
	}
	event := db.Event{
		ID:          eventID, //Sort this out
		Title:       scrapedEvent.Summary,
		Description: scrapedEvent.Description,
		StartDate:   scrapedEvent.Start,
		EndDate:     scrapedEvent.End,
	}

	currentEvent, err := db.GetEvent(c, event)
	if err != nil {
		return nil, err
	}
	if currentEvent != nil {
		event.UID = currentEvent.UID
	}
	_, er1 := db.UpsertEvent(c, event)
	if er1 != nil {
		return nil, err
	}
	//Exits here if it created a new event, and has then retrieved that event from the database
	return db.GetEvent(c, event)
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
