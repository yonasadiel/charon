package exam

import (
	"testing"
	"time"

	"github.com/yonasadiel/helios"

	"github.com/stretchr/testify/assert"
	"github.com/yonasadiel/charon/backend/auth"
)

func TestGetAllVenue(t *testing.T) {
	helios.App.BeforeTest()

	VenueFactorySaved(Venue{})
	VenueFactorySaved(Venue{})

	type getAllEventOfUserTestCase struct {
		user           auth.User
		expectedLength int
		expectedError  helios.Error
	}
	testCases := []getAllEventOfUserTestCase{{
		user:           auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		expectedLength: 2,
	}, {
		user:           auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		expectedLength: 2,
	}, {
		user:          auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal}),
		expectedError: errVenueAccessNotAuthorized,
	}, {
		user:          auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant}),
		expectedError: errVenueAccessNotAuthorized,
	}}
	for i, testCase := range testCases {
		t.Logf("Test GetAllVenue testcase: %d", i)
		var venues []Venue
		var err helios.Error
		venues, err = GetAllVenue(testCase.user)
		if testCase.expectedError == nil {
			assert.Nil(t, err)
			assert.Equal(t, testCase.expectedLength, len(venues))
		} else {
			assert.Equal(t, testCase.expectedError, err)
		}
	}
}

func TestUpsertVenue(t *testing.T) {
	helios.App.BeforeTest()

	type upsertVenueTestCase struct {
		user               auth.User
		venue              Venue
		expectedError      helios.Error
		expectedVenueCount int
	}
	testCases := []upsertVenueTestCase{{
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant}),
		venue:              VenueFactory(Venue{}),
		expectedError:      errVenueAccessNotAuthorized,
		expectedVenueCount: 1,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal}),
		venue:              VenueFactory(Venue{}),
		expectedError:      errVenueAccessNotAuthorized,
		expectedVenueCount: 1,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		venue:              VenueFactory(Venue{}),
		expectedError:      nil,
		expectedVenueCount: 2,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		venue:              VenueFactorySaved(Venue{Name: "New Title"}),
		expectedError:      nil,
		expectedVenueCount: 2,
	}}
	for i, testCase := range testCases {
		var venueCount int
		var venueSaved Venue
		t.Logf("Test UpsertVenue testcase: %d", i)
		err := UpsertVenue(testCase.user, &testCase.venue)
		helios.DB.Model(Venue{}).Count(&venueCount)
		helios.DB.Where("id = ?", testCase.venue.ID).First(&venueSaved)
		assert.Equal(t, testCase.expectedVenueCount, venueCount)
		if testCase.expectedError == nil {
			assert.Nil(t, err)
			assert.Equal(t, testCase.venue.Name, venueSaved.Name, "If the venue has already existed, it should be updated")
		} else {
			assert.Equal(t, testCase.expectedError, err)
		}
	}
}

func TestDeleteVenue(t *testing.T) {
	helios.App.BeforeTest()

	var venue1 Venue = VenueFactorySaved(Venue{})
	var venue2 Venue = VenueFactorySaved(Venue{})
	ParticipationFactorySaved(Participation{Venue: &venue2})

	type deleteVenueTestCase struct {
		user               auth.User
		venueID            uint
		expectedVenue      Venue
		expectedVenueCount int
		expectedError      helios.Error
	}
	testCases := []deleteVenueTestCase{{
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant}),
		venueID:            venue1.ID,
		expectedVenueCount: 2,
		expectedError:      errVenueAccessNotAuthorized,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal}),
		venueID:            venue1.ID,
		expectedVenueCount: 2,
		expectedError:      errVenueAccessNotAuthorized,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		venueID:            999999,
		expectedVenueCount: 2,
		expectedError:      errVenueNotFound,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		venueID:            venue1.ID,
		expectedVenue:      venue1,
		expectedVenueCount: 1,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		venueID:            venue2.ID,
		expectedVenueCount: 1,
		expectedError:      errVenueCantDeletedEventExists,
	}}

	for i, testCase := range testCases {
		t.Logf("Test DeleteVenue testcase: %d", i)
		var venueCount int
		var venue *Venue
		var err helios.Error
		venue, err = DeleteVenue(testCase.user, testCase.venueID)
		helios.DB.Model(&Venue{}).Count(&venueCount)
		assert.Equal(t, testCase.expectedVenueCount, venueCount)
		if testCase.expectedError == nil {
			assert.Nil(t, err)
			assert.Equal(t, testCase.expectedVenue.ID, venue.ID)
			assert.Equal(t, testCase.expectedVenue.Name, venue.Name)
		} else {
			assert.Equal(t, testCase.expectedError, err)
		}
	}
}

