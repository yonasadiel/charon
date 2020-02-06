package main

import (
	"github.com/gorilla/mux"

	"github.com/yonasadiel/charon/app"
	"github.com/yonasadiel/charon/auth"
)

// CreateRouter returns the router
func CreateRouter() (router *mux.Router) {
	router = mux.NewRouter()

	router.HandleFunc("/login", app.Middleware(auth.LoginView)).Methods("POST")

	return router
}
