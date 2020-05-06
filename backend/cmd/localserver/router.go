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
	router.HandleFunc("/exam/", helios.WithMiddleware(exam.EventListView, loggedInMiddlewares)).Methods(http.MethodGet)
	router.HandleFunc("/exam/", helios.WithMiddleware(exam.EventCreateView, loggedInMiddlewares)).Methods(http.MethodPost)
	router.HandleFunc("/exam/{eventSlug}/participation/", helios.WithMiddleware(exam.ParticipationListView, loggedInMiddlewares)).Methods(http.MethodGet)
	router.HandleFunc("/exam/{eventSlug}/participation/", helios.WithMiddleware(exam.ParticipationCreateView, loggedInMiddlewares)).Methods(http.MethodPost)
	router.HandleFunc("/exam/{eventSlug}/participation/", helios.WithMiddleware(optionHandler, basicMiddlewares)).Methods(http.MethodOptions)
	router.HandleFunc("/exam/{eventSlug}/participation-status/", helios.WithMiddleware(exam.ParticipationStatusListView, loggedInMiddlewares)).Methods(http.MethodGet)
	router.HandleFunc("/exam/{eventSlug}/participation-status/", helios.WithMiddleware(optionHandler, basicMiddlewares)).Methods(http.MethodOptions)
	router.HandleFunc("/exam/{eventSlug}/participation-status/{sessionID}/", helios.WithMiddleware(exam.ParticipationStatusDeleteView, loggedInMiddlewares)).Methods(http.MethodDelete)
	router.HandleFunc("/exam/{eventSlug}/participation-status/{sessionID}/", helios.WithMiddleware(optionHandler, basicMiddlewares)).Methods(http.MethodOptions)
	router.HandleFunc("/exam/{eventSlug}/verify/", helios.WithMiddleware(exam.ParticipationVerifyView, loggedInMiddlewares)).Methods(http.MethodPost)
	router.HandleFunc("/exam/{eventSlug}/verify/", helios.WithMiddleware(optionHandler, basicMiddlewares)).Methods(http.MethodOptions)
	router.HandleFunc("/exam/{eventSlug}/participation/{participationID}/", helios.WithMiddleware(exam.ParticipationDeleteView, loggedInMiddlewares)).Methods(http.MethodDelete)
	router.HandleFunc("/exam/{eventSlug}/participation/{participationID}/", helios.WithMiddleware(optionHandler, basicMiddlewares)).Methods(http.MethodOptions)
	router.HandleFunc("/exam/", helios.WithMiddleware(optionHandler, basicMiddlewares)).Methods(http.MethodOptions)
	router.HandleFunc("/exam/{eventSlug}/sync/", helios.WithMiddleware(exam.GetSynchronizationDataView, loggedInMiddlewares)).Methods(http.MethodGet)
	router.HandleFunc("/exam/{eventSlug}/sync/", helios.WithMiddleware(exam.PutSynchronizationDataView, loggedInMiddlewares)).Methods(http.MethodPost)
	router.HandleFunc("/exam/{eventSlug}/sync/", helios.WithMiddleware(optionHandler, basicMiddlewares)).Methods(http.MethodOptions)
	router.HandleFunc("/exam/{eventSlug}/decrypt/", helios.WithMiddleware(exam.DecryptEventDataView, loggedInMiddlewares)).Methods(http.MethodPost)
	router.HandleFunc("/exam/{eventSlug}/decrypt/", helios.WithMiddleware(optionHandler, basicMiddlewares)).Methods(http.MethodOptions)
	router.HandleFunc("/exam/{eventSlug}/question/", helios.WithMiddleware(exam.QuestionListView, loggedInMiddlewares)).Methods(http.MethodGet)
	router.HandleFunc("/exam/{eventSlug}/question/", helios.WithMiddleware(exam.QuestionCreateView, loggedInMiddlewares)).Methods(http.MethodPost)
	router.HandleFunc("/exam/{eventSlug}/question/", helios.WithMiddleware(optionHandler, basicMiddlewares)).Methods(http.MethodOptions)
	router.HandleFunc("/exam/{eventSlug}/question/{questionNumber}/", helios.WithMiddleware(exam.QuestionDetailView, loggedInMiddlewares)).Methods(http.MethodGet)
	router.HandleFunc("/exam/{eventSlug}/question/{questionNumber}/", helios.WithMiddleware(exam.QuestionDeleteView, loggedInMiddlewares)).Methods(http.MethodDelete)
	router.HandleFunc("/exam/{eventSlug}/question/{questionNumber}/", helios.WithMiddleware(optionHandler, basicMiddlewares)).Methods(http.MethodOptions)
	router.HandleFunc("/exam/{eventSlug}/question/{questionNumber}/submit/", helios.WithMiddleware(exam.SubmissionCreateView, loggedInMiddlewares)).Methods(http.MethodPost)
	router.HandleFunc("/exam/{eventSlug}/question/{questionNumber}/submit/", helios.WithMiddleware(optionHandler, basicMiddlewares)).Methods(http.MethodOptions)

	router.Use(mux.CORSMethodMiddleware(router))

	return router
}