func TestGetAllEventOfUser(t *testing.T) {
	helios.App.BeforeTest()

	var userParticipant auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant})
	var userLocal auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal})
	var event1 Event = EventFactorySaved(Event{StartsAt: time.Now().Add(48 * time.Hour)})
	var event2 Event = EventFactorySaved(Event{StartsAt: time.Now().Add(24 * time.Hour)})
	EventFactorySaved(Event{StartsAt: time.Now().Add(72 * time.Hour)})
	ParticipationFactorySaved(Participation{Event: &event1, User: &userParticipant})
	ParticipationFactorySaved(Participation{Event: &event2, User: &userParticipant})
	ParticipationFactorySaved(Participation{Event: &event1, User: &userLocal})

	type getAllEventOfUserTestCase struct {
		user               auth.User
		expectedLength     int
		expectedFirstTitle string
	}
	testCases := []getAllEventOfUserTestCase{{
		user:               userParticipant,
		expectedLength:     2,
		expectedFirstTitle: event2.Title,
	}, {
		user:               userLocal,
		expectedLength:     1,
		expectedFirstTitle: event1.Title,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		expectedLength:     3,
		expectedFirstTitle: event2.Title,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		expectedLength:     3,
		expectedFirstTitle: event2.Title,
	}}
	for i, testCase := range testCases {
		t.Logf("Test GetAllEventOfUser testcase: %d", i)
		events := GetAllEventOfUser(testCase.user)
		assert.Equal(t, testCase.expectedLength, len(events))
		assert.Equal(t, testCase.expectedFirstTitle, events[0].Title, "Events received should be ordered by start time")
	}
}

func TestUpsertEvent(t *testing.T) {
	helios.App.BeforeTest()

	type upsertEventTestCase struct {
		user               auth.User
		event              Event
		expectedError      helios.Error
		expectedEventCount int
	}
	testCases := []upsertEventTestCase{{
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant}),
		event:              EventFactory(Event{}),
		expectedError:      errEventChangeNotAuthorized,
		expectedEventCount: 1,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal}),
		event:              EventFactory(Event{}),
		expectedError:      errEventChangeNotAuthorized,
		expectedEventCount: 1,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		event:              EventFactory(Event{}),
		expectedError:      nil,
		expectedEventCount: 2,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		event:              EventFactorySaved(Event{Title: "New Title"}),
		expectedError:      nil,
		expectedEventCount: 2,
	}}
	for i, testCase := range testCases {
		var eventCount int
		var eventSaved Event
		t.Logf("Test UpsertEvent testcase: %d", i)
		err := UpsertEvent(testCase.user, &testCase.event)
		helios.DB.Model(Event{}).Count(&eventCount)
		helios.DB.Where("id = ?", testCase.event.ID).First(&eventSaved)
		assert.Equal(t, testCase.expectedEventCount, eventCount)
		if testCase.expectedError == nil {
			assert.Nil(t, err)
			assert.Equal(t, testCase.event.Title, eventSaved.Title, "If the event has already existed, it should be updated")
		} else {
			assert.Equal(t, testCase.expectedError, err)
		}
	}
}

