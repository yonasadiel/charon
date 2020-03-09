package exam

// isAnswerValidChoice returns true if the answer is exist in one of question choices.
// Returns true if the choices are empty. It means that the question is not multiple choices type of question.
func isAnswerValidChoice(answer string, choices []QuestionChoice) bool {
	if len(choices) == 0 {
		return true
	}
	var exist bool = false
	for _, choice := range choices {
		if choice.Text == answer {
			exist = true
		}
	}
	return exist
}
