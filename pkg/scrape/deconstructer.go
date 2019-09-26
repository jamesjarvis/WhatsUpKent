package scrape

import (
	"fmt"
	"os"
	"time"

	"github.com/apognu/gocal"
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

// Opens the file and starts the parsing
func ParseCal(fid FilesIds) {
	f, _ := os.Open(fid.filename)
	defer f.Close()

	// start, end := time.Now(), time.Now().Add(12*30*24*time.Hour)

	c := gocal.NewParser(f)
	// c.Start, c.End = &start, &end
	c.Parse()

	scrapeEvent := db.Scrape{
		ID:          fid.id,
		LastScraped: time.Now(),
	}

	for _, e := range c.Events {
		fmt.Printf("%s on %s by %s", e.Summary, e.Start, e.Organizer.Cn)
	}
}
