package exam

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yonasadiel/charon/auth"
	"github.com/yonasadiel/helios"
)

func TestGetAllEventOfUser(t *testing.T) {
	beforeTest(true)

	eventsParticipant := GetAllEventOfUser(userParticipant)
	assert.Equal(t, 2, len(eventsParticipant), "eventUnparticipated should not be included")
	assert.Equal(t, event1.Title, eventsParticipant[0].Title, "Event should be ordered by start time, so event1 will be the first")
	assert.Equal(t, event2.Title, eventsParticipant[1].Title, "Event should be ordered by start time, so event2 will be the second")

	eventsLocal := GetAllEventOfUser(userLocal)
	assert.Equal(t, 1, len(eventsLocal), "eventUnparticipated should not be shown to local user")

	eventsOrganizer := GetAllEventOfUser(userOrganizer)
	assert.Equal(t, 3, len(eventsOrganizer), "eventUnparticipated should be shown to organizer user")

	eventsAdmin := GetAllEventOfUser(userAdmin)
	assert.Equal(t, 3, len(eventsAdmin), "eventUnparticipated should be shown to admin user")
}

func TestUpsertEvent(t *testing.T) {
	beforeTest(false)
	var eventCount int
	var eventSaved Event

	err1 := UpsertEvent(user1, &event1)
	assert.Equal(t, errEventChangeNotAuthorized, err1, "participant user should not be able to upsert event")

	user1.SetAsLocal()
	err2 := UpsertEvent(user1, &event1)
	assert.Equal(t, errEventChangeNotAuthorized, err2, "local user should not be able to upsert event")

	user1.SetAsOrganizer()
	event1.ID = 0
	err3 := UpsertEvent(user1, &event1)
	helios.DB.Model(&Event{}).Count(&eventCount)
	assert.Nil(t, err3, "organizer user should be able to upsert event")
	assert.Equal(t, 1, eventCount, "Event should be created")

	user1.SetAsAdmin()
	event1.ID = 0
	err4 := UpsertEvent(user1, &event1)
	helios.DB.Model(&Event{}).Count(&eventCount)
	helios.DB.Where("id = ?", event1.ID).First(&eventSaved)
	assert.Nil(t, err4, "organizer user should be able to upsert event")
	assert.Equal(t, 2, eventCount, "Event should be created")
	assert.Equal(t, event1.Title, eventSaved.Title, "Event title is not equal")

	event1.Title = "ABC"
	err5 := UpsertEvent(user1, &event1)
	helios.DB.Model(&Event{}).Count(&eventCount)
	helios.DB.Where("id = ?", event1.ID).First(&eventSaved)
	assert.Nil(t, err5, "organizer user should be able to upsert event")
	assert.Equal(t, 2, eventCount, "Event should be updated, not created")
	assert.Equal(t, event1.Title, eventSaved.Title, "Event title is not equal")
}

func TestUpsertQuestion(t *testing.T) {
	beforeTest(true)
	var questionCountBefore int
	var choiceCountBefore int

	helios.DB.Model(&Question{}).Count(&questionCountBefore)
	helios.DB.Model(&QuestionChoice{}).Count(&choiceCountBefore)

	type questionUpsertTestCase struct {
		user          auth.User
		question      Question
		questionCount int
		choiceCount   int
		errExpected   helios.Error
	}
	testCases := []questionUpsertTestCase{
		questionUpsertTestCase{
			user:          userParticipant,
			question:      Question{Content: "Content 1", EventID: event1.ID},
			questionCount: questionCountBefore,
			choiceCount:   choiceCountBefore,
			errExpected:   errQuestionChangeNotAuthorized,
		},
		questionUpsertTestCase{
			user:          userLocal,
			question:      Question{Content: "Content 2", EventID: event1.ID},
			questionCount: questionCountBefore,
			choiceCount:   choiceCountBefore,
			errExpected:   errQuestionChangeNotAuthorized,
		},
		questionUpsertTestCase{
			user:          userOrganizer,
			question:      Question{Content: "Content 3", EventID: event1.ID},
			questionCount: questionCountBefore + 1,
			choiceCount:   choiceCountBefore,
			errExpected:   nil,
		},
		questionUpsertTestCase{
			user:          userAdmin,
			question:      Question{ID: questionSimple.ID, Content: "Content 4", EventID: event2.ID},
			questionCount: questionCountBefore + 1,
			choiceCount:   choiceCountBefore,
			errExpected:   nil,
		},
		questionUpsertTestCase{
			user: userAdmin,
			question: Question{
				Content: "Content 5",
				EventID: event1.ID,
				Choices: []QuestionChoice{
					QuestionChoice{ID: questionWithChoice.Choices[0].ID, Text: "Choice 5.1"}, // the ID will be ignored
					QuestionChoice{Text: "Choice 5.2"},
				},
			},
			questionCount: questionCountBefore + 2,
			choiceCount:   choiceCountBefore + 2,
			errExpected:   nil,
		},
		questionUpsertTestCase{
			user: userAdmin,
			question: Question{
				ID:      questionWithChoice.ID,
				Content: "Content 6",
				EventID: event1.ID,
				Choices: []QuestionChoice{
					QuestionChoice{ID: questionWithChoice.Choices[0].ID, Text: "Choice 6.1"},
					QuestionChoice{Text: "Choice 6.2"},
					QuestionChoice{Text: "Choice 6.3"},
					QuestionChoice{Text: "Choice 6.4"},
				},
			},
			questionCount: questionCountBefore + 2,
			choiceCount:   choiceCountBefore + 2 - len(questionWithChoice.Choices) + 4,
			errExpected:   nil,
		},
	}

	for i, testCase := range testCases {
		var questionCount int
		var choiceCount int
		var questionSaved Question
		t.Logf("Test UpsertQuestion testcase: %d", i)
		err := UpsertQuestion(testCase.user, &testCase.question)
		helios.DB.Model(&Question{}).Count(&questionCount)
		helios.DB.Model(&QuestionChoice{}).Count(&choiceCount)
		helios.DB.Where("id = ?", testCase.question.ID).First(&questionSaved)
		if testCase.errExpected == nil {
			assert.Nil(t, err, "There should be no error")
			assert.Equal(t, testCase.questionCount, questionCount, "Different number of questions expected")
			assert.Equal(t, testCase.choiceCount, choiceCount, "Different number of question choices expected")
			assert.Equal(t, testCase.question.Content, questionSaved.Content, "Different question content")
			assert.Equal(t, testCase.question.EventID, questionSaved.EventID, "Different question event id")
		} else {
			assert.Equal(t, testCase.errExpected, err, "Different error expected")
			assert.Equal(t, testCase.questionCount, questionCount, "Different number of questions expected")
			assert.Equal(t, testCase.choiceCount, choiceCount, "Different number of question choices expected")
		}
	}
}

