package exam

import (
	"net/http"

	"github.com/yonasadiel/charon/backend/auth"
	"github.com/yonasadiel/helios"
)

// VenueListView send list of venues
func VenueListView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.GetStatusCode())
		return
	}
	var venues []Venue
	var errGetVenue helios.Error

	venues, errGetVenue = GetAllVenue(user)
	if errGetVenue != nil {
		req.SendJSON(errGetVenue.GetMessage(), errGetVenue.GetStatusCode())
	}

	serializedVenues := make([]VenueData, 0)
	for _, venue := range venues {
		serializedVenues = append(serializedVenues, SerializeVenue(venue))
	}
	req.SendJSON(serializedVenues, http.StatusOK)
}

// VenueCreateView creates the venue
func VenueCreateView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.GetStatusCode())
		return
	}

	var venueData VenueData
	var venue Venue
	var err helios.Error
	err = req.DeserializeRequestData(&venueData)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}
	err = DeserializeVenue(venueData, &venue)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}
	venue.ID = 0
	UpsertVenue(user, &venue)
	req.SendJSON(SerializeVenue(venue), http.StatusCreated)
}

// VenueDeleteView delete the question
func VenueDeleteView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.GetStatusCode())
		return
	}

	venueID, errParseVenueID := req.GetURLParamUint("venueID")
	if errParseVenueID != nil {
		req.SendJSON(errVenueNotFound.GetMessage(), errVenueNotFound.GetStatusCode())
		return
	}

	var venue *Venue
	var err helios.Error
	venue, err = DeleteVenue(user, venueID)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}
	var serializedVenue VenueData = SerializeVenue(*venue)
	req.SendJSON(serializedVenue, http.StatusOK)
}

// EventListView send list of questions
func EventListView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.GetStatusCode())
		return
	}

	var events []Event = GetAllEventOfUser(user)
	serializedEvents := make([]EventData, 0)
	for _, event := range events {
		serializedEvents = append(serializedEvents, SerializeEvent(event))
	}
	req.SendJSON(serializedEvents, http.StatusOK)
}

// EventCreateView creates the event
func EventCreateView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.GetStatusCode())
		return
	}

	var eventData EventData
	var event Event
	var err helios.Error
	err = req.DeserializeRequestData(&eventData)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}
	err = DeserializeEvent(eventData, &event)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}
	event.ID = 0
	UpsertEvent(user, &event)
	req.SendJSON(SerializeEvent(event), http.StatusCreated)
}

// QuestionListView send list of questions
func QuestionListView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.GetStatusCode())
		return
	}

	eventID, errParseEventID := req.GetURLParamUint("eventID")
	if errParseEventID != nil {
		req.SendJSON(errEventNotFound.GetMessage(), errEventNotFound.GetStatusCode())
		return
	}

	var questions []Question
	var err helios.Error
	questions, err = GetAllQuestionOfEventAndUser(user, eventID)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}

	serializedQuestions := make([]QuestionData, 0)
	for _, question := range questions {
		serializedQuestions = append(serializedQuestions, SerializeQuestion(question))
	}
	req.SendJSON(serializedQuestions, http.StatusOK)
}

// QuestionCreateView creates the question
func QuestionCreateView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.GetStatusCode())
		return
	}

	eventID, errParseEventID := req.GetURLParamUint("eventID")
	if errParseEventID != nil {
		req.SendJSON(errEventNotFound.GetMessage(), errEventNotFound.GetStatusCode())
		return
	}

	var questionData QuestionData
	var question Question
	var err helios.Error
	err = req.DeserializeRequestData(&questionData)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}

	// Ignoring error because currently this function wont return any error
	DeserializeQuestion(questionData, &question)

	question.ID = 0
	question.EventID = eventID
	err = UpsertQuestion(user, &question)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}

	req.SendJSON(SerializeQuestion(question), http.StatusCreated)
}

// QuestionDetailView send the question
func QuestionDetailView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.GetStatusCode())
		return
	}

	eventID, errParseEventID := req.GetURLParamUint("eventID")
	if errParseEventID != nil {
		req.SendJSON(errEventNotFound.GetMessage(), errEventNotFound.GetStatusCode())
		return
	}

	questionID, errParseQuestionID := req.GetURLParamUint("questionID")
	if errParseQuestionID != nil {
		req.SendJSON(errQuestionNotFound.GetMessage(), errQuestionNotFound.GetStatusCode())
		return
	}

	var question *Question
	var err helios.Error
	question, err = GetQuestionOfEventAndUser(user, eventID, questionID)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}
	var serializedQuestion QuestionData = SerializeQuestion(*question)
	req.SendJSON(serializedQuestion, http.StatusOK)
}

// QuestionDeleteView delete the question
func QuestionDeleteView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.GetStatusCode())
		return
	}

	eventID, errParseEventID := req.GetURLParamUint("eventID")
	if errParseEventID != nil {
		req.SendJSON(errEventNotFound.GetMessage(), errEventNotFound.GetStatusCode())
		return
	}

	questionID, errParseQuestionID := req.GetURLParamUint("questionID")
	if errParseQuestionID != nil {
		req.SendJSON(errQuestionNotFound.GetMessage(), errQuestionNotFound.GetStatusCode())
		return
	}

	var question *Question
	var err helios.Error
	question, err = DeleteQuestion(user, eventID, questionID)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}
	var serializedQuestion QuestionData = SerializeQuestion(*question)
	req.SendJSON(serializedQuestion, http.StatusOK)
}

// SubmissionCreateView create a submission of a question
func SubmissionCreateView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.GetStatusCode())
		return
	}

	eventID, errParseEventID := req.GetURLParamUint("eventID")
	if errParseEventID != nil {
		req.SendJSON(errEventNotFound.GetMessage(), errEventNotFound.GetStatusCode())
		return
	}

	questionID, errParseQuestionID := req.GetURLParamUint("questionID")
	if errParseQuestionID != nil {
		req.SendJSON(errQuestionNotFound.GetMessage(), errQuestionNotFound.GetStatusCode())
		return
	}

	var submitSubmissionRequest SubmitSubmissionRequest
	var errDeserialization helios.Error = req.DeserializeRequestData(&submitSubmissionRequest)
	if errDeserialization != nil {
		req.SendJSON(errDeserialization.GetMessage(), errDeserialization.GetStatusCode())
		return
	}

	var question *Question
	var err helios.Error
	question, err = SubmitSubmission(user, eventID, questionID, submitSubmissionRequest.Answer)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
	} else {
		var questionData QuestionData = SerializeQuestion(*question)
		req.SendJSON(questionData, http.StatusCreated)
	}
}
