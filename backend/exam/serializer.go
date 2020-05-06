package exam

import (
	"strings"
	"time"

	"github.com/yonasadiel/charon/backend/auth"
	"github.com/yonasadiel/helios"
)

// EventData is JSON representation of exam event.
type EventData struct {
	ID                  uint   `json:"id"`
	Slug                string `json:"slug"`
	Title               string `json:"title"`
	Description         string `json:"description"`
	StartsAt            string `json:"startsAt"`
	EndsAt              string `json:"endsAt"`
	SimKey              string `json:"simKey"`
	SimKeySign          string `json:"simKeySign"`
	PubKey              string `json:"pubKey"`
	IsDecrypted         bool   `json:"isDecrypted"`
	LastSynchronization string `json:"lastSynchronization"`
}

// VenueData is JSON representation of venue.
type VenueData struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ParticipationData is JSON representation of participation.
type ParticipationData struct {
	ID           uint   `json:"id"`
	UserUsername string `json:"userUsername"`
	VenueID      uint   `json:"venueId"`
	KeyPlain     string `json:"key,omitempty"`
	KeyTwice     string `json:"keyTwice"`
}

// VerificationData used for client submitting hashed once participation key
type VerificationData struct {
	KeyHashedOnce string `json:"key"`
}

// QuestionData is JSON representation of question.
// Answer is the user's answer of the question, equals to Submission.Answer
type QuestionData struct {
	Number  uint     `json:"number"`
	Content string   `json:"content"`
	Choices []string `json:"choices"`
	Answer  string   `json:"answer"`
}

// SubmitSubmissionRequest is JSON representation of request data
// when user wants to answer a question
type SubmitSubmissionRequest struct {
	Answer string `json:"answer"`
}

// SynchronizationData is JSON representation of encrypted data when
// event data passed before exam starts
type SynchronizationData struct {
	Event     EventData                   `json:"event"`
	Venue     VenueData                   `json:"venue"`
	Questions []QuestionData              `json:"questions"`
	Users     []auth.UserWithPasswordData `json:"users"`
	UsersKey  map[string]string           `json:"usersKey"`
}

// DecryptRequest is JSON representation of submitting key for
// decrypting event data
type DecryptRequest struct {
	SimKey string `json:"key"`
}

// SerializeVenue converts Venue object venue to JSON of venue
func SerializeVenue(venue Venue) VenueData {
	venueData := VenueData{
		ID:   venue.ID,
		Name: venue.Name,
	}
	return venueData
}

// DeserializeVenue returns the Venue from VenueData
func DeserializeVenue(venueData VenueData, venue *Venue) helios.Error {
	var err helios.ErrorForm = helios.NewErrorForm()
	venue.ID = venueData.ID
	venue.Name = venueData.Name

	if venue.Name == "" {
		err.FieldError["name"] = helios.ErrorFormFieldAtomic{"Name can't be empty"}
	}
	if err.IsError() {
		return err
	}
	return nil
}

// SerializeEvent converts Event object event to JSON of event
func SerializeEvent(event Event) EventData {
	var lastSynchronization string
	if !event.LastSynchronization.IsZero() {
		lastSynchronization = event.LastSynchronization.Local().Format(time.RFC3339)
	}
	eventData := EventData{
		ID:                  event.ID,
		Slug:                event.Slug,
		Title:               event.Title,
		Description:         event.Description,
		StartsAt:            event.StartsAt.Local().Format(time.RFC3339),
		EndsAt:              event.EndsAt.Local().Format(time.RFC3339),
		SimKey:              event.SimKey,
		SimKeySign:          event.SimKeySign,
		PubKey:              event.PubKey,
		IsDecrypted:         !event.DecryptedAt.IsZero(),
		LastSynchronization: lastSynchronization,
	}
	return eventData
}

