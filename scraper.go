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

	scrape.FuckIt(client)
}
