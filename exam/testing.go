package exam

import (
	"time"

	"github.com/yonasadiel/charon/auth"
	"github.com/yonasadiel/helios"
)

var event1 Event
var event2 Event
var eventUnparticipated Event
var user1 auth.User
var user2 auth.User
var userLocal auth.User
var userParticipant auth.User
var userOrganizer auth.User
var userAdmin auth.User
var questionSimple Question
var questionWithChoice Question
var questionUnanswered Question
var questionUnowned Question
var questionEvent2 Question
var questionEventUnparticipated Question
var submissionUser1QuestionSimple1 Submission
var submissionUser1QuestionSimple2 Submission
var submissionUser1QuestionWithChoice1 Submission
var submissionUser2QuestionSimple1 Submission

func beforeTest(populate bool) {
	helios.App.BeforeTest()

	if populate {
		utcTZ := time.FixedZone("UTC", 0)
		event1 = Event{
			Title:    "Event #1",
			StartsAt: time.Date(2020, 8, 12, 9, 30, 10, 0, utcTZ),
			EndsAt:   time.Date(2020, 8, 12, 4, 30, 10, 0, utcTZ),
		}
		event2 = Event{
			Title:    "Event #2",
			StartsAt: time.Date(2020, 8, 13, 9, 30, 10, 0, utcTZ),
			EndsAt:   time.Date(2020, 8, 13, 4, 30, 10, 0, utcTZ),
		}
		eventUnparticipated = Event{
			Title:    "Event #3",
			StartsAt: time.Date(2020, 8, 14, 9, 30, 10, 0, utcTZ),
			EndsAt:   time.Date(2020, 8, 14, 4, 30, 10, 0, utcTZ),
		}
		helios.DB.Create(&eventUnparticipated)
		helios.DB.Create(&event2)
		helios.DB.Create(&event1)

		user1 = auth.User{Username: "user1"}
		user2 = auth.User{Username: "user2"}
		helios.DB.Create(&user1)
		helios.DB.Create(&user2)

		userLocal = auth.User{Username: "userLocal"}
		userParticipant = auth.User{Username: "userParticipant"}
		userOrganizer = auth.User{Username: "userOrganizer"}
		userAdmin = auth.User{Username: "userAdmin"}
		userLocal.SetAsLocal()
		userParticipant.SetAsParticipant()
		userOrganizer.SetAsOrganizer()
		userAdmin.SetAsAdmin()
		helios.DB.Create(&userLocal)
		helios.DB.Create(&userParticipant)
		helios.DB.Create(&userOrganizer)
		helios.DB.Create(&userAdmin)

		// Connect all user to all events, except eventUnparticipated.
		helios.DB.Create(&UserEvent{UserID: user1.ID, EventID: event1.ID})
		helios.DB.Create(&UserEvent{UserID: user1.ID, EventID: event2.ID})
		helios.DB.Create(&UserEvent{UserID: user2.ID, EventID: event1.ID})
		helios.DB.Create(&UserEvent{UserID: user2.ID, EventID: event2.ID})
		helios.DB.Create(&UserEvent{UserID: userParticipant.ID, EventID: event1.ID})
		helios.DB.Create(&UserEvent{UserID: userParticipant.ID, EventID: event2.ID})
		helios.DB.Create(&UserEvent{UserID: userLocal.ID, EventID: event1.ID})

		questionSimple = Question{
			EventID: event1.ID,
			Content: "abc",
			Choices: []QuestionChoice{},
		}
		questionWithChoice = Question{
			EventID: event1.ID,
			Content: "def",
			Choices: []QuestionChoice{
				QuestionChoice{Text: "choice1"},
				QuestionChoice{Text: "choice2"},
			},
		}
		questionUnanswered = Question{
			EventID: event1.ID,
			Content: "ghi",
			Choices: []QuestionChoice{},
		}
		questionUnowned = Question{
			EventID: event1.ID,
			Content: "jkl",
			Choices: []QuestionChoice{},
		}
		questionEvent2 = Question{
			EventID: event2.ID,
			Content: "mno",
			Choices: []QuestionChoice{},
		}
		questionEventUnparticipated = Question{
			EventID: eventUnparticipated.ID,
			Content: "pqr",
			Choices: []QuestionChoice{},
		}
		helios.DB.Create(&questionSimple)
		helios.DB.Create(&questionWithChoice)
		helios.DB.Create(&questionUnanswered)
		helios.DB.Create(&questionUnowned)
		helios.DB.Create(&questionEvent2)
		helios.DB.Create(&questionEventUnparticipated)

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
