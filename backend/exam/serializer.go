package exam

import (
	"time"

	"github.com/yonasadiel/charon/backend/auth"
	"github.com/yonasadiel/helios"
)

// EventData is JSON representation of exam event.
type EventData struct {
	ID          uint   `json:"id"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartsAt    string `json:"startsAt"`
	EndsAt      string `json:"endsAt"`
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
}

// QuestionData is JSON representation of question.
// Answer is the user's answer of the question, equals to Submission.Answer
type QuestionData struct {
	ID      uint     `json:"id"`
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
	Event          EventData           `json:"event"`
	Questions      []QuestionData      `json:"questions"`
	Participations []ParticipationData `json:"participations"`
	Users          []auth.UserData     `json:"users"`
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
	eventData := EventData{
		ID:          event.ID,
		Slug:        event.Slug,
		Title:       event.Title,
		Description: event.Description,
		StartsAt:    event.StartsAt.Local().Format(time.RFC3339),
		EndsAt:      event.EndsAt.Local().Format(time.RFC3339),
	}
	return eventData
}

// DeserializeEvent returns the Event from EventData
func DeserializeEvent(eventData EventData, event *Event) helios.Error {
	var err helios.ErrorForm = helios.NewErrorForm()
	var errStartsAt, errEndsAt error
	event.ID = eventData.ID
	event.Slug = eventData.Slug
	event.Description = eventData.Description
	event.Title = eventData.Title
	event.StartsAt, errStartsAt = time.Parse(time.RFC3339, eventData.StartsAt)
	event.EndsAt, errEndsAt = time.Parse(time.RFC3339, eventData.EndsAt)

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

// SerializeQuestion converts Question object question to JSON of question
func SerializeQuestion(question Question) QuestionData {
	choices := make([]string, 0)
	for _, choice := range question.Choices {
		choices = append(choices, choice.Text)
	}
	questionData := QuestionData{
		ID:      question.ID,
		Content: question.Content,
		Choices: choices,
		Answer:  question.UserAnswer,
	}

	return questionData
}

// DeserializeQuestion convert JSON of question to Question object
func DeserializeQuestion(questionData QuestionData, question *Question) helios.Error {
	var err helios.ErrorForm = helios.NewErrorForm()
	var questionChoices []QuestionChoice
	for _, choiceText := range questionData.Choices {
		if choiceText != "" {
			questionChoices = append(questionChoices, QuestionChoice{
				Text: choiceText,
			})
		}
	}
	question.ID = questionData.ID
	question.Content = questionData.Content
	question.Choices = questionChoices

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
func SerializeSynchronizationData(event Event, questions []Question, participations []Participation, users []auth.User) SynchronizationData {
	var questionsData []QuestionData = make([]QuestionData, 0)
	var participationsData []ParticipationData = make([]ParticipationData, 0)
	var usersData []auth.UserData = make([]auth.UserData, 0)
	for _, question := range questions {
		questionsData = append(questionsData, SerializeQuestion(question))
	}
	for _, participation := range participations {
		participationsData = append(participationsData, SerializeParticipation(participation))
	}
	for _, user := range users {
		usersData = append(usersData, auth.SerializeUser(user))
	}
	return SynchronizationData{
		Event:          SerializeEvent(event),
		Questions:      questionsData,
		Participations: participationsData,
		Users:          usersData,
	}
}

// DeserializeSynchronizationData converts event, questions, participations, and users
// into SynchronizationData
func DeserializeSynchronizationData(synchronizationData SynchronizationData, event *Event, questions []Question, participations []Participation, users []auth.User) helios.Error {
	var err helios.ErrorForm = helios.NewErrorForm()
	var errEvent helios.Error = DeserializeEvent(synchronizationData.Event, event)
	if errEvent != nil {
		var errEventForm helios.ErrorForm = errEvent.(helios.ErrorForm)
		err.FieldError["event"] = errEventForm.FieldError
		err.NonFieldError = errEventForm.NonFieldError
	}

	var errQuestions helios.ErrorFormFieldArray = make(helios.ErrorFormFieldArray, 0)
	for _, questionData := range synchronizationData.Questions {
		var question Question
		var errQuestion helios.Error = DeserializeQuestion(questionData, &question)
		if errQuestion == nil {
			questions = append(questions, question)
			errQuestions = append(errQuestions, helios.ErrorFormFieldNested{})
		} else {
			var errorQuestionForm helios.ErrorForm = errQuestion.(helios.ErrorForm)
			errQuestions = append(errQuestions, errorQuestionForm.FieldError)
			// Currently this is commented out because the deserialization doesn't have any non field error
			// for _, nonFieldError := range errorQuestionForm.NonFieldError {
			// 	err.NonFieldError = append(err.NonFieldError, nonFieldError)
			// }
		}
	}
	err.FieldError["questions"] = errQuestions

	var errParticipations helios.ErrorFormFieldArray = make(helios.ErrorFormFieldArray, 0)
	for _, participationData := range synchronizationData.Participations {
		var participation Participation
		var errParticipation helios.Error = DeserializeParticipation(participationData, &participation)
		if errParticipation == nil {
			participations = append(participations, participation)
			errParticipations = append(errParticipations, helios.ErrorFormFieldNested{})
		} else {
			var errParticipationForm helios.ErrorForm = errParticipation.(helios.ErrorForm)
			errParticipations = append(errParticipations, errParticipationForm.FieldError)
			// Currently this is commented out because the deserialization doesn't have any non field error
			// for _, nonFieldError := range errParticipationForm.NonFieldError {
			// 	err.NonFieldError = append(err.NonFieldError, nonFieldError)
			// }
		}
	}
	err.FieldError["participations"] = errParticipations

	var errUsers helios.ErrorFormFieldArray = make(helios.ErrorFormFieldArray, 0)
	for _, userData := range synchronizationData.Users {
		var user auth.User
		var errUser helios.Error = auth.DeserializeUser(userData, &user)
		if errUser == nil {
			users = append(users, user)
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

	if err.IsError() {
		return err
	}
	return nil
}
