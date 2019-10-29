package api

import (
	"log"
	"net/http"
	"sync"

	badger "github.com/dgraph-io/badger"
	"github.com/dgraph-io/dgo/v2"
	"github.com/jamesjarvis/WhatsUpKent/pkg/db"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
func Starter(url string) {
	URL = url
	Lock = &sync.Mutex{}
	var err error

	// Set up a new DB client
	Client = db.NewClient(URL)

	// Set up a new cache client
	CacheDB, err = badger.Open(badger.DefaultOptions("/cache"))
	HandleError(err)
	defer CacheDB.Close()

	router := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST"})
	origins := handlers.AllowedOrigins([]string{"*"})

	router.HandleFunc("/", Info).Methods("GET")
	router.HandleFunc("/", Query).Methods("POST")

	log.Println("🤖 Starting api service on port 4000 .......")
	HandleError(http.ListenAndServe(":4000", handlers.CORS(headers, methods, origins)(router)))
}