func TestGetAllQuestionOfEventAndUser(t *testing.T) {
	helios.App.BeforeTest()

	var userParticipant auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant})
	var userLocal auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal})
	var event1 Event = EventFactorySaved(Event{})
	var event2 Event = EventFactorySaved(Event{StartsAt: time.Now().Add(2 * time.Hour)})
	var participation1 Participation = ParticipationFactorySaved(Participation{Event: &event1, User: &userParticipant})
	var participation2 Participation = ParticipationFactorySaved(Participation{Event: &event1, User: &userLocal})
	ParticipationFactorySaved(Participation{Event: &event2, User: &userParticipant})
	ParticipationFactorySaved(Participation{Event: &event2, User: &userLocal})
	var question1 Question = QuestionFactorySaved(Question{Event: &event1})
	var question2 Question = QuestionFactorySaved(Question{Event: &event1})
	var question3 Question = QuestionFactorySaved(Question{Event: &event1})
	QuestionFactorySaved(Question{Event: &event1})
	QuestionFactorySaved(Question{})
	UserQuestionFactorySaved(UserQuestion{Participation: &participation1, Question: &question1, Ordering: 20, Answer: "abc"})
	UserQuestionFactorySaved(UserQuestion{Participation: &participation1, Question: &question2, Ordering: 10, Answer: "def"})
	UserQuestionFactorySaved(UserQuestion{Participation: &participation2, Question: &question3, Ordering: 10, Answer: "def"})
	type getAllQuestionOfEventAndUserTestCase struct {
		user                        auth.User
		eventID                     uint
		expectedError               helios.Error
		expectedQuestionLen         int
		expectedFirstQuestionAnswer string
	}
	testCases := []getAllQuestionOfEventAndUserTestCase{{
		user:                auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventID:             event1.ID,
		expectedQuestionLen: 4,
	}, {
		user:                auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventID:             event1.ID,
		expectedQuestionLen: 4,
	}, {
		user:                userLocal,
		eventID:             event1.ID,
		expectedQuestionLen: 4,
	}, {
		user:                        userParticipant,
		eventID:                     event1.ID,
		expectedQuestionLen:         2,
		expectedFirstQuestionAnswer: "def",
	}, {
		user:          userParticipant,
		eventID:       event2.ID,
		expectedError: errEventIsNotYetStarted,
	}, {
		user:          userLocal,
		eventID:       event2.ID,
		expectedError: errEventIsNotYetStarted,
	}, {
		user:          auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventID:       99999,
		expectedError: errEventNotFound,
	}, {
		user:                auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventID:             event2.ID,
		expectedQuestionLen: 0,
	}}
	for i, testCase := range testCases {
		t.Logf("Test GetAllQuestionOfEventAndUser testcase: %d", i)
		var questions []Question
		var err helios.Error
		questions, err = GetAllQuestionOfEventAndUser(testCase.user, testCase.eventID)
		if testCase.expectedError == nil {
			assert.Nil(t, err)
			assert.Equal(t, testCase.expectedQuestionLen, len(questions))
			if testCase.expectedFirstQuestionAnswer != "" {
				assert.Equal(t, testCase.expectedFirstQuestionAnswer, questions[0].UserAnswer)
			}
		} else {
			assert.Equal(t, testCase.expectedError, err)
		}
	}
}

