package main

import (
	"os"

	"github.com/jamesjarvis/WhatsUpKent/pkg/api"
)

func main() {
	url := os.Getenv("DGRAPH_URL")
	if url == "" {
		url = "localhost:9080"
	}
	api.Start(url)
}
