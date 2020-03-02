package problem

import (
	"github.com/yonasadiel/charon/auth"
	"github.com/yonasadiel/helios"
)

// GetAllQuestionOfUser returns all questions in database
// that belongs to an user, including the current submission
func GetAllQuestionOfUser(user auth.User) []Question {
	var questions []Question
	var userSubmissions []Submission
	var userSubmissionByQuestionID = make(map[uint]Submission)

	// Querying for user questions and user submissions
	helios.DB.
		Table("user_questions").
		Select("questions.*").
		Joins("left join questions on questions.id = user_questions.question_id").
		Where("user_id = ?", user.ID).
		Order("user_questions.ordering asc").
		Find(&questions)
	helios.DB.Where("user_id = ?", user.ID).Order("created_at asc").Find(&userSubmissions)

	// Mapping user submission by the question id
	for _, userSubmission := range userSubmissions {
		// Here, we safely assume that if the userSubmission is latter
		// in array, then it has latest creation time
		userSubmissionByQuestionID[userSubmission.QuestionID] = userSubmission
	}

	// Set the question answer to the user submission
	for i := range questions {
		userSubmission := userSubmissionByQuestionID[questions[i].ID]
		if userSubmission.ID != 0 {
			questions[i].UserAnswer = userSubmission.Answer
		}
	}

	return questions
}

// GetQuestionOfUser returns a question with given id, but first check
// if the user has rights to the question
func GetQuestionOfUser(questionID uint, user auth.User) *Question {
	var question Question
	var userSubmission Submission

	helios.DB.
		Table("user_questions").
		Select("questions.*").
		Joins("left join questions on questions.id = user_questions.question_id").
		Where("user_id = ?", user.ID).
		Where("question_id = ?", questionID).
		First(&question)
	helios.DB.Where("user_id = ?", user.ID).Where("question_id = ?", questionID).Order("created_at desc").First(&userSubmission)

	if question.ID == 0 {
		return nil
	}
	if userSubmission.ID != 0 {
		question.UserAnswer = userSubmission.Answer
	}
	return &question
}

// SubmitSubmission submit a submission from user to a question.
func SubmitSubmission(questionID uint, user auth.User, answer string) (*Submission, *helios.APIError) {
	var question Question
	var choices []QuestionChoice
	var submission Submission

	helios.DB.Where("id = ?", questionID).First(&question)
	helios.DB.Where("question_id = ?", questionID).Find(&choices)
	if question.ID == 0 {
		return nil, &errQuestionNotFound
	}
	if !isAnswerValidChoice(answer, choices) {
		return nil, &errAnswerNotValid
	}
	submission = Submission{
		Answer:     answer,
		QuestionID: question.ID,
		UserID:     user.ID,
		Question:   &question,
		User:       &user,
	}
	helios.DB.Create(&submission)
	return &submission, nil
}
