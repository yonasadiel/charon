package exam

import (
	"time"

	"github.com/yonasadiel/helios"
)

// EventData is JSON representation of exam event.
type EventData struct {
	ID       uint      `json:"id"`
	Title    string    `json:"title"`
	StartsAt time.Time `json:"startsAt"`
	EndsAt   time.Time `json:"endsAt"`
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

// SerializeEvent converts Event object event to JSON of event
func SerializeEvent(event Event) EventData {
	eventData := EventData{
		ID:       event.ID,
		Title:    event.Title,
		StartsAt: event.StartsAt,
		EndsAt:   event.EndsAt,
	}
	return eventData
}

// DeserializeEvent returns the Event from EventData
func DeserializeEvent(eventData EventData, event *Event) helios.Error {
	var err helios.FormError
	var valid bool = true
	event.ID = eventData.ID
	event.Title = eventData.Title
	event.StartsAt = eventData.StartsAt
	event.EndsAt = eventData.EndsAt

	err = helios.FormError{}
	if event.Title == "" {
		err.AddFieldError("title", "Title can't be empty")
		valid = false
	}
	if event.StartsAt.IsZero() {
		err.AddFieldError("startsAt", "Start time must be provided")
		valid = false
	}
	if event.EndsAt.IsZero() {
		err.AddFieldError("endsAt", "End time must be provided")
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
	var questionChoices []QuestionChoice
	for _, choiceText := range questionData.Choices {
		questionChoices = append(questionChoices, QuestionChoice{
			Text: choiceText,
		})
	}
	question.ID = questionData.ID
	question.Content = questionData.Content
	question.Choices = questionChoices
	return nil
}
