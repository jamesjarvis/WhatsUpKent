package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// HandleError simply logs and exits the program if the error exists
func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//ErrorJSON is the response sent back in case of an error
type ErrorJSON struct {
	Status, Error string
}

//Query performs a read only query
func Query(w http.ResponseWriter, r *http.Request) {

	// Read request body and close it
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	HandleError(err)
	defer r.Body.Close()

	//Retrieve query result
	result, err := PerformCachedQuery(string(body))
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(500)
		errorJSON := ErrorJSON{Status: "Query Failed", Error: "Query not correctly formatted."}
		marshalled, marshallErr := json.Marshal(errorJSON)
		HandleError(marshallErr)
		fmt.Fprintf(w, string(marshalled))
		HandleError(err)
	} else {
		fmt.Fprintf(w, *result)
	}
}

//GetQuery performs a read only query without badger
func GetQuery(w http.ResponseWriter, r *http.Request) {
	// Read request body and close it
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	HandleError(err)
	defer r.Body.Close()

	//Retrieve query result
	result, err := PerformQuery(string(body))
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		HandleError(err)
		w.WriteHeader(500)
		errorJSON := ErrorJSON{Status: "Query Failed", Error: "Query probably not correctly formatted."}
		marshalled, marshallErr := json.Marshal(errorJSON)
		HandleError(marshallErr)
		fmt.Fprintf(w, string(marshalled))
	} else {
		fmt.Fprintf(w, *result)
	}
}
