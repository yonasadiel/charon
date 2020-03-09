package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yonasadiel/helios"

	"github.com/yonasadiel/charon/auth"
	"github.com/yonasadiel/charon/problem"
)

// CreateRouter returns the router
func CreateRouter() (router *mux.Router) {
	router = mux.NewRouter()
	allowedOrigins := []string{"*"}

	headerMiddleware := func(f helios.HTTPHandler) helios.HTTPHandler {
		return func(req helios.Request) {
			req.SetHeader("Access-Control-Max-Age", "86400")
			req.SetHeader("Access-Control-Allow-Headers", "Content-Type")
			req.SetHeader("Access-Control-Allow-Credentials", "true")
			f(req)
		}
	}

	basicMiddlewares := []helios.Middleware{helios.CreateCORSMiddleware(allowedOrigins), headerMiddleware}
	loggedInMiddlewares := []helios.Middleware{auth.LoggedInMiddleware}

	optionHandler := func(req helios.Request) {
		// do nothing
	}

	router.HandleFunc("/auth/login/", helios.WithMiddleware(auth.LoginView, basicMiddlewares)).Methods(http.MethodPost)
	router.HandleFunc("/auth/login/", helios.WithMiddleware(optionHandler, basicMiddlewares)).Methods(http.MethodOptions)
	router.HandleFunc("/auth/logout/", helios.WithMiddleware(auth.LogoutView, loggedInMiddlewares)).Methods(http.MethodPost)

	router.HandleFunc("/problem/question/", helios.WithMiddleware(problem.QuestionListView, loggedInMiddlewares)).Methods(http.MethodGet)
	router.HandleFunc("/problem/question/:questionId/", helios.WithMiddleware(problem.QuestionDetailView, loggedInMiddlewares)).Methods(http.MethodGet)
	router.HandleFunc("/problem/question/:questionId/submit/", helios.WithMiddleware(problem.SubmissionCreateView, loggedInMiddlewares)).Methods(http.MethodPost)

	router.Use(mux.CORSMethodMiddleware(router))

	return router
}