func TestUpsertQuestion(t *testing.T) {
	helios.App.BeforeTest()
	var questionCountBefore, choiceCountBefore int

	var userParticipant auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant})
	var userLocal auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal})
	var event1 Event = EventFactorySaved(Event{})
	var event2 Event = EventFactorySaved(Event{})
	var participation1 Participation = ParticipationFactorySaved(Participation{Event: &event1, User: &userParticipant})
	ParticipationFactorySaved(Participation{Event: &event1, User: &userLocal})
	var question1 Question = QuestionFactorySaved(Question{Event: &event1})
	var question2 Question = QuestionFactorySaved(Question{Event: &event1})
	QuestionFactorySaved(Question{Event: &event1})
	QuestionFactorySaved(Question{})
	UserQuestionFactorySaved(UserQuestion{Participation: &participation1, Question: &question1, Ordering: 20, Answer: "abc"})
	UserQuestionFactorySaved(UserQuestion{Participation: &participation1, Question: &question2, Ordering: 10, Answer: "def"})

	helios.DB.Model(&Question{}).Count(&questionCountBefore)
	helios.DB.Model(&QuestionChoice{}).Count(&choiceCountBefore)

	type questionUpsertTestCase struct {
		user          auth.User
		question      Question
		questionCount int
		choiceCount   int
		expectedError helios.Error
	}
	testCases := []questionUpsertTestCase{{
		user:          userParticipant,
		question:      Question{Content: "Content 1", EventID: event1.ID},
		questionCount: questionCountBefore,
		choiceCount:   choiceCountBefore,
		expectedError: errQuestionChangeNotAuthorized,
	}, {
		user:          userLocal,
		question:      Question{Content: "Content 2", EventID: event1.ID},
		questionCount: questionCountBefore,
		choiceCount:   choiceCountBefore,
		expectedError: errQuestionChangeNotAuthorized,
	}, {
		user:          auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		question:      Question{Content: "Content 3", EventID: 9999},
		questionCount: questionCountBefore,
		choiceCount:   choiceCountBefore,
		expectedError: errEventNotFound,
	}, {
		user:          auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		question:      Question{Content: "Content 4", EventID: event1.ID},
		questionCount: questionCountBefore + 1,
		choiceCount:   choiceCountBefore,
		expectedError: nil,
	}, {
		user:          auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		question:      Question{ID: question1.ID, Content: "Content 5", EventID: event2.ID},
		questionCount: questionCountBefore + 1,
		choiceCount:   choiceCountBefore - len(question1.Choices),
		expectedError: nil,
	}, {
		user: auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		question: Question{
			Content: "Content 6",
			EventID: event1.ID,
			Choices: []QuestionChoice{
				{ID: question2.Choices[0].ID, Text: "Choice 6.1"}, // the ID will be ignored
				{Text: "Choice 6.2"},
			},
		},
		questionCount: questionCountBefore + 2,
		choiceCount:   choiceCountBefore - len(question1.Choices) + 2,
		expectedError: nil,
	}, {
		user: auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		question: Question{
			ID:      question2.ID,
			Content: "Content 7",
			EventID: event1.ID,
			Choices: []QuestionChoice{
				{ID: question2.Choices[0].ID, Text: "Choice 7.1"},
				{Text: "Choice 7.2"},
			},
		},
		questionCount: questionCountBefore + 2,
		choiceCount:   choiceCountBefore - len(question1.Choices) + 2 - len(question2.Choices) + 2,
		expectedError: nil,
	}}

	for i, testCase := range testCases {
		var questionCount int
		var choiceCount int
		var questionSaved Question
		t.Logf("Test UpsertQuestion testcase: %d", i)
		err := UpsertQuestion(testCase.user, &testCase.question)
		helios.DB.Model(&Question{}).Count(&questionCount)
		helios.DB.Model(&QuestionChoice{}).Count(&choiceCount)
		helios.DB.Where("id = ?", testCase.question.ID).First(&questionSaved)
		assert.Equal(t, testCase.questionCount, questionCount, "Different number of questions expected")
		assert.Equal(t, testCase.choiceCount, choiceCount, "Different number of question choices expected")
		if testCase.expectedError == nil {
			assert.Nil(t, err, "There should be no error")
			assert.Equal(t, testCase.question.Content, questionSaved.Content, "Different question content")
			assert.Equal(t, testCase.question.EventID, questionSaved.EventID, "Different question event id")
		} else {
			assert.Equal(t, testCase.expectedError, err, "Different error expected")
		}
	}
}

