package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	mux "github.com/julienschmidt/httprouter"
)

//Query performs a read only query
func Query(w http.ResponseWriter, r *http.Request, _ mux.Params) {

	// Read request body and close it
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	HandleError(err)
	defer r.Body.Close()

	//Retrieve query result
	result, err := PerformCachedQuery(string(body))
	HandleError(err)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, *result)
}
