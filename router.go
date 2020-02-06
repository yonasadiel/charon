package main

import (
	"github.com/gorilla/mux"
)

// CreateRouter returns the router
func CreateRouter() (router *mux.Router) {
	router = mux.NewRouter()

	return router
}