func TestGetQuestionOfEventAndUser(t *testing.T) {
	helios.App.BeforeTest()

	var userParticipant auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant})
	var userLocal auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal})
	var event1 Event = EventFactorySaved(Event{})
	var event2 Event = EventFactorySaved(Event{StartsAt: time.Now().Add(2 * time.Hour)})
	var participation1 Participation = ParticipationFactorySaved(Participation{Event: &event1, User: &userParticipant})
	var participation2 Participation = ParticipationFactorySaved(Participation{Event: &event1, User: &userLocal})
	ParticipationFactorySaved(Participation{Event: &event2, User: &userParticipant})
	ParticipationFactorySaved(Participation{Event: &event2, User: &userLocal})
	var question1 Question = QuestionFactorySaved(Question{Event: &event1})
	var question2 Question = QuestionFactorySaved(Question{Event: &event1})
	var question3 Question = QuestionFactorySaved(Question{Event: &event1})
	var question4 Question = QuestionFactorySaved(Question{Event: &event2})
	QuestionFactorySaved(Question{})
	UserQuestionFactorySaved(UserQuestion{Participation: &participation1, Question: &question1, Ordering: 20, Answer: "abc"})
	UserQuestionFactorySaved(UserQuestion{Participation: &participation1, Question: &question2, Ordering: 10, Answer: "def"})
	UserQuestionFactorySaved(UserQuestion{Participation: &participation2, Question: &question3, Ordering: 10, Answer: "def"})
	type getQuestionOfEventAndUserTestCase struct {
		user                    auth.User
		eventID                 uint
		questionID              uint
		expectedError           helios.Error
		expectedQuestionContent string
		expectedQuestionAnswer  string
	}
	testCases := []getQuestionOfEventAndUserTestCase{{
		user:                    auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventID:                 event2.ID,
		questionID:              question4.ID,
		expectedQuestionContent: question4.Content,
	}, {
		user:                    auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventID:                 event2.ID,
		questionID:              question4.ID,
		expectedQuestionContent: question4.Content,
	}, {
		user:                    userLocal,
		eventID:                 event1.ID,
		questionID:              question3.ID,
		expectedQuestionContent: question3.Content,
	}, {
		user:          userLocal,
		eventID:       event2.ID,
		questionID:    question4.ID,
		expectedError: errEventIsNotYetStarted,
	}, {
		user:          userParticipant,
		eventID:       event2.ID,
		questionID:    question4.ID,
		expectedError: errEventIsNotYetStarted,
	}, {
		user:                    userParticipant,
		eventID:                 event1.ID,
		questionID:              question1.ID,
		expectedQuestionContent: question1.Content,
		expectedQuestionAnswer:  "abc",
	}, {
		user:                    userParticipant,
		eventID:                 event1.ID,
		questionID:              question2.ID,
		expectedQuestionContent: question2.Content,
		expectedQuestionAnswer:  "def",
	}, {
		user:          auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventID:       99999,
		questionID:    question1.ID,
		expectedError: errEventNotFound,
	}, {
		user:          auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventID:       event2.ID,
		questionID:    999999,
		expectedError: errQuestionNotFound,
	}}
	for i, testCase := range testCases {
		t.Logf("Test GetQuestionOfEventAndUser testcase: %d", i)
		var question *Question
		var err helios.Error
		question, err = GetQuestionOfEventAndUser(testCase.user, testCase.eventID, testCase.questionID)
		if testCase.expectedError == nil {
			assert.Nil(t, err)
			assert.Equal(t, testCase.expectedQuestionContent, question.Content)
			assert.Equal(t, testCase.expectedQuestionAnswer, question.UserAnswer)
		} else {
			assert.Equal(t, testCase.expectedError, err)
		}
	}
}

