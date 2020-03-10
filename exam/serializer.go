package exam

import "time"

// EventResponse is JSON representation of exam event.
type EventResponse struct {
	ID       uint      `json:"id"`
	Title    string    `json:"title"`
	StartsAt time.Time `json:"startsAt"`
	EndsAt   time.Time `json:"endsAt"`
}

// SubmitSubmissionRequest is JSON representation of request data
// when user wants to answer a question
type SubmitSubmissionRequest struct {
	Answer string `json:"answer"`
}

// QuestionResponse is JSON representation of question.
// Answer is the user's answer of the question, equals to Submission.Answer
type QuestionResponse struct {
	ID      uint     `json:"id"`
	Content string   `json:"content"`
	Choices []string `json:"choices"`
	Answer  string   `json:"answer"`
}

// SerializeEvent converts Event object event to JSON of event
func SerializeEvent(event Event) EventResponse {
	eventData := EventResponse{
		ID:       event.ID,
		Title:    event.Title,
		StartsAt: event.StartsAt,
		EndsAt:   event.EndsAt,
	}
	return eventData
}

// SerializeQuestion converts Question object question to JSON of question
func SerializeQuestion(question Question) QuestionResponse {
	choices := make([]string, 0)
	for _, choice := range question.Choices {
		choices = append(choices, choice.Text)
	}
	questionData := QuestionResponse{
		ID:      question.ID,
		Content: question.Content,
		Choices: choices,
		Answer:  question.UserAnswer,
	}

	return questionData
}

// DeserializeQuestion convert JSON of question to Question object
func DeserializeQuestion(questionData QuestionResponse) Question {
	var questionChoices []QuestionChoice
	for _, choiceText := range questionData.Choices {
		questionChoices = append(questionChoices, QuestionChoice{
			Text: choiceText,
		})
	}
	question := Question{
		ID:      questionData.ID,
		Content: questionData.Content,
		Choices: questionChoices,
	}
	return question
}