func TestGetAllQuestionOfEventAndUser(t *testing.T) {
	beforeTest(true)

	var questions []Question
	var err helios.Error

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

	question1, err1 := GetQuestionOfUser(user1, event1.ID, questionSimple.ID)
	assert.Nil(t, err1, "Failed to get question")
	assert.NotNil(t, question1, "Question is not found")
	assert.Equal(t, questionSimple.ID, question1.ID, "Different question content")
	assert.Equal(t, questionSimple.Content, question1.Content, "Different question content")
	assert.Equal(t, submissionUser1QuestionSimple2.Answer, question1.UserAnswer, "The answer should be latest submission")

	_, err2 := GetQuestionOfUser(user1, event1.ID, 4567)
	assert.Equal(t, errQuestionNotFound, err2, "Unknwon question id returns errQuestionNotFound")

	_, err3 := GetQuestionOfUser(user1, event1.ID, questionUnowned.ID)
	assert.Equal(t, errQuestionNotFound, err3, "Question unowned by the user should not be found")

	question4, err4 := GetQuestionOfUser(user1, event1.ID, questionUnanswered.ID)
	assert.Nil(t, err4, "Failed to get question")
	assert.NotNil(t, question4, "Question is not found")
	assert.Equal(t, questionUnanswered.ID, question4.ID, "Different question content")
	assert.Equal(t, questionUnanswered.Content, question4.Content, "Different question content")
	assert.Equal(t, "", question4.UserAnswer, "The answer should be latest submission")

}

func TestSubmitSubmissionSuccess(t *testing.T) {
	beforeTest(true)

	submission1Answer := "answer2"
	submission1Returned, err := SubmitSubmission(user1, event1.ID, questionSimple.ID, submission1Answer)
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
	submission2Returned, err := SubmitSubmission(user2, event1.ID, questionWithChoice.ID, submission2Answer)
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

func TestSubmitInvalidQuestionIDOrEventID(t *testing.T) {
	beforeTest(true)

	var submissionCountBefore, submissionCountAfter int
	helios.DB.Model(&Submission{}).Count(&submissionCountBefore)
	submission1Answer := "answer1"
	submission1Returned, err := SubmitSubmission(user1, event1.ID, 30, submission1Answer)
	helios.DB.Model(&Submission{}).Count(&submissionCountAfter)
	assert.Equal(t, errQuestionNotFound, err, "Submission should be fail")
	assert.Nil(t, submission1Returned, "Fail to submit should return nil submission")
	assert.Equal(t, submissionCountBefore, submissionCountAfter, "Submission should not be made")
}

func TestSubmitSubmissionInvalidChoice(t *testing.T) {
	beforeTest(true)

	var submissionCountBefore, submissionCountAfter int
	helios.DB.Model(&Submission{}).Count(&submissionCountBefore)
	submission1Answer := "not in choice"
	submission1Returned, err := SubmitSubmission(user1, event1.ID, questionWithChoice.ID, submission1Answer)
	helios.DB.Model(&Submission{}).Count(&submissionCountAfter)
	assert.Equal(t, errAnswerNotValid, err, "Submission should be fail")
	assert.Nil(t, submission1Returned, "Fail to submit should return nil submission")
	assert.Equal(t, submissionCountBefore, submissionCountAfter, "Submission should not be made")
}
