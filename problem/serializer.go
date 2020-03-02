package problem

// QuestionResponse is JSON representation of question.
// Answer is the user's answer of the question, equals to Submission.Answer
type QuestionResponse struct {
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