// DeserializeEvent returns the Event from EventData
func DeserializeEvent(eventData EventData, event *Event) helios.Error {
	var err helios.ErrorForm = helios.NewErrorForm()
	var errStartsAt, errEndsAt, errLastSynchronization error
	event.ID = eventData.ID
	event.Slug = eventData.Slug
	event.Description = eventData.Description
	event.Title = eventData.Title
	event.SimKeySign = eventData.SimKeySign
	event.PubKey = eventData.PubKey
	event.StartsAt, errStartsAt = time.Parse(time.RFC3339, eventData.StartsAt)
	event.EndsAt, errEndsAt = time.Parse(time.RFC3339, eventData.EndsAt)
	event.LastSynchronization, errLastSynchronization = time.Parse(time.RFC3339, eventData.LastSynchronization)

	if event.Title == "" {
		err.FieldError["title"] = helios.ErrorFormFieldAtomic{"Title can't be empty"}
	}
	if event.Slug == "" {
		err.FieldError["slug"] = helios.ErrorFormFieldAtomic{"Slug can't be empty"}
	}
	if eventData.StartsAt == "" {
		err.FieldError["startsAt"] = helios.ErrorFormFieldAtomic{"Start time must be provided"}
	} else if errStartsAt != nil {
		err.FieldError["startsAt"] = helios.ErrorFormFieldAtomic{"Failed to parse time"}
	}
	if eventData.EndsAt == "" {
		err.FieldError["endsAt"] = helios.ErrorFormFieldAtomic{"End time must be provided"}
	} else if errEndsAt != nil {
		err.FieldError["endsAt"] = helios.ErrorFormFieldAtomic{"Failed to parse time"}
	}
	if event.EndsAt.Before(event.StartsAt) {
		err.FieldError["endsAt"] = helios.ErrorFormFieldAtomic{"End time should be after start time"}
	}
	if eventData.LastSynchronization == "" {
		event.LastSynchronization = time.Time{}
	} else if errLastSynchronization != nil {
		err.FieldError["lastSynchronization"] = helios.ErrorFormFieldAtomic{"Failed to parse time"}
	}
	if err.IsError() {
		return err
	}
	return nil
}

// SerializeParticipation converts Participation object participation to JSON of participation
func SerializeParticipation(participation Participation) ParticipationData {
	participationData := ParticipationData{
		ID:           participation.ID,
		UserUsername: participation.User.Username,
		VenueID:      participation.Venue.ID,
		KeyTwice:     participation.KeyHashedTwice,
	}
	return participationData
}

// DeserializeParticipation convert JSON of participation to Participation object
func DeserializeParticipation(participationData ParticipationData, participation *Participation) helios.Error {
	var err helios.ErrorForm = helios.NewErrorForm()
	participation.ID = participationData.ID
	participation.VenueID = participationData.VenueID

	if participation.VenueID == 0 {
		err.FieldError["venueId"] = helios.ErrorFormFieldAtomic{"Venue can't be empty"}
	}
	if participationData.UserUsername == "" {
		err.FieldError["userUsername"] = helios.ErrorFormFieldAtomic{"Username can't be empty"}
	}

	if err.IsError() {
		return err
	}
	return nil
}

// DeserializeParticipationWithKey convert JSON like DeserializeParticipation but with key
// used in creating participation
func DeserializeParticipationWithKey(participationData ParticipationData, participation *Participation) helios.Error {
	var err = DeserializeParticipation(participationData, participation)
	var errForm helios.ErrorForm = helios.NewErrorForm()
	if err != nil {
		errForm, _ = err.(helios.ErrorForm)
	}
	participation.KeyPlain = participationData.KeyPlain
	if len(participation.KeyPlain) != 32 {
		errForm.FieldError["key"] = helios.ErrorFormFieldAtomic{"Key length must be 32 chars"}
	}
	if errForm.IsError() {
		return errForm
	}
	return nil
}

// SerializeQuestion converts Question object question to JSON of question
func SerializeQuestion(question Question) QuestionData {
	var choicesArr []string = strings.Split(question.Choices, "|")
	var choices []string = make([]string, 0)
	for _, choice := range choicesArr {
		if len(choice) > 0 {
			choices = append(choices, choice)
		}
	}
	questionData := QuestionData{
		Number:  question.ID,
		Content: question.Content,
		Choices: choices,
		Answer:  question.UserAnswer,
	}

	return questionData
}

