package main

import (
	"log"
	"time"

	"github.com/jamesjarvis/WhatsUpKent/pkg/db"
	"github.com/jamesjarvis/WhatsUpKent/pkg/scrape"
)

func main() {
	// Setup database connection
	client := db.NewClient()
	err := db.Setup(client)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Schema successfully updated")

	// Setup Scraper
	scrape.CreateDownloadDir()
	config := scrape.InitialConfig{
		StartRange:   130000,
		EndRange:     134000,
		SlowInterval: time.Second * 45,
		MaxAge:       time.Hour * 24 * 7,
	}

	s, errOld := db.GetOldestScrape(client)
	if errOld != nil {
		log.Fatal(errOld)
	}
	oldestAge := time.Since(*s.LastScraped)
	//Only enter the big scraping if the oldest scrape is over a week old. Helps if it ever crashes (shouldnt do!)
	if oldestAge > config.MaxAge {
		// Update locations
		errLoc := scrape.Locations(client)
		if errLoc != nil {
			log.Fatal(errLoc)
		}
		log.Println("------------- Location scraping complete -------------")

		// Update Modules
		errMod := scrape.Modules(client)
		if errMod != nil {
			log.Fatal(errMod)
		}
		log.Println("------------- Module scraping complete -------------")

		// Update the ical feeds
		scrape.FuckIt(&config, client)
		log.Println("------------- Event scraping complete -------------")
	}

	// Now the main scrape is complete, enter a "slow mode"
	continuousErr := scrape.Continuous(&config, client)
	if continuousErr != nil {
		log.Fatal(continuousErr)
	}
}