func TestDeleteQuestion(t *testing.T) {
	helios.App.BeforeTest()

	var userParticipant auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant})
	var userLocal auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal})
	var event1 Event = EventFactorySaved(Event{})
	var event2 Event = EventFactorySaved(Event{})
	var participation1 Participation = ParticipationFactorySaved(Participation{Event: &event1, User: &userParticipant})
	var question1 Question = QuestionFactorySaved(Question{Event: &event1})
	var question2 Question = QuestionFactorySaved(Question{Event: &event1})
	UserQuestionFactorySaved(UserQuestion{Participation: &participation1, Question: &question1, Ordering: 20, Answer: "abc"})
	UserQuestionFactorySaved(UserQuestion{Participation: &participation1, Question: &question2, Ordering: 10, Answer: "def"})
	var questionCountBefore, choiceCountBefore, userQuestionCountBefore int
	helios.DB.Model(&Question{}).Count(&questionCountBefore)
	helios.DB.Model(&QuestionChoice{}).Count(&choiceCountBefore)
	helios.DB.Model(&UserQuestion{}).Count(&userQuestionCountBefore)

	type deleteQuestionTestCase struct {
		user                      auth.User
		eventID                   uint
		questionID                uint
		expectedQuestion          Question
		expectedQuestionCount     int
		expectedChoiceCount       int
		expectedUserQuestionCount int
		expectedError             helios.Error
	}
	testCases := []deleteQuestionTestCase{{
		user:                      userParticipant,
		eventID:                   event1.ID,
		questionID:                question1.ID,
		expectedQuestion:          question1,
		expectedQuestionCount:     questionCountBefore,
		expectedChoiceCount:       choiceCountBefore,
		expectedUserQuestionCount: userQuestionCountBefore,
		expectedError:             errQuestionChangeNotAuthorized,
	}, {
		user:                      userLocal,
		eventID:                   event1.ID,
		questionID:                question1.ID,
		expectedQuestion:          question1,
		expectedQuestionCount:     questionCountBefore,
		expectedChoiceCount:       choiceCountBefore,
		expectedUserQuestionCount: userQuestionCountBefore,
		expectedError:             errQuestionChangeNotAuthorized,
	}, {
		user:                      auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventID:                   event1.ID,
		questionID:                23987,
		expectedQuestion:          question1,
		expectedQuestionCount:     questionCountBefore,
		expectedChoiceCount:       choiceCountBefore,
		expectedUserQuestionCount: userQuestionCountBefore,
		expectedError:             errQuestionNotFound,
	}, {
		user:                      auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventID:                   23987,
		questionID:                question1.ID,
		expectedQuestion:          question1,
		expectedQuestionCount:     questionCountBefore,
		expectedChoiceCount:       choiceCountBefore,
		expectedUserQuestionCount: userQuestionCountBefore,
		expectedError:             errEventNotFound,
	}, {
		user:                      auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventID:                   event2.ID,
		questionID:                question1.ID,
		expectedQuestion:          question1,
		expectedQuestionCount:     questionCountBefore,
		expectedChoiceCount:       choiceCountBefore,
		expectedUserQuestionCount: userQuestionCountBefore,
		expectedError:             errQuestionNotFound,
	}, {
		user:                      auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventID:                   event1.ID,
		questionID:                question1.ID,
		expectedQuestion:          question1,
		expectedQuestionCount:     questionCountBefore - 1,
		expectedChoiceCount:       choiceCountBefore - len(question1.Choices),
		expectedUserQuestionCount: userQuestionCountBefore - 1,
		expectedError:             nil,
	}}

	for i, testCase := range testCases {
		var questionCount int
		var choiceCount int
		var userQuestionCount int
		t.Logf("Test DeleteQuestion testcase: %d", i)
		questionDeleted, err := DeleteQuestion(testCase.user, testCase.eventID, testCase.questionID)
		helios.DB.Model(&Question{}).Count(&questionCount)
		helios.DB.Model(&QuestionChoice{}).Count(&choiceCount)
		helios.DB.Model(&UserQuestion{}).Count(&userQuestionCount)
		assert.Equal(t, testCase.expectedQuestionCount, questionCount)
		assert.Equal(t, testCase.expectedChoiceCount, choiceCount)
		assert.Equal(t, testCase.expectedUserQuestionCount, userQuestionCount)
		if testCase.expectedError == nil {
			assert.Nil(t, err)
			assert.Equal(t, testCase.expectedQuestion.ID, questionDeleted.ID)
			assert.Equal(t, testCase.expectedQuestion.Content, questionDeleted.Content)
		} else {
			assert.Equal(t, testCase.expectedError, err)
		}
	}
}

