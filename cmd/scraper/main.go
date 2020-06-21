package main

import (
	"log"
	"os"
	"time"

	"github.com/jamesjarvis/WhatsUpKent/pkg/db"
	"github.com/jamesjarvis/WhatsUpKent/pkg/scrape"
)

func main() {
	// Setup Scraper
	url := os.Getenv("DGRAPH_URL")
	if url == "" {
		url = "localhost:9080"
	}

	// Setup database connection
	log.Println("Setting up DB Connection")
	client, err := db.NewClient(url)
	if err != nil {
		log.Fatal(err)
	}

	config := scrape.InitialConfig{
		Url:              url,
		StartRange:       110000,
		EndRange:         150000,
		SlowInterval:     time.Second * 30,
		MaxAge:           time.Hour * 24 * 7,
		DownloadPool:     1,
		ProcessPool:      3,
		EventProcessPool: 5,
		DBClient:         client,
	}

	log.Println("Install schema into DB")
	err = config.DBClient.Setup()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Schema successfully updated")

	s, errOld := config.DBClient.GetOldestScrape()
	if errOld != nil {
		log.Fatal(errOld)
	}
	oldestAge := time.Since(*s.LastScraped)
	//Only enter the big scraping if the oldest scrape is over a week old. Helps if it ever crashes (shouldnt do!)
	if oldestAge > config.MaxAge {
		// Update locations
		errLoc := config.Locations()
		if errLoc != nil {
			log.Fatal(errLoc)
		}
		log.Println("------------- Location scraping complete -------------")

		// Update Modules
		errMod := config.Modules()
		if errMod != nil {
			log.Fatal(errMod)
		}
		log.Println("------------- Module scraping complete -------------")

		// Update the ical feeds
		config.FuckIt()
		log.Println("------------- Event scraping complete -------------")
	}

	//Increase the number of processes available to scrape
	config.EventProcessPool = 10

	// Now the main scrape is complete, enter a "slow mode"
	continuousErr := config.Continuous()
	if continuousErr != nil {
		log.Fatal(continuousErr)
	}
}
