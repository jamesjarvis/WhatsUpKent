package main

import (
	"log"
	"os"

	"github.com/jamesjarvis/WhatsUpKent/pkg/api"
)

func main() {
	url := os.Getenv("DGRAPH_URL")
	if url == "" {
		url = "localhost:9080"
	}

	err := api.Starter(url)
	if err != nil {
		log.Fatal(err)
	}
}