func TestSubmitSubmission(t *testing.T) {
	helios.App.BeforeTest()
	var userParticipant auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant})
	var event1 Event = EventFactorySaved(Event{})
	var event2 Event = EventFactorySaved(Event{StartsAt: time.Now().Add(2 * time.Hour)})
	var question1 Question = QuestionFactorySaved(Question{Event: &event1})
	var question2 Question = QuestionFactorySaved(Question{Event: &event1, Choices: []QuestionChoice{}})
	var question3 Question = QuestionFactorySaved(Question{Event: &event1})
	var question4 Question = QuestionFactorySaved(Question{Event: &event2})
	var randomChoice QuestionChoice = QuestionChoiceFactorySaved(QuestionChoice{})
	var participation1 Participation = ParticipationFactorySaved(Participation{User: &userParticipant, Event: &event1})
	ParticipationFactorySaved(Participation{Event: &event2, User: &userParticipant})
	UserQuestionFactorySaved(UserQuestion{Participation: &participation1, Question: &question1})
	UserQuestionFactorySaved(UserQuestion{Participation: &participation1, Question: &question2})
	type submitSubmissionTestCase struct {
		user          auth.User
		eventID       uint
		questionID    uint
		answer        string
		expectedError helios.Error
	}
	testCases := []submitSubmissionTestCase{{
		user:          auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventID:       event1.ID,
		questionID:    question1.ID,
		answer:        question1.Choices[0].Text,
		expectedError: errSubmissionNotAuthorized,
	}, {
		user:          auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventID:       event1.ID,
		questionID:    question1.ID,
		answer:        question1.Choices[0].Text,
		expectedError: errSubmissionNotAuthorized,
	}, {
		user:          auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal}),
		eventID:       event1.ID,
		questionID:    question1.ID,
		answer:        question1.Choices[0].Text,
		expectedError: errSubmissionNotAuthorized,
	}, {
		user:          userParticipant,
		eventID:       event2.ID,
		questionID:    question4.ID,
		answer:        randomChoice.Text,
		expectedError: errEventIsNotYetStarted,
	}, {
		user:          userParticipant,
		eventID:       event1.ID,
		questionID:    question1.ID,
		answer:        randomChoice.Text,
		expectedError: errAnswerNotValid,
	}, {
		user:       userParticipant,
		eventID:    event1.ID,
		questionID: question1.ID,
		answer:     question1.Choices[0].Text,
	}, {
		user:       userParticipant,
		eventID:    event1.ID,
		questionID: question1.ID,
		answer:     question1.Choices[1].Text,
	}, {
		user:          userParticipant,
		eventID:       999999,
		questionID:    question1.ID,
		answer:        question1.Choices[1].Text,
		expectedError: errEventNotFound,
	}, {
		user:          userParticipant,
		eventID:       event1.ID,
		questionID:    999999,
		answer:        question1.Choices[1].Text,
		expectedError: errQuestionNotFound,
	}, {
		user:       userParticipant,
		eventID:    event1.ID,
		questionID: question2.ID,
		answer:     "answer",
	}, {
		user:          userParticipant,
		eventID:       event1.ID,
		questionID:    question3.ID,
		answer:        question3.Choices[0].Text,
		expectedError: errQuestionNotFound,
	}}
	for i, testCase := range testCases {
		t.Logf("Test SubmitSubmission testcase: %d", i)
		var question *Question
		var errSubmit helios.Error
		var userQuestion UserQuestion
		question, errSubmit = SubmitSubmission(testCase.user, testCase.eventID, testCase.questionID, testCase.answer)
		helios.DB.
			Table("user_questions").
			Joins("inner join participations on participations.id = user_questions.participation_id").
			Where("user_questions.question_id = ?", testCase.questionID).
			Where("participations.user_id = ?", testCase.user.ID).
			First(&userQuestion)
		if testCase.expectedError == nil {
			assert.Nil(t, errSubmit)
			assert.Equal(t, testCase.answer, question.UserAnswer)
			assert.NotEqual(t, 0, userQuestion.ID)
			assert.Equal(t, testCase.answer, userQuestion.Answer)
		} else {
			assert.Equal(t, testCase.expectedError, errSubmit)
		}
	}
}
