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
func (config *Config) Query() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read request body and close it
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		// HandleError(err)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		defer r.Body.Close()

		//Retrieve query result
		result, err := config.PerformCachedQuery(string(body))
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.WriteHeader(500)
			errorJSON := ErrorJSON{Status: "Query Failed", Error: "Query not correctly formatted."}
			marshalled, marshallErr := json.Marshal(errorJSON)
			HandleError(marshallErr)
			fmt.Fprintf(w, string(marshalled))
		} else {
			fmt.Fprintf(w, *result)
		}
	}
}

//GetQuery performs a read only query without badger
func (config *Config) GetQuery() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read request body and close it
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		HandleError(err)
		defer r.Body.Close()

		//Retrieve query result
		result, err := config.PerformQuery(string(body))
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
}

//Info simply returns a pretty ASCII art
func Info(w http.ResponseWriter, r *http.Request) {
	str :=
		`
__          ___           _       _    _       _  __          _
\ \        / / |         | |     | |  | |     | |/ /         | |
 \ \  /\  / /| |__   __ _| |_ ___| |  | |_ __ | ' / ___ _ __ | |_
  \ \/  \/ / | '_ \ / _' | __/ __| |  | | '_ \|  < / _ \ '_ \| __|
   \  /\  /  | | | | (_| | |_\__ \ |__| | |_) | . \  __/ | | | |_
    \/  \/   |_| |_|\__,_|\__|___/\____/| .__/|_|\_\___|_| |_|\__|
                                        | |
                                        |_|
Welcome.
You're in the wrong area.
To get back to seeing everything going on at the University of Kent, go back to:
https://whatsupkent.com

If you are still curious about the project, find my contact details on https://jamesjarvis.io
`

	fmt.Fprintf(w, str)
}
