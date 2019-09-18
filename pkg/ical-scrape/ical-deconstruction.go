package ical-scrape

// The purpose of this section is to deconstruct the cached ical files and to create event objects from that.
// Then it should call the

// Check if database is already populated, if it is then try and get the next id that needs to be updated.
// If the database is unpopulated, then just use the list of 50,000 initial ids
// Send them into the pool of workers
// Only actually download 2 or 3 at once
// Send the actual filename over to another pool of workers
// This pool has no limit, it shuold read the file, deconstruct the ical into individual events and then add it to the database, after that, it should delete the cached version


import (
	"github.com/apognu/gocal"
)

func parseICal(filename string) {
	
}