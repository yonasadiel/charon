package exam

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yonasadiel/helios"
)

func TestGetAllEventOfUser(t *testing.T) {
	beforeTest(true)

	eventsParticipant := GetAllEventOfUser(user1)
	assert.Equal(t, 2, len(eventsParticipant), "eventUnparticipated should not be included")
	assert.Equal(t, event1.Title, eventsParticipant[0].Title, "Event should be ordered by start time, so event1 will be the first")
	assert.Equal(t, event2.Title, eventsParticipant[1].Title, "Event should be ordered by start time, so event2 will be the second")

	eventsLocal := GetAllEventOfUser(userLocal)
	assert.Equal(t, 3, len(eventsLocal), "eventUnparticipated should be shown because the user has local role")
}

func TestGetAllQuestionOfEventAndUser(t *testing.T) {
	beforeTest(true)

	var questions []Question
	var err *helios.APIError

	questions, err = GetAllQuestionOfEventAndUser(user1, event1.ID)
	assert.Nil(t, err, "Failed to get all question")
	assert.Equal(t, 3, len(questions), "Different number of questions. Maybe questionEvent2 or questionUnowned is included?")
	assert.Equal(t, questionSimple.Content, questions[0].Content, "Different question content")
	assert.Equal(t, submissionUser1QuestionSimple2.Answer, questions[0].UserAnswer, "The answer should be latest submission")
	assert.Equal(t, questionWithChoice.Content, questions[1].Content, "Different question content")
	assert.Equal(t, submissionUser1QuestionWithChoice1.Answer, questions[1].UserAnswer, "Different answer on question with choice")
	assert.Empty(t, questions[2].UserAnswer, "Quesion unanswered should be unanswered")

	questions, err = GetAllQuestionOfEventAndUser(user1, 1234)
	assert.NotNil(t, err, "Random event ID will not found")
}

func TestGetQuestionOfUser(t *testing.T) {
	beforeTest(true)

	question1 := GetQuestionOfUser(user1, questionSimple.ID)
	assert.NotNil(t, question1, "Question is not found")
	assert.Equal(t, questionSimple.ID, question1.ID, "Different question content")
	assert.Equal(t, questionSimple.Content, question1.Content, "Different question content")
	assert.Equal(t, submissionUser1QuestionSimple2.Answer, question1.UserAnswer, "The answer should be latest submission")

	question2 := GetQuestionOfUser(user1, 4567)
	assert.Nil(t, question2, "Question invalid ID should not be found")

	question3 := GetQuestionOfUser(user1, questionUnowned.ID)
	assert.Nil(t, question3, "Question unowned by the user should not be found")

	question4 := GetQuestionOfUser(user1, questionUnanswered.ID)
	assert.NotNil(t, question4, "Question is not found")
	assert.Equal(t, questionUnanswered.ID, question4.ID, "Different question content")
	assert.Equal(t, questionUnanswered.Content, question4.Content, "Different question content")
	assert.Equal(t, "", question4.UserAnswer, "The answer should be latest submission")

}

func TestSubmitSubmissionSuccess(t *testing.T) {
	beforeTest(true)

	submission1Answer := "answer2"
	submission1Returned, err := SubmitSubmission(user1, questionSimple.ID, submission1Answer)
	assert.Nil(t, err, "Failed to submit submission")
	assert.Equal(t, submission1Answer, submission1Returned.Answer, "Answer returned different with answer submitted")

	var submission1Stored Submission
	helios.DB.Where("question_id = ?", questionSimple.ID).Where("user_id = ?", user1.ID).Order("created_at desc").First(&submission1Stored)
	assert.NotEqual(t, 0, submission1Stored.ID, "Submission is not stored to database")
	assert.Equal(t, submission1Returned.ID, submission1Stored.ID, "Different submission ID returned with stored")
	assert.Equal(t, submission1Answer, submission1Stored.Answer, "Different answer stored in database")
	assert.Equal(t, questionSimple.ID, submission1Stored.QuestionID, "Different question ID stored in database")
	assert.Equal(t, user1.ID, submission1Stored.UserID, "Different user ID stored in database")

	submission2Answer := "choice2"
	submission2Returned, err := SubmitSubmission(user2, questionWithChoice.ID, submission2Answer)
	assert.Nil(t, err, "Failed to submit submission")
	assert.Equal(t, submission2Answer, submission2Returned.Answer, "Answer returned different with answer submitted")

	var submission2Stored Submission
	helios.DB.Where("question_id = ?", questionWithChoice.ID).Where("user_id = ?", user2.ID).Order("created_at desc").First(&submission2Stored)
	assert.NotEqual(t, 0, submission2Stored.ID, "Submission is not stored to database")
	assert.Equal(t, submission2Returned.ID, submission2Stored.ID, "Different submission ID returned with stored")
	assert.Equal(t, submission2Answer, submission2Stored.Answer, "Different answer stored in database")
	assert.Equal(t, questionWithChoice.ID, submission2Stored.QuestionID, "Different question ID stored in database")
	assert.Equal(t, user2.ID, submission2Stored.UserID, "Different user ID stored in database")
}

func TestSubmitInvalidQuestionID(t *testing.T) {
	beforeTest(true)

	var submissionCountBefore, submissionCountAfter int
	helios.DB.Model(&Submission{}).Count(&submissionCountBefore)
	submission1Answer := "answer1"
	submission1Returned, err := SubmitSubmission(user1, 30, submission1Answer)
	helios.DB.Model(&Submission{}).Count(&submissionCountAfter)
	assert.Equal(t, errQuestionNotFound.Code, err.Code, "Submission should be fail")
	assert.Nil(t, submission1Returned, "Fail to submit should return nil submission")
	assert.Equal(t, submissionCountBefore, submissionCountAfter, "Submission should not be made")
}

func TestSubmitSubmissionInvalidChoice(t *testing.T) {
	beforeTest(true)

	var submissionCountBefore, submissionCountAfter int
	helios.DB.Model(&Submission{}).Count(&submissionCountBefore)
	submission1Answer := "not in choice"
	submission1Returned, err := SubmitSubmission(user1, questionWithChoice.ID, submission1Answer)
	helios.DB.Model(&Submission{}).Count(&submissionCountAfter)
	assert.Equal(t, errAnswerNotValid.Code, err.Code, "Submission should be fail")
	assert.Nil(t, submission1Returned, "Fail to submit should return nil submission")
	assert.Equal(t, submissionCountBefore, submissionCountAfter, "Submission should not be made")
}
