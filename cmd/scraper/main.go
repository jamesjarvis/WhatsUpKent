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
	duration, durationErr := time.ParseDuration("1m")
	if durationErr != nil {
		log.Fatal(durationErr)
	}
	config := scrape.InitialConfig{
		StartRange:   130000,
		EndRange:     134000,
		SlowInterval: duration,
	}

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

	// Now the main scrape is complete, enter a "slow mode"
	scrape.Continuous(&config, client)
}
