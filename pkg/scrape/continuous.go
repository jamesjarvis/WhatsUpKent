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
	var processWG sync.WaitGroup
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
			return err
		}

		//Scrape file
		processWG.Add(1)
		go ProcessFile(c, fid, eventMX, &processWG)
		duration := time.Since(*oldestScrape.LastScraped)
		log.Printf("Rescraping %d, after %s minutes", fid.id, duration)
	}
	processWG.Wait()
	return nil
}