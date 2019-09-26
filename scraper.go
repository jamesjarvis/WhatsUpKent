package main

import (
	"log"

	"github.com/jamesjarvis/WhatsUpKent/pkg/db"
)

func main() {
	// scrape.FuckIt()

	// Setup database connection
	client := db.NewClient()
	err := db.Setup(client)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Schema successfully updated")
}
