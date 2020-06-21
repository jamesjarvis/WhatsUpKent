package api

import (
	"log"
	"net/http"
	"sync"

	badger "github.com/dgraph-io/badger"
	"github.com/dgraph-io/dgo/v200"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jamesjarvis/WhatsUpKent/pkg/db"
)

//URL is the url for the dgraph database
var URL = "localhost:9080"

// Client is the database client
var Client *dgo.Dgraph

// CacheDB is the cache client
var CacheDB *badger.DB

// Lock is a global lock for database operations, just makes it a bit nicer.
var Lock *sync.Mutex

// Starter starts the server
func Starter(url string) error {
	URL = url
	Lock = &sync.Mutex{}
	var err error

	log.Println("Setting up DB Client")
	// Set up a new DB client
	Client, err = db.NewClient(URL)
	if err != nil {
		return err
	}

	log.Println("Setting up Cache client")
	// Set up a new cache client
	CacheDB, err = badger.Open(badger.DefaultOptions("/cache"))
	if err != nil {
		return err
	}

	defer CacheDB.Close()

	router := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST"})
	origins := handlers.AllowedOrigins([]string{"*"})

	router.HandleFunc("/", Info).Methods("GET")
	router.HandleFunc("/", Query).Methods("POST")

	log.Println("🤖 Starting api service on port 4000 .......")
	return http.ListenAndServe(":4000", handlers.CORS(headers, methods, origins)(router))
}
