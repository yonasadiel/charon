package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yonasadiel/helios"

	"github.com/yonasadiel/charon/backend/auth"
	"github.com/yonasadiel/charon/backend/exam"
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
	loggedInMiddlewares := []helios.Middleware{helios.CreateCORSMiddleware(allowedOrigins), headerMiddleware, auth.LoggedInMiddleware}

	optionHandler := func(req helios.Request) {
		// do nothing
	}

	router.HandleFunc("/auth/login/", helios.WithMiddleware(auth.LoginView, basicMiddlewares)).Methods(http.MethodPost)
	router.HandleFunc("/auth/login/", helios.WithMiddleware(optionHandler, basicMiddlewares)).Methods(http.MethodOptions)
	router.HandleFunc("/auth/logout/", helios.WithMiddleware(auth.LogoutView, loggedInMiddlewares)).Methods(http.MethodPost)
	router.HandleFunc("/auth/logout/", helios.WithMiddleware(optionHandler, basicMiddlewares)).Methods(http.MethodOptions)
	router.HandleFunc("/auth/user/", helios.WithMiddleware(auth.UserListView, loggedInMiddlewares)).Methods(http.MethodGet)
	router.HandleFunc("/auth/user/", helios.WithMiddleware(auth.UserCreateView, loggedInMiddlewares)).Methods(http.MethodPost)
	router.HandleFunc("/auth/user/", helios.WithMiddleware(optionHandler, basicMiddlewares)).Methods(http.MethodOptions)

	router.HandleFunc("/exam/venue/", helios.WithMiddleware(exam.VenueListView, loggedInMiddlewares)).Methods(http.MethodGet)
	router.HandleFunc("/exam/venue/", helios.WithMiddleware(exam.VenueCreateView, loggedInMiddlewares)).Methods(http.MethodPost)
	router.HandleFunc("/exam/venue/", helios.WithMiddleware(optionHandler, basicMiddlewares)).Methods(http.MethodOptions)
	router.HandleFunc("/exam/venue/{venueID}/", helios.WithMiddleware(exam.VenueDeleteView, loggedInMiddlewares)).Methods(http.MethodDelete)
	router.HandleFunc("/exam/venue/{venueID}/", helios.WithMiddleware(optionHandler, basicMiddlewares)).Methods(http.MethodOptions)
	router.HandleFunc("/exam/participation/", helios.WithMiddleware(exam.ParticipationListView, loggedInMiddlewares)).Methods(http.MethodGet)
	router.HandleFunc("/exam/participation/", helios.WithMiddleware(exam.ParticipationCreateView, loggedInMiddlewares)).Methods(http.MethodPost)
	router.HandleFunc("/exam/participation/", helios.WithMiddleware(optionHandler, basicMiddlewares)).Methods(http.MethodOptions)
	router.HandleFunc("/exam/participation/{participationID}/", helios.WithMiddleware(exam.ParticipationDeleteView, loggedInMiddlewares)).Methods(http.MethodDelete)
	router.HandleFunc("/exam/participation/{participationID}/", helios.WithMiddleware(optionHandler, basicMiddlewares)).Methods(http.MethodOptions)
	router.HandleFunc("/exam/", helios.WithMiddleware(exam.EventListView, loggedInMiddlewares)).Methods(http.MethodGet)
	router.HandleFunc("/exam/", helios.WithMiddleware(exam.EventCreateView, loggedInMiddlewares)).Methods(http.MethodPost)
	router.HandleFunc("/exam/", helios.WithMiddleware(optionHandler, basicMiddlewares)).Methods(http.MethodOptions)
	router.HandleFunc("/exam/{eventSlug}/question/", helios.WithMiddleware(exam.QuestionListView, loggedInMiddlewares)).Methods(http.MethodGet)
	router.HandleFunc("/exam/{eventSlug}/question/", helios.WithMiddleware(exam.QuestionCreateView, loggedInMiddlewares)).Methods(http.MethodPost)
	router.HandleFunc("/exam/{eventSlug}/question/", helios.WithMiddleware(optionHandler, basicMiddlewares)).Methods(http.MethodOptions)
	router.HandleFunc("/exam/{eventSlug}/question/{questionID}/", helios.WithMiddleware(exam.QuestionDetailView, loggedInMiddlewares)).Methods(http.MethodGet)
	router.HandleFunc("/exam/{eventSlug}/question/{questionID}/", helios.WithMiddleware(exam.QuestionDeleteView, loggedInMiddlewares)).Methods(http.MethodDelete)
	router.HandleFunc("/exam/{eventSlug}/question/{questionID}/", helios.WithMiddleware(optionHandler, basicMiddlewares)).Methods(http.MethodOptions)
	router.HandleFunc("/exam/{eventSlug}/question/{questionID}/submit/", helios.WithMiddleware(exam.SubmissionCreateView, loggedInMiddlewares)).Methods(http.MethodPost)
	router.HandleFunc("/exam/{eventSlug}/question/{questionID}/submit/", helios.WithMiddleware(optionHandler, basicMiddlewares)).Methods(http.MethodOptions)

	router.Use(mux.CORSMethodMiddleware(router))

	return router
}