// DeserializeQuestion convert JSON of question to Question object
func DeserializeQuestion(questionData QuestionData, question *Question) helios.Error {
	var err helios.ErrorForm = helios.NewErrorForm()
	question.ID = questionData.Number
	question.Content = questionData.Content
	var choices []string = make([]string, 0)
	for _, choice := range questionData.Choices {
		if len(choice) > 0 {
			choices = append(choices, choice)
		}
	}
	question.Choices = strings.Join(choices, "|")

	if question.Content == "" {
		err.FieldError["content"] = helios.ErrorFormFieldAtomic{"Content can't be empty"}
	}
	if err.IsError() {
		return err
	}
	return nil
}

// SerializeSynchronizationData converts event, questions, participations, and users
// into SynchronizationData
func SerializeSynchronizationData(event Event, venue Venue, questions []Question, users []auth.User, usersKey map[string]string) SynchronizationData {
	var questionsData []QuestionData = make([]QuestionData, 0)
	var usersData []auth.UserWithPasswordData = make([]auth.UserWithPasswordData, 0)
	for _, question := range questions {
		questionsData = append(questionsData, SerializeQuestion(question))
	}
	for _, user := range users {
		usersData = append(usersData, auth.SerializeUserWithPassword(user))
	}
	return SynchronizationData{
		Event:     SerializeEvent(event),
		Venue:     SerializeVenue(venue),
		Questions: questionsData,
		Users:     usersData,
		UsersKey:  usersKey,
	}
}

// DeserializeSynchronizationData converts event, questions, participations, and users
// into SynchronizationData
func DeserializeSynchronizationData(synchronizationData SynchronizationData, event *Event, venue *Venue, questions *[]Question, users *[]auth.User, usersKey *map[string]string) helios.Error {
	var err helios.ErrorForm = helios.NewErrorForm()
	var errEvent helios.Error = DeserializeEvent(synchronizationData.Event, event)
	if errEvent != nil {
		var errEventForm helios.ErrorForm = errEvent.(helios.ErrorForm)
		err.FieldError["event"] = errEventForm.FieldError
		err.NonFieldError = errEventForm.NonFieldError
	}

	var errVenue helios.Error = DeserializeVenue(synchronizationData.Venue, venue)
	if errVenue != nil {
		var errVenueForm helios.ErrorForm = errVenue.(helios.ErrorForm)
		err.FieldError["venue"] = errVenueForm.FieldError
		// Currently this is commented out because the deserialization doesn't have any non field error
		// for _, nonFieldError := range errVenueForm.NonFieldError {
		// 	err.NonFieldError = append(err.NonFieldError, nonFieldError)
		// }
	}

	var errQuestions helios.ErrorFormFieldArray = make(helios.ErrorFormFieldArray, 0)
	for _, questionData := range synchronizationData.Questions {
		var question Question
		var errQuestion helios.Error = DeserializeQuestion(questionData, &question)
		if errQuestion == nil {
			*questions = append(*questions, question)
			errQuestions = append(errQuestions, helios.ErrorFormFieldNested{})
		} else {
			var errQuestionForm helios.ErrorForm = errQuestion.(helios.ErrorForm)
			errQuestions = append(errQuestions, errQuestionForm.FieldError)
			// Currently this is commented out because the deserialization doesn't have any non field error
			// for _, nonFieldError := range errQuestionForm.NonFieldError {
			// 	err.NonFieldError = append(err.NonFieldError, nonFieldError)
			// }
		}
	}
	err.FieldError["questions"] = errQuestions

	var errUsers helios.ErrorFormFieldArray = make(helios.ErrorFormFieldArray, 0)
	for _, userData := range synchronizationData.Users {
		var user auth.User
		var errUser helios.Error = auth.DeserializeUserWithPassword(userData, &user)
		if errUser == nil {
			*users = append(*users, user)
			errUsers = append(errUsers, helios.ErrorFormFieldNested{})
		} else {
			var errUserForm helios.ErrorForm = errUser.(helios.ErrorForm)
			errUsers = append(errUsers, errUserForm.FieldError)
			// Currently this is commented out because the deserialization doesn't have any non field error
			// for _, nonFieldError := range errUserForm.NonFieldError {
			// 	err.NonFieldError = append(err.NonFieldError, nonFieldError)
			// }
		}
	}
	err.FieldError["users"] = errUsers

	*usersKey = make(map[string]string)
	for k, v := range synchronizationData.UsersKey {
		(*usersKey)[k] = v
	}

	if err.IsError() {
		return err
	}
	return nil
}
