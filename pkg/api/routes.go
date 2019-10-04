package api

import (
	mux "github.com/julienschmidt/httprouter"
)

type Route struct {
	Method string
	Path   string
	Handle mux.Handle // httprouter package as mux
}

type Routes []Route

var routes = Routes{
	Route{
		"POST",
		"/query",
		Query,
	},
}
