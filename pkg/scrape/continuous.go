package scrape

import (
	"log"
	"sync"
	"time"
)

//Continuous is the continous scraper
func (config *InitialConfig) Continuous() error {
	var eventMX = &sync.Mutex{}

	for {
		time.Sleep(config.SlowInterval)

		//Get oldest scrape
		oldestScrape, oldErr := config.DBClient.GetOldestScrape()
		if oldErr != nil {
			return oldErr
		}

		//Download file
		fid, err := DownloadFile(oldestScrape.ID)

		if err != nil {
			if err != ErrInvalidID {
				return err
			}
			//Remove the dead scrape
			log.Printf("Scrape %d seems dead, removing from database...", oldestScrape.ID)
			removeScrapeErr := config.DBClient.RemoveScrape(*oldestScrape)
			if removeScrapeErr != nil {
				return removeScrapeErr
			}
		} else {
			//Scrape file
			err = config.ProcessFile(fid, eventMX)
			duration := time.Since(*oldestScrape.LastScraped)
			log.Printf("Rescraped %d, after %s minutes", fid.id, duration)
		}
	}
}
