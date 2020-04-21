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

// ParticipationListView send list of participations
func ParticipationListView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.GetStatusCode())
		return
	}

	var eventSlug string = req.GetURLParam("eventSlug")
	var participations []Participation
	var err helios.Error
	participations, err = GetAllParticipationOfUserAndEvent(user, eventSlug)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}

	serializedParticipations := make([]ParticipationData, 0)
	for _, participation := range participations {
		serializedParticipations = append(serializedParticipations, SerializeParticipation(participation))
	}
	req.SendJSON(serializedParticipations, http.StatusOK)
}

// ParticipationCreateView creates the participation
func ParticipationCreateView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.GetStatusCode())
		return
	}

	var eventSlug string = req.GetURLParam("eventSlug")
	var participationData ParticipationData
	var participation Participation
	var err helios.Error
	err = req.DeserializeRequestData(&participationData)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}

	err = DeserializeParticipation(participationData, &participation)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}

	err = UpsertParticipation(user, eventSlug, participationData.UserUsername, &participation)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}

	req.SendJSON(SerializeParticipation(participation), http.StatusOK)
}

// ParticipationDeleteView delete the participation
func ParticipationDeleteView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.GetStatusCode())
		return
	}

	var eventSlug string = req.GetURLParam("eventSlug")
	participationID, errParseParticipationID := req.GetURLParamUint("participationID")
	if errParseParticipationID != nil {
		req.SendJSON(errParticipationNotFound.GetMessage(), errParticipationNotFound.GetStatusCode())
		return
	}

	var participation *Participation
	var err helios.Error
	participation, err = DeleteParticipation(user, eventSlug, participationID)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}
	var serializedParticipation ParticipationData = SerializeParticipation(*participation)
	req.SendJSON(serializedParticipation, http.StatusOK)
}

// QuestionListView send list of questions
func QuestionListView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.GetStatusCode())
		return
	}

	var eventSlug string = req.GetURLParam("eventSlug")
	var questions []Question
	var err helios.Error
	questions, err = GetAllQuestionOfUserAndEvent(user, eventSlug)
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

	var eventSlug string = req.GetURLParam("eventSlug")
	var questionData QuestionData
	var question Question
	var err helios.Error
	err = req.DeserializeRequestData(&questionData)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}

	err = DeserializeQuestion(questionData, &question)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}

	question.ID = 0
	err = UpsertQuestion(user, eventSlug, &question)
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
	var eventSlug string = req.GetURLParam("eventSlug")
	questionID, errParseQuestionID := req.GetURLParamUint("questionID")
	if errParseQuestionID != nil {
		req.SendJSON(errQuestionNotFound.GetMessage(), errQuestionNotFound.GetStatusCode())
		return
	}

	var question *Question
	var err helios.Error
	question, err = GetQuestionOfEventAndUser(user, eventSlug, questionID)
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

	var eventSlug string = req.GetURLParam("eventSlug")
	questionID, errParseQuestionID := req.GetURLParamUint("questionID")
	if errParseQuestionID != nil {
		req.SendJSON(errQuestionNotFound.GetMessage(), errQuestionNotFound.GetStatusCode())
		return
	}

	var question *Question
	var err helios.Error
	question, err = DeleteQuestion(user, eventSlug, questionID)
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

	var eventSlug string = req.GetURLParam("eventSlug")
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
	question, err = SubmitSubmission(user, eventSlug, questionID, submitSubmissionRequest.Answer)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
	} else {
		var questionData QuestionData = SerializeQuestion(*question)
		req.SendJSON(questionData, http.StatusCreated)
	}
}

// GetSynchronizationDataView gets the synchronization data of event
func GetSynchronizationDataView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.GetStatusCode())
		return
	}
	var eventSlug string = req.GetURLParam("eventSlug")
	var event *Event
	var venue *Venue
	var questions []Question
	var users []auth.User
	var err helios.Error

	event, venue, questions, users, err = GetSynchronizationData(user, eventSlug)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
	} else {
		var synchronizationData SynchronizationData = SerializeSynchronizationData(*event, *venue, questions, users)
		req.SendJSON(synchronizationData, http.StatusOK)
	}
}

// PutSynchronizationDataView gets the synchronization data of event
func PutSynchronizationDataView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.GetStatusCode())
		return
	}

	var synchronizationData SynchronizationData
	var errDeserialization helios.Error = req.DeserializeRequestData(&synchronizationData)
	if errDeserialization != nil {
		req.SendJSON(errDeserialization.GetMessage(), errDeserialization.GetStatusCode())
		return
	}

	var event Event
	var venue Venue
	var questions []Question
	var users []auth.User
	var err helios.Error

	err = DeserializeSynchronizationData(synchronizationData, &event, &venue, &questions, &users)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}

	err = PutSynchronizationData(user, event, venue, questions, users)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
	} else {
		req.SendJSON("OK", http.StatusCreated)
	}
}
