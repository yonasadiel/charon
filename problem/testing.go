package problem

import (
	"github.com/yonasadiel/charon/auth"
	"github.com/yonasadiel/helios"
)

var user1 auth.User
var user2 auth.User
var questionSimple Question
var questionWithChoice Question
var questionUnanswered Question
var questionUnowned Question
var submissionUser1QuestionSimple1 Submission
var submissionUser1QuestionSimple2 Submission
var submissionUser1QuestionWithChoice1 Submission
var submissionUser2QuestionSimple1 Submission

func beforeTest(populate bool) {
	helios.App.BeforeTest()

	if populate {
		user1 = auth.User{Email: "user1"}
		user2 = auth.User{Email: "user2"}
		helios.DB.Create(&user1)
		helios.DB.Create(&user2)

		questionSimple = Question{
			Content: "abc",
			Choices: []QuestionChoice{},
		}
		questionWithChoice = Question{
			Content: "def",
			Choices: []QuestionChoice{
				QuestionChoice{Text: "choice1"},
				QuestionChoice{Text: "choice2"},
			},
		}
		questionUnanswered = Question{
			Content: "ghi",
			Choices: []QuestionChoice{},
		}
		questionUnowned = Question{
			Content: "jkl",
			Choices: []QuestionChoice{},
		}
		helios.DB.Create(&questionSimple)
		helios.DB.Create(&questionWithChoice)
		helios.DB.Create(&questionUnanswered)
		helios.DB.Create(&questionUnowned)

		submissionUser1QuestionSimple1 = Submission{
			Answer:     "answer1",
			QuestionID: questionSimple.ID,
			Question:   &questionSimple,
			UserID:     user1.ID,
			User:       &user1,
		}
		submissionUser1QuestionSimple2 = Submission{
			Answer:     "answer2",
			QuestionID: questionSimple.ID,
			Question:   &questionSimple,
			UserID:     user1.ID,
			User:       &user1,
		}
		submissionUser1QuestionWithChoice1 = Submission{
			Answer:     questionWithChoice.Choices[0].Text,
			QuestionID: questionWithChoice.ID,
			Question:   &questionWithChoice,
			UserID:     user1.ID,
			User:       &user1,
		}
		submissionUser2QuestionSimple1 = Submission{
			Answer:     "answer3",
			QuestionID: questionSimple.ID,
			Question:   &questionSimple,
			UserID:     user2.ID,
			User:       &user2,
		}
		helios.DB.Create(&submissionUser1QuestionSimple1)
		helios.DB.Create(&submissionUser1QuestionSimple2)
		helios.DB.Create(&submissionUser1QuestionWithChoice1)
		helios.DB.Create(&submissionUser2QuestionSimple1)

		// Connect all user to all questions, except questionUnowned.
		helios.DB.Create(&UserQuestion{UserID: user1.ID, QuestionID: questionUnanswered.ID, Ordering: 3})
		helios.DB.Create(&UserQuestion{UserID: user1.ID, QuestionID: questionWithChoice.ID, Ordering: 2})
		helios.DB.Create(&UserQuestion{UserID: user1.ID, QuestionID: questionSimple.ID, Ordering: 1})
		helios.DB.Create(&UserQuestion{UserID: user2.ID, QuestionID: questionSimple.ID, Ordering: 1})
		helios.DB.Create(&UserQuestion{UserID: user2.ID, QuestionID: questionWithChoice.ID, Ordering: 2})
		helios.DB.Create(&UserQuestion{UserID: user2.ID, QuestionID: questionUnanswered.ID, Ordering: 3})
	}
}
