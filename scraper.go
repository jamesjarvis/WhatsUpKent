package main

import (
	"log"

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
	config := scrape.InitialConfig{
		StartRange: 130000,
		EndRange:   134000,
	}

	// Update locations
	errLoc := scrape.ScrapeLocations(client)
	if errLoc != nil {
		log.Fatal(errLoc)
	}
	log.Println("------------- Location scraping complete -------------")

	// TODO: Update Modules

	// Update the ical feeds
	scrape.FuckIt(&config, client)
	log.Println("------------- Event scraping complete -------------")
}
