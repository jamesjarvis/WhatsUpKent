package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/jamesjarvis/WhatsUpKent/pkg/db"
	mux "github.com/julienschmidt/httprouter"
)

func Query(w http.ResponseWriter, r *http.Request, _ mux.Params) {

	// Read request body and close it
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	HandleError(err)
	defer r.Body.Close()

	//Get client connection
	client := db.NewClient()
	result, queryErr := db.ReadOnly(client, string(body))
	HandleError(queryErr)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(result))
}
