package api

import (
	"log"
	"net/http"
	"sync"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jamesjarvis/WhatsUpKent/pkg/db"
)

// Config is the API configuration
type Config struct {
	//URL is the url for the dgraph database
	URL string
	// DBClient is the database client
	DBClient *db.ConfigDB
	// CacheDB is the cache client
	CacheDB *badger.DB
	// Lock is a global lock for database operations, just makes it a bit nicer.
	Lock *sync.Mutex
}

// SetupRouter returns a router with all the routes attached
func (config *Config) SetupRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", Info).Methods("GET")
	router.HandleFunc("/", config.Query()).Methods("POST")

	return router
}

// Starter starts the server
func Starter(url string) error {

	log.Println("Setting up DB Client")
	// Set up a new DB client
	Client, err := db.NewClient(url)
	if err != nil {
		return err
	}

	log.Println("Setting up Cache client")
	// Set up a new cache client
	CacheDB, err := badger.Open(badger.DefaultOptions("/cache"))
	if err != nil {
		return err
	}

	defer CacheDB.Close()

	config := &Config{
		URL:      url,
		DBClient: Client,
		CacheDB:  CacheDB,
		Lock:     &sync.Mutex{},
	}

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST"})
	origins := handlers.AllowedOrigins([]string{"*"})

	router := config.SetupRouter()

	log.Println("🤖 Starting api service on port 4000 .......")
	return http.ListenAndServe(":4000", handlers.CORS(headers, methods, origins)(router))
}
