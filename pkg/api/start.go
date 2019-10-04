package api

import (
	"log"
	"net/http"
)

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Start() {
	router := NewRouter()

	log.Println("Starting api service on port 4000 .......")

	log.Fatal(http.ListenAndServe(":4000", router))
}
