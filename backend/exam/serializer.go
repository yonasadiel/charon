package exam

import (
	"time"

	"github.com/yonasadiel/helios"
)

// EventData is JSON representation of exam event.
type EventData struct {
	ID          uint   `json:"id"`
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
	var err helios.FormError
	var valid bool = true
	venue.ID = venueData.ID
	venue.Name = venueData.Name

	err = helios.FormError{}
	if venue.Name == "" {
		err.AddFieldError("name", "Name can't be empty")
		valid = false
	}
	if !valid {
		return err
	}
	return nil
}

// SerializeEvent converts Event object event to JSON of event
func SerializeEvent(event Event) EventData {
	eventData := EventData{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		StartsAt:    event.StartsAt.Local().Format(time.RFC3339),
		EndsAt:      event.EndsAt.Local().Format(time.RFC3339),
	}
	return eventData
}

// DeserializeEvent returns the Event from EventData
func DeserializeEvent(eventData EventData, event *Event) helios.Error {
	var err helios.FormError
	var errStartsAt, errEndsAt error
	var valid bool = true
	event.ID = eventData.ID
	event.Description = eventData.Description
	event.Title = eventData.Title
	event.StartsAt, errStartsAt = time.Parse(time.RFC3339, eventData.StartsAt)
	event.EndsAt, errEndsAt = time.Parse(time.RFC3339, eventData.EndsAt)

	err = helios.FormError{}
	if event.Title == "" {
		err.AddFieldError("title", "Title can't be empty")
		valid = false
	}
	if eventData.StartsAt == "" {
		err.AddFieldError("startsAt", "Start time must be provided")
		valid = false
	} else if errStartsAt != nil {
		err.AddFieldError("startsAt", "Failed to parse time")
		valid = false
	}
	if eventData.EndsAt == "" {
		err.AddFieldError("endsAt", "End time must be provided")
		valid = false
	} else if errEndsAt != nil {
		err.AddFieldError("endsAt", "Failed to parse time")
		valid = false
	}
	if event.EndsAt.Before(event.StartsAt) {
		err.AddFieldError("endsAt", "End time should be after start time")
		valid = false
	}
	if !valid {
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
	var err helios.FormError
	var valid bool = true
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

	err = helios.FormError{}
	if question.Content == "" {
		err.AddFieldError("content", "Content can't be empty")
		valid = false
	}
	if !valid {
		return err
	}
	return nil
}
