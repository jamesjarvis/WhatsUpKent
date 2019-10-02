package scrape

import (
	"log"
	"os"
	"regexp"
	"sync"
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
func ParseCal(c *dgo.Dgraph, fid FilesIds, mx *sync.Mutex) error {
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
	eventsChan := make(chan gocal.Event, 10000)
	resultsChan := make(chan db.Event, 10000)
	var wg sync.WaitGroup

	numberOfWorkers := 20

	for i := 0; i <= numberOfWorkers; i++ {
		wg.Add(1)
		go handleGenerator(c, mx, eventsChan, resultsChan, &wg)
	}

	for _, e := range parser.Events {
		eventsChan <- e
	}
	close(eventsChan)

	wg.Wait()
	close(resultsChan)

	for ev := range resultsChan {
		tempEvent := db.Event{
			UID: ev.UID,
		}
		events = append(events, tempEvent)
	}

	log.Printf("Finally scraped %d, with %d events", fid.id, len(events))

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

func handleGenerator(c *dgo.Dgraph, mx *sync.Mutex, eventsChan <-chan gocal.Event, resultsChan chan<- db.Event, wg *sync.WaitGroup) {
	for e := range eventsChan {
		event, genErr := generateEvent(c, &e, mx)
		if genErr != nil {
			log.Fatal(genErr)
		}
		tempEvent := db.Event{
			UID: event.UID,
		}
		resultsChan <- tempEvent
	}
	wg.Done()
}

func generateEvent(c *dgo.Dgraph, scrapedEvent *gocal.Event, mx *sync.Mutex) (*db.Event, error) {
	eventID, idErr := generateEventID(scrapedEvent.Uid)
	if idErr != nil {
		return nil, idErr
	}

	//Locations connecting
	locations := make([]db.Location, 0)
	loc, locErr := db.GetLocationFromKentSlug(c, scrapedEvent.Location)
	if locErr != nil {
		return nil, locErr
	}
	if loc != nil {
		locations = append(locations, *loc)
	} else {
		locations = append(locations, db.Location{
			Name: scrapedEvent.Location,
		})
	}

	//Modules connecting
	modules := make([]db.Module, 0)
	sdsCode, sdsErr := getModuleCodeFromEvent(scrapedEvent.Summary)
	if sdsErr != nil {
		return nil, sdsErr
	}
	mod, modErr := db.GetModuleFromSDSCode(c, sdsCode)
	if modErr != nil {
		return nil, modErr
	}
	if mod != nil {
		modules = append(modules, *mod)
	}

	event := db.Event{
		ID:           eventID, //Sort this out
		Title:        scrapedEvent.Summary,
		Description:  scrapedEvent.Description,
		StartDate:    scrapedEvent.Start,
		EndDate:      scrapedEvent.End,
		Location:     locations,
		PartOfModule: modules,
	}

	//Mutually exclude read,write operations on the database
	mx.Lock()
	storedEvent, storingErr := StoreEvent(c, &event)
	mx.Unlock()
	if storingErr != nil {
		return nil, storingErr
	}
	if storedEvent != nil {
		return storedEvent, nil
	}

	//Exits here if it created a new event, and has then retrieved that event from the database
	return db.GetEvent(c, event)
}

//StoreEvent handles the read and write operations
//Returns the event if it already exists, or nil, with a nil error if it has just been created
func StoreEvent(c *dgo.Dgraph, e *db.Event) (*db.Event, error) {
	currentEvent, getErr := db.GetEvent(c, *e)
	if getErr != nil {
		return nil, getErr
	}
	if currentEvent != nil {
		e.UID = currentEvent.UID
		//Check if the event is basically the same
		//If it is, then dont bother upserting it.
		if e.Equal(*currentEvent) {
			return currentEvent, nil
		}
	}
	_, upsertErr := db.UpsertEvent(c, *e)
	if upsertErr != nil {
		return nil, upsertErr
	}
	return nil, nil
}

func generateEventID(currentID string) (string, error) {
	r1, err1 := regexp.Compile(`\A\d{6}_`)
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

type ModuleMetaInfo struct {
	SDSCode string
	Subject string
}

func getModuleInfoFromEvent(e *gocal.Event) (*ModuleMetaInfo, error) {
	sds, sdsErr := getModuleCodeFromEvent(e.Summary)
	if sdsErr != nil {
		return nil, sdsErr
	}
	subject, subjectErr := getSubjectFromModuleCode(sds)
	if subjectErr != nil {
		return nil, subjectErr
	}
	return &ModuleMetaInfo{
		SDSCode: sds,
		Subject: subject,
	}, nil
}

func getModuleCodeFromEvent(s string) (string, error) {
	r, err := regexp.Compile(`[A-Z]{1,4}\d{1,4}`)
	if err != nil {
		return "", err
	}
	temp := r.FindString(s)
	return temp, nil
}

func getSubjectFromModuleCode(s string) (string, error) {
	r, err := regexp.Compile(`[A-Z]*`)
	if err != nil {
		return "", err
	}
	temp := r.FindString(s)
	subject := SubjectsMap[temp]
	return subject, nil
}
