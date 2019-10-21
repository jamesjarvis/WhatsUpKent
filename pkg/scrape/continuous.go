package scrape

import (
	"log"
	"sync"
	"time"

	"github.com/dgraph-io/dgo/v2"
	"github.com/jamesjarvis/WhatsUpKent/pkg/db"
)

//Continuous is the continous scraper
func Continuous(config *InitialConfig, c *dgo.Dgraph) error {
	var eventMX = &sync.Mutex{}

	for {
		time.Sleep(config.SlowInterval)

		//Get oldest scrape
		oldestScrape, oldErr := db.GetOldestScrape(c)
		if oldErr != nil {
			return oldErr
		}

		//Download file
		fid := FilesIds{
			id:       oldestScrape.ID,
			filename: FormatFilename(oldestScrape.ID),
		}
		err := DownloadFile(fid)

		if err != nil {
			if err != ErrInvalidID {
				return err
			}
			//Remove the dead scrape
			log.Printf("Scrape %d seems dead, removing from database...", fid.id)
			removeScrapeErr := db.RemoveScrape(c, *oldestScrape)
			if removeScrapeErr != nil {
				return removeScrapeErr
			}
		} else {
			//Scrape file
			err = ProcessFile(c, fid, eventMX, config)
			duration := time.Since(*oldestScrape.LastScraped)
			log.Printf("Rescraped %d, after %s minutes", fid.id, duration)
		}
	}
	return nil
}
