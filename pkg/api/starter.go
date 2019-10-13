package api

import (
	"log"
	"net/http"
	"time"

	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
)

//URL is the url for the dgraph database
var URL = "localhost:9080"

// Starter starts the server
func Starter(url string) {
	URL = url

	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.MFU),
		memory.AdapterWithCapacity(10000000),
	)
	HandleError(err)

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(10*time.Hour),
		cache.ClientWithRefreshKey("opn"),
	)
	HandleError(err)

	handler := http.HandlerFunc(GetQuery)

	http.Handle("/", cacheClient.Middleware(handler))

	log.Println("Starting api service on port 4000 .......")
	HandleError(http.ListenAndServe(":4000", nil))
}
