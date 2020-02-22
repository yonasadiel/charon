package main

import (
	"github.com/gorilla/mux"

	"github.com/yonasadiel/helios"
	"github.com/yonasadiel/charon/auth"
	"github.com/yonasadiel/charon/problem"
)

// CreateRouter returns the router
func CreateRouter() (router *mux.Router) {
	router = mux.NewRouter()

	basicMiddlewares := []helios.Middleware{}
	loggedInMiddlewares := []helios.Middleware{auth.LoggedInMiddleware}

	router.HandleFunc("/login", helios.WithMiddleware(auth.LoginView, basicMiddlewares)).Methods("POST")
	router.HandleFunc("/logout", helios.WithMiddleware(auth.LogoutView, loggedInMiddlewares)).Methods("POST")
	router.HandleFunc("/question", helios.WithMiddleware(problem.QuestionListView, loggedInMiddlewares)).Methods("GET")

	return router
}
