package exam

import (
	"crypto/sha256"
	"fmt"
	"strings"
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

	type getAllVenueTestCase struct {
		user           auth.User
		expectedLength int
		expectedError  helios.Error
	}
	testCases := []getAllVenueTestCase{{
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
		user                       auth.User
		event                      Event
		expectedError              helios.Error
		expectedEventCount         int
		expectedParticipationCount int
	}
	testCases := []upsertEventTestCase{{
		user:                       auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant}),
		event:                      EventFactory(Event{}),
		expectedError:              errEventChangeNotAuthorized,
		expectedEventCount:         1,
		expectedParticipationCount: 0,
	}, {
		user:                       auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal}),
		event:                      Event{},
		expectedEventCount:         2,
		expectedParticipationCount: 1,
	}, {
		user:                       auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		event:                      EventFactory(Event{}),
		expectedError:              nil,
		expectedEventCount:         3,
		expectedParticipationCount: 1,
	}, {
		user:                       auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		event:                      EventFactorySaved(Event{Title: "New Title"}),
		expectedError:              nil,
		expectedEventCount:         3,
		expectedParticipationCount: 1,
	}}
	for i, testCase := range testCases {
		var eventCount int
		var eventSaved Event
		var participationCount int
		t.Logf("Test UpsertEvent testcase: %d", i)
		err := UpsertEvent(testCase.user, &testCase.event)
		helios.DB.Model(Event{}).Count(&eventCount)
		helios.DB.Model(Participation{}).Count(&participationCount)
		helios.DB.Where("id = ?", testCase.event.ID).First(&eventSaved)
		assert.Equal(t, testCase.expectedEventCount, eventCount)
		assert.Equal(t, testCase.expectedParticipationCount, participationCount)
		if testCase.expectedError == nil {
			assert.Nil(t, err)
			assert.Equal(t, testCase.event.Title, eventSaved.Title, "If the event has already existed, it should be updated")
			if !testCase.user.IsLocal() {
				assert.NotEmpty(t, eventSaved.SimKey)
				assert.NotEmpty(t, eventSaved.PubKey)
				assert.NotEmpty(t, eventSaved.PrvKey)
				assert.NotEmpty(t, testCase.event.SimKey)
				assert.NotEmpty(t, testCase.event.PubKey)
				assert.NotEmpty(t, testCase.event.PrvKey)
			} else {
				assert.Empty(t, eventSaved.SimKey)
				assert.Empty(t, eventSaved.PubKey)
				assert.Empty(t, eventSaved.PrvKey)
				assert.Empty(t, testCase.event.SimKey)
				assert.Empty(t, testCase.event.PubKey)
				assert.Empty(t, testCase.event.PrvKey)
			}
		} else {
			assert.Equal(t, testCase.expectedError, err)
		}
	}
}

func TestGetAllParticipationOfUserAndEvent(t *testing.T) {
	helios.App.BeforeTest()

	var event1 Event = EventFactorySaved(Event{})
	var event2 Event = EventFactorySaved(Event{})
	var userParticipant auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant})
	var userLocal auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal})
	ParticipationFactorySaved(Participation{Event: &event1, User: &userLocal})
	ParticipationFactorySaved(Participation{Event: &event1, User: &userParticipant})
	ParticipationFactorySaved(Participation{Event: &event2})

	type getAllParticipationOfUserAndEventTestCase struct {
		user           auth.User
		eventSlug      string
		expectedLength int
		expectedError  helios.Error
	}
	testCases := []getAllParticipationOfUserAndEventTestCase{{
		user:           auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventSlug:      event1.Slug,
		expectedLength: 2, // userParticipant + userLocal
	}, {
		user:           auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:      event2.Slug,
		expectedLength: 1, // userLocal
	}, {
		user:           userLocal,
		eventSlug:      event1.Slug,
		expectedLength: 2, // userParticipant + userLocal
	}, {
		user:          userLocal,
		eventSlug:     event2.Slug,
		expectedError: errEventNotFound,
	}, {
		user:           userParticipant,
		eventSlug:      event1.Slug,
		expectedLength: 1, // userParticipant
	}}
	for i, testCase := range testCases {
		t.Logf("Test GetAllParticipation testcase: %d", i)
		var participations []Participation
		var err helios.Error
		participations, err = GetAllParticipationOfUserAndEvent(testCase.user, testCase.eventSlug)
		if testCase.expectedError == nil {
			assert.Nil(t, err)
			assert.Equal(t, testCase.expectedLength, len(participations))
		} else {
			assert.Equal(t, testCase.expectedError, err)
		}
	}
}

func TestUpsertParticipation(t *testing.T) {
	helios.App.BeforeTest()
	var participationCountBefore int

	var venue1 Venue = VenueFactorySaved(Venue{})
	var venue2 Venue = VenueFactorySaved(Venue{})
	var userParticipant auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant})
	var userLocal auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal})
	var event1 Event = EventFactorySaved(Event{})
	var event2 Event = EventFactorySaved(Event{})
	ParticipationFactorySaved(Participation{Event: &event1, User: &userLocal})
	var key = "secret_key_per_user_per_event"
	var keyHashedOnce = "e33a9931ec1ba26e9acd8957b597595ce7e336e4df534ac83bc4102e963c4814"
	var keyHashedTwice = "cfa42ce14740fb597b001bdc9c6a2569c027f53358f7fd2ebdc80d0888737530"

	helios.DB.Model(&Participation{}).Count(&participationCountBefore)

	type participationUpsertTestCase struct {
		user                       auth.User
		eventSlug                  string
		userUsername               string
		participation              Participation
		expectedParticipationCount int
		expectedError              helios.Error
	}
	testCases := []participationUpsertTestCase{{
		user:                       auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal}),
		eventSlug:                  event2.Slug,
		userUsername:               userParticipant.Username,
		participation:              Participation{VenueID: venue1.ID, KeyPlain: key},
		expectedParticipationCount: participationCountBefore,
		expectedError:              errEventNotFound,
	}, {
		user:                       auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:                  event1.Slug,
		userUsername:               "random_username",
		participation:              Participation{VenueID: venue1.ID, KeyPlain: key},
		expectedParticipationCount: participationCountBefore,
		expectedError:              errUserNotFound,
	}, {
		user:                       userLocal,
		eventSlug:                  event1.Slug,
		userUsername:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal}).Username,
		participation:              Participation{VenueID: venue1.ID, KeyPlain: key},
		expectedParticipationCount: participationCountBefore,
		expectedError:              errParticipationChangeNotAuthorized,
	}, {
		user:                       userLocal,
		eventSlug:                  event1.Slug,
		userUsername:               userParticipant.Username,
		participation:              Participation{VenueID: 123, KeyPlain: key},
		expectedParticipationCount: participationCountBefore,
		expectedError:              errVenueNotFound,
	}, {
		user:                       userLocal,
		eventSlug:                  event1.Slug,
		userUsername:               userParticipant.Username,
		participation:              Participation{EventID: event2.ID, UserID: userLocal.ID, VenueID: venue1.ID, KeyPlain: key},
		expectedParticipationCount: participationCountBefore + 1,
	}, {
		user:                       auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:                  event1.Slug,
		userUsername:               userLocal.Username,
		participation:              Participation{EventID: event2.ID, UserID: userLocal.ID, VenueID: venue2.ID, KeyPlain: key},
		expectedParticipationCount: participationCountBefore + 1,
	}}

	for i, testCase := range testCases {
		var participationCount int
		var participationSaved Participation
		var tempVenueID uint = testCase.participation.VenueID
		var err helios.Error
		t.Logf("Test UpsertParticipation testcase: %d", i)
		err = UpsertParticipation(testCase.user, testCase.eventSlug, testCase.userUsername, &testCase.participation)
		helios.DB.Model(&Participation{}).Count(&participationCount)
		helios.DB.Preload("User").Preload("Event").Preload("Venue").Where("id = ?", testCase.participation.ID).First(&participationSaved)
		assert.Equal(t, testCase.expectedParticipationCount, participationCount)
		if testCase.expectedError == nil {
			assert.Nil(t, err)
			assert.Equal(t, testCase.eventSlug, participationSaved.Event.Slug)
			assert.Equal(t, testCase.eventSlug, testCase.participation.Event.Slug)
			assert.Equal(t, testCase.userUsername, participationSaved.User.Username)
			assert.Equal(t, testCase.userUsername, testCase.participation.User.Username)
			assert.Equal(t, keyHashedOnce, participationSaved.KeyHashedOnce)
			assert.Equal(t, keyHashedTwice, participationSaved.KeyHashedTwice)
			assert.Equal(t, tempVenueID, participationSaved.Venue.ID)
			assert.Equal(t, tempVenueID, testCase.participation.Venue.ID)
		} else {
			assert.Equal(t, testCase.expectedError, err)
		}
	}
}

func TestVerifyParticipation(t *testing.T) {
	helios.App.BeforeTest()

	var userParticipant1 auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant})
	var userParticipant2 auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant})
	var event1 Event = EventFactorySaved(Event{})
	var event2 Event = EventFactorySaved(Event{})
	var event1User1Key = "event_1_user1"
	var event1User2Key = "event_1_user2"
	var event1User1KeyHashedOnce = fmt.Sprintf("%x", sha256.Sum256([]byte(event1User1Key)))
	var event1User2KeyHashedOnce = fmt.Sprintf("%x", sha256.Sum256([]byte(event1User2Key)))
	var event1User1KeyHashedTwice = fmt.Sprintf("%x", sha256.Sum256([]byte(event1User1KeyHashedOnce)))
	var event1User2KeyHashedTwice = fmt.Sprintf("%x", sha256.Sum256([]byte(event1User2KeyHashedOnce)))
	var participation1 Participation = ParticipationFactorySaved(Participation{Event: &event1, User: &userParticipant1, KeyHashedTwice: event1User1KeyHashedTwice})
	var participation2 Participation = ParticipationFactorySaved(Participation{Event: &event1, User: &userParticipant2, KeyHashedTwice: event1User2KeyHashedTwice})
	type verifyParticipationTestCase struct {
		user           auth.User
		eventSlug      string
		keyHashedOnce  string
		participation  Participation
		expectedResult bool
		expectedError  helios.Error
	}
	testCases := []verifyParticipationTestCase{{
		user:          userParticipant1,
		eventSlug:     event1.Slug,
		keyHashedOnce: "wrong_key",
		participation: participation1,
		expectedError: errParticipationWrongKey,
	}, {
		user:          userParticipant1,
		eventSlug:     event1.Slug,
		keyHashedOnce: event1User1KeyHashedOnce,
		participation: participation1,
	}, {
		user:          userParticipant2,
		eventSlug:     event1.Slug,
		keyHashedOnce: event1User1KeyHashedOnce,
		participation: participation2,
		expectedError: errParticipationWrongKey,
	}, {
		user:          userParticipant2,
		eventSlug:     event1.Slug,
		keyHashedOnce: event1User2KeyHashedOnce,
		participation: participation2,
	}, {
		user:          userParticipant2,
		eventSlug:     event2.Slug,
		keyHashedOnce: "wrong_key",
		participation: participation2,
		expectedError: errEventNotFound,
	}}
	for i, testCase := range testCases {
		t.Logf("Test VerifyParticipation testcase: %d", i)
		var err helios.Error
		var participationSaved Participation
		err = VerifyParticipation(testCase.user, testCase.eventSlug, testCase.keyHashedOnce)
		helios.DB.Where("id = ?", testCase.participation.ID).First(&participationSaved)
		if testCase.expectedError == nil {
			assert.Nil(t, err)
			assert.Equal(t, testCase.keyHashedOnce, participationSaved.KeyHashedOnce)
		} else {
			assert.Equal(t, testCase.expectedError, err)
		}
	}
}

func TestDeleteParticipation(t *testing.T) {
	helios.App.BeforeTest()

	var userParticipant auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant})
	var userLocal1 auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal})
	var event1 Event = EventFactorySaved(Event{})
	var event2 Event = EventFactorySaved(Event{})
	var participation1 Participation = ParticipationFactorySaved(Participation{Event: &event1, User: &userParticipant})
	var participation2 Participation = ParticipationFactorySaved(Participation{Event: &event1, User: &userLocal1})
	UserQuestionFactorySaved(UserQuestion{Participation: &participation1})
	UserQuestionFactorySaved(UserQuestion{Participation: &participation1})
	var participationCountBefore, userQuestionCountBefore int
	helios.DB.Model(&Participation{}).Count(&participationCountBefore)
	helios.DB.Model(&UserQuestion{}).Count(&userQuestionCountBefore)

	type deleteParticipationTestCase struct {
		user                       auth.User
		eventSlug                  string
		participationID            uint
		expectedParticipation      Participation
		expectedParticipationCount int
		expectedUserQuestionCount  int
		expectedError              helios.Error
	}
	testCases := []deleteParticipationTestCase{{
		user:                       auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:                  "random_slug",
		participationID:            participation1.ID,
		expectedParticipationCount: participationCountBefore,
		expectedUserQuestionCount:  userQuestionCountBefore,
		expectedError:              errEventNotFound,
	}, {
		user:                       auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:                  event2.Slug,
		participationID:            participation1.ID,
		expectedParticipationCount: participationCountBefore,
		expectedUserQuestionCount:  userQuestionCountBefore,
		expectedError:              errParticipationNotFound,
	}, {
		user:                       userParticipant,
		eventSlug:                  event1.Slug,
		participationID:            participation2.ID,
		expectedParticipationCount: participationCountBefore,
		expectedUserQuestionCount:  userQuestionCountBefore,
		expectedError:              errParticipationChangeNotAuthorized,
	}, {
		user:                       userLocal1,
		eventSlug:                  event1.Slug,
		participationID:            participation1.ID,
		expectedParticipation:      participation1,
		expectedParticipationCount: participationCountBefore - 1,
		expectedUserQuestionCount:  userQuestionCountBefore - 2,
	}, {
		user:                       auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:                  event1.Slug,
		participationID:            participation2.ID,
		expectedParticipation:      participation2,
		expectedParticipationCount: participationCountBefore - 2,
		expectedUserQuestionCount:  userQuestionCountBefore - 2,
	}}

	for i, testCase := range testCases {
		var questionCount int
		var userQuestionCount int
		t.Logf("Test DeleteParticipation testcase: %d", i)
		questionDeleted, err := DeleteParticipation(testCase.user, testCase.eventSlug, testCase.participationID)
		helios.DB.Model(&Participation{}).Count(&questionCount)
		helios.DB.Model(&UserQuestion{}).Count(&userQuestionCount)
		assert.Equal(t, testCase.expectedParticipationCount, questionCount)
		assert.Equal(t, testCase.expectedUserQuestionCount, userQuestionCount)
		if testCase.expectedError == nil {
			assert.Nil(t, err)
			assert.Equal(t, testCase.expectedParticipation.ID, questionDeleted.ID)
			assert.Equal(t, testCase.expectedParticipation.User.Username, questionDeleted.User.Username)
			assert.Equal(t, testCase.expectedParticipation.Venue.ID, questionDeleted.Venue.ID)
		} else {
			assert.Equal(t, testCase.expectedError, err)
		}
	}
}

func TestGetAllQuestionOfUserAndEvent(t *testing.T) {
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
	type getAllQuestionOfUserAndEventTestCase struct {
		user                        auth.User
		eventSlug                   string
		expectedError               helios.Error
		expectedQuestionLen         int
		expectedFirstQuestionAnswer string
	}
	testCases := []getAllQuestionOfUserAndEventTestCase{{
		user:                auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventSlug:           event1.Slug,
		expectedQuestionLen: 4,
	}, {
		user:                auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:           event1.Slug,
		expectedQuestionLen: 4,
	}, {
		user:                userLocal,
		eventSlug:           event1.Slug,
		expectedQuestionLen: 4,
	}, {
		user:                        userParticipant,
		eventSlug:                   event1.Slug,
		expectedQuestionLen:         2,
		expectedFirstQuestionAnswer: "def",
	}, {
		user:          userParticipant,
		eventSlug:     event2.Slug,
		expectedError: errEventIsNotYetStarted,
	}, {
		user:          userLocal,
		eventSlug:     event2.Slug,
		expectedError: errEventIsNotYetStarted,
	}, {
		user:          auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventSlug:     "99999",
		expectedError: errEventNotFound,
	}, {
		user:                auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventSlug:           event2.Slug,
		expectedQuestionLen: 0,
	}}
	for i, testCase := range testCases {
		t.Logf("Test GetAllQuestionOfUserAndEvent testcase: %d", i)
		var questions []Question
		var err helios.Error
		questions, err = GetAllQuestionOfUserAndEvent(testCase.user, testCase.eventSlug)
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
	var questionCountBefore int

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

	type questionUpsertTestCase struct {
		user                  auth.User
		eventSlug             string
		question              Question
		expectedQuestionCount int
		expectedError         helios.Error
	}
	testCases := []questionUpsertTestCase{{
		user:                  userParticipant,
		eventSlug:             event1.Slug,
		question:              Question{Content: "Content 1", EventID: event1.ID},
		expectedQuestionCount: questionCountBefore,
		expectedError:         errQuestionChangeNotAuthorized,
	}, {
		user:                  userLocal,
		eventSlug:             event1.Slug,
		question:              Question{Content: "Content 2", EventID: event1.ID},
		expectedQuestionCount: questionCountBefore,
		expectedError:         errQuestionChangeNotAuthorized,
	}, {
		user:                  auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventSlug:             "9999",
		question:              Question{Content: "Content 3", EventID: event1.ID},
		expectedQuestionCount: questionCountBefore,
		expectedError:         errEventNotFound,
	}, {
		user:                  auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventSlug:             event1.Slug,
		question:              Question{Content: "Content 4", EventID: event1.ID},
		expectedQuestionCount: questionCountBefore + 1,
		expectedError:         nil,
	}, {
		user:                  auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:             event1.Slug,
		question:              Question{ID: question1.ID, Content: "Content 5", EventID: event2.ID},
		expectedQuestionCount: questionCountBefore + 1,
		expectedError:         nil,
	}, {
		user:      auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug: event1.Slug,
		question: Question{
			Content: "Content 6",
			EventID: event1.ID,
			Choices: "Choice 6.1|Choice 6.2",
		},
		expectedQuestionCount: questionCountBefore + 2,
		expectedError:         nil,
	}, {
		user:      auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug: event1.Slug,
		question: Question{
			ID:      question2.ID,
			Content: "Content 7",
			EventID: event1.ID,
			Choices: "Choice 7.1|Choice 7.2",
		},
		expectedQuestionCount: questionCountBefore + 2,
		expectedError:         nil,
	}}

	for i, testCase := range testCases {
		var questionCount int
		var questionSaved Question
		t.Logf("Test UpsertQuestion testcase: %d", i)
		err := UpsertQuestion(testCase.user, testCase.eventSlug, &testCase.question)
		helios.DB.Model(&Question{}).Count(&questionCount)
		helios.DB.Where("id = ?", testCase.question.ID).First(&questionSaved)
		assert.Equal(t, testCase.expectedQuestionCount, questionCount)
		if testCase.expectedError == nil {
			assert.Nil(t, err)
			assert.Equal(t, testCase.question.Content, questionSaved.Content)
			assert.Equal(t, testCase.question.EventID, questionSaved.EventID)
		} else {
			assert.Equal(t, testCase.expectedError, err)
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
	UserQuestionFactorySaved(UserQuestion{Participation: &participation1, Question: &question1, Ordering: 20, Answer: "abc"}) // questionNumber 2
	UserQuestionFactorySaved(UserQuestion{Participation: &participation1, Question: &question2, Ordering: 10, Answer: "def"}) // questionNumber 1
	UserQuestionFactorySaved(UserQuestion{Participation: &participation2, Question: &question3, Ordering: 10, Answer: "def"}) // questionNumber 1
	type getQuestionOfEventAndUserTestCase struct {
		user                    auth.User
		eventSlug               string
		questionNumber          uint
		expectedError           helios.Error
		expectedQuestionContent string
		expectedQuestionAnswer  string
	}
	testCases := []getQuestionOfEventAndUserTestCase{{
		user:                    auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventSlug:               event2.Slug,
		questionNumber:          1,
		expectedQuestionContent: question4.Content,
	}, {
		user:                    auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:               event2.Slug,
		questionNumber:          1,
		expectedQuestionContent: question4.Content,
	}, {
		user:                    userLocal,
		eventSlug:               event1.Slug,
		questionNumber:          3,
		expectedQuestionContent: question3.Content,
	}, {
		user:           userLocal,
		eventSlug:      event2.Slug,
		questionNumber: 1,
		expectedError:  errEventIsNotYetStarted,
	}, {
		user:           userParticipant,
		eventSlug:      event2.Slug,
		questionNumber: 1,
		expectedError:  errEventIsNotYetStarted,
	}, {
		user:                    userParticipant,
		eventSlug:               event1.Slug,
		questionNumber:          2,
		expectedQuestionContent: question1.Content,
		expectedQuestionAnswer:  "abc",
	}, {
		user:                    userParticipant,
		eventSlug:               event1.Slug,
		questionNumber:          1,
		expectedQuestionContent: question2.Content,
		expectedQuestionAnswer:  "def",
	}, {
		user:           auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventSlug:      "99999",
		questionNumber: 2,
		expectedError:  errEventNotFound,
	}, {
		user:           auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventSlug:      event2.Slug,
		questionNumber: 999999,
		expectedError:  errQuestionNotFound,
	}}
	for i, testCase := range testCases {
		t.Logf("Test GetQuestionOfEventAndUser testcase: %d", i)
		var question *Question
		var err helios.Error
		question, err = GetQuestionOfEventAndUser(testCase.user, testCase.eventSlug, testCase.questionNumber)
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
	UserQuestionFactorySaved(UserQuestion{Participation: &participation1, Question: &question1, Ordering: 20, Answer: "abc"}) // 2
	UserQuestionFactorySaved(UserQuestion{Participation: &participation1, Question: &question2, Ordering: 10, Answer: "def"}) // 1
	var questionCountBefore, userQuestionCountBefore int
	helios.DB.Model(&Question{}).Count(&questionCountBefore)
	helios.DB.Model(&UserQuestion{}).Count(&userQuestionCountBefore)

	type deleteQuestionTestCase struct {
		user                      auth.User
		eventSlug                 string
		questionNumber            uint
		expectedQuestion          Question
		expectedQuestionCount     int
		expectedUserQuestionCount int
		expectedError             helios.Error
	}
	testCases := []deleteQuestionTestCase{{
		user:                      userParticipant,
		eventSlug:                 event1.Slug,
		questionNumber:            1,
		expectedQuestionCount:     questionCountBefore,
		expectedUserQuestionCount: userQuestionCountBefore,
		expectedError:             errQuestionChangeNotAuthorized,
	}, {
		user:                      userLocal,
		eventSlug:                 event1.Slug,
		questionNumber:            1,
		expectedQuestionCount:     questionCountBefore,
		expectedUserQuestionCount: userQuestionCountBefore,
		expectedError:             errQuestionChangeNotAuthorized,
	}, {
		user:                      auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:                 event1.Slug,
		questionNumber:            23987,
		expectedQuestionCount:     questionCountBefore,
		expectedUserQuestionCount: userQuestionCountBefore,
		expectedError:             errQuestionNotFound,
	}, {
		user:                      auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventSlug:                 "23987",
		questionNumber:            1,
		expectedQuestionCount:     questionCountBefore,
		expectedUserQuestionCount: userQuestionCountBefore,
		expectedError:             errEventNotFound,
	}, {
		user:                      auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventSlug:                 event2.Slug,
		questionNumber:            1,
		expectedQuestionCount:     questionCountBefore,
		expectedUserQuestionCount: userQuestionCountBefore,
		expectedError:             errQuestionNotFound,
	}, {
		user:                      auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventSlug:                 event1.Slug,
		questionNumber:            1,
		expectedQuestion:          question1,
		expectedQuestionCount:     questionCountBefore - 1,
		expectedUserQuestionCount: userQuestionCountBefore - 1,
	}}

	for i, testCase := range testCases {
		var questionCount int
		var userQuestionCount int
		t.Logf("Test DeleteQuestion testcase: %d", i)
		questionDeleted, err := DeleteQuestion(testCase.user, testCase.eventSlug, testCase.questionNumber)
		helios.DB.Model(&Question{}).Count(&questionCount)
		helios.DB.Model(&UserQuestion{}).Count(&userQuestionCount)
		assert.Equal(t, testCase.expectedQuestionCount, questionCount)
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
	var question2 Question = QuestionFactorySaved(Question{Event: &event1, Choices: "|"})
	var question3 Question = QuestionFactorySaved(Question{Event: &event1})
	var question4 Question = QuestionFactorySaved(Question{Event: &event2})
	var participation1 Participation = ParticipationFactorySaved(Participation{User: &userParticipant, Event: &event1})
	ParticipationFactorySaved(Participation{Event: &event2, User: &userParticipant})
	UserQuestionFactorySaved(UserQuestion{Participation: &participation1, Question: &question1, Ordering: 20}) // questionNumber 2
	UserQuestionFactorySaved(UserQuestion{Participation: &participation1, Question: &question2, Ordering: 10}) // questionNumber 1
	type submitSubmissionTestCase struct {
		user           auth.User
		eventSlug      string
		questionNumber uint
		answer         string
		expectedError  helios.Error
	}
	testCases := []submitSubmissionTestCase{{
		user:           auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventSlug:      event1.Slug,
		questionNumber: 2,
		answer:         strings.Split(question1.Choices, "|")[0],
		expectedError:  errSubmissionNotAuthorized,
	}, {
		user:           auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:      event1.Slug,
		questionNumber: 2,
		answer:         strings.Split(question1.Choices, "|")[0],
		expectedError:  errSubmissionNotAuthorized,
	}, {
		user:           auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal}),
		eventSlug:      event1.Slug,
		questionNumber: 2,
		answer:         strings.Split(question1.Choices, "|")[0],
		expectedError:  errSubmissionNotAuthorized,
	}, {
		user:           userParticipant,
		eventSlug:      event2.Slug,
		questionNumber: question4.ID,
		answer:         "random",
		expectedError:  errEventIsNotYetStarted,
	}, {
		user:           userParticipant,
		eventSlug:      event1.Slug,
		questionNumber: 2,
		answer:         strings.Split(question1.Choices, "|")[0],
	}, {
		user:           userParticipant,
		eventSlug:      event1.Slug,
		questionNumber: 2,
		answer:         strings.Split(question1.Choices, "|")[1],
	}, {
		user:           userParticipant,
		eventSlug:      "999999",
		questionNumber: 2,
		answer:         strings.Split(question1.Choices, "|")[1],
		expectedError:  errEventNotFound,
	}, {
		user:           userParticipant,
		eventSlug:      event1.Slug,
		questionNumber: 999999,
		answer:         strings.Split(question1.Choices, "|")[1],
		expectedError:  errQuestionNotFound,
	}, {
		user:           userParticipant,
		eventSlug:      event1.Slug,
		questionNumber: 1,
		answer:         "answer",
	}, {
		user:           userParticipant,
		eventSlug:      event1.Slug,
		questionNumber: 3,
		answer:         strings.Split(question3.Choices, "|")[0],
		expectedError:  errQuestionNotFound,
	}}
	for i, testCase := range testCases {
		t.Logf("Test SubmitSubmission testcase: %d", i)
		var question *Question
		var errSubmit helios.Error
		var userQuestion UserQuestion
		question, errSubmit = SubmitSubmission(testCase.user, testCase.eventSlug, testCase.questionNumber, testCase.answer)
		helios.DB.
			Table("user_questions").
			Joins("inner join participations on participations.id = user_questions.participation_id").
			Order("user_questions.ordering asc").
			Offset(testCase.questionNumber - 1).
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

func TestGetParticipationStatus(t *testing.T) {
	helios.App.BeforeTest()

	var user1 auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant})
	var user2 auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant})
	var userLocal auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal})
	var event1 Event = EventFactorySaved(Event{})
	var event2 Event = EventFactorySaved(Event{})
	var notNilTime time.Time = time.Now()
	ParticipationFactorySaved(Participation{Event: &event1, User: &userLocal})
	ParticipationFactorySaved(Participation{Event: &event1, User: &user1})
	ParticipationFactorySaved(Participation{Event: &event1, User: &user2})
	var session auth.Session = auth.Session{
		ID:        1,
		User:      &user2,
		Token:     "abc",
		IPAddress: "192.168.0.2",
	}
	helios.DB.Create(&session)
	type getParticipationStatusTestCase struct {
		user           auth.User
		eventSlug      string
		expectedStatus []ParticipationStatus
		expectedError  helios.Error
	}
	testCases := []getParticipationStatusTestCase{{
		user:          auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:     event1.Slug,
		expectedError: errParticipationStatusAccessNotAuthorized,
	}, {
		user:          userLocal,
		eventSlug:     event2.Slug,
		expectedError: errEventNotFound,
	}, {
		user:      userLocal,
		eventSlug: event1.Slug,
		expectedStatus: []ParticipationStatus{{
			UserUsername: user1.Username,
			IPAddress:    "",
			LoginAt:      nil,
			SessionID:    0,
		}, {
			UserUsername: user2.Username,
			IPAddress:    "192.168.0.2",
			LoginAt:      &notNilTime,
			SessionID:    session.ID,
		}},
	}}
	for i, testCase := range testCases {
		t.Logf("Test GetParticipationStatus testcase: %d", i)
		var status []ParticipationStatus
		var err helios.Error
		status, err = GetParticipationStatus(testCase.user, testCase.eventSlug)
		if testCase.expectedError == nil {
			assert.Equal(t, len(testCase.expectedStatus), len(status))
			for j := range status {
				assert.Equal(t, testCase.expectedStatus[j].UserUsername, status[j].UserUsername)
				if testCase.expectedStatus[j].LoginAt == nil {
					assert.Empty(t, status[j].LoginAt)
					assert.Empty(t, status[j].IPAddress)
				} else {
					assert.NotNil(t, status[j].LoginAt)
					assert.Equal(t, testCase.expectedStatus[j].IPAddress, status[j].IPAddress)
				}
			}
		} else {
			assert.Equal(t, testCase.expectedError, err)
		}
	}
}

func TestRemoveParticipationSession(t *testing.T) {
	helios.App.BeforeTest()

	var user1 auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant})
	var userLocal auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal})
	var event1 Event = EventFactorySaved(Event{})
	var event2 Event = EventFactorySaved(Event{})
	ParticipationFactorySaved(Participation{Event: &event1, User: &userLocal})
	ParticipationFactorySaved(Participation{Event: &event1, User: &user1})
	var session auth.Session = auth.Session{
		User:      &user1,
		Token:     "abc",
		IPAddress: "192.168.0.2",
	}
	helios.DB.Create(&session)
	type removeParticipationSessionTestCase struct {
		user          auth.User
		eventSlug     string
		sessionID     uint
		expectedError helios.Error
	}
	testCases := []removeParticipationSessionTestCase{{
		user:          auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:     event1.Slug,
		sessionID:     session.ID,
		expectedError: errParticipationStatusAccessNotAuthorized,
	}, {
		user:          userLocal,
		eventSlug:     event2.Slug,
		sessionID:     session.ID,
		expectedError: errEventNotFound,
	}, {
		user:          userLocal,
		eventSlug:     event1.Slug,
		sessionID:     12,
		expectedError: errParticipationStatusNotFound,
	}, {
		user:      userLocal,
		eventSlug: event1.Slug,
		sessionID: session.ID,
	}, {
		user:          userLocal,
		eventSlug:     event1.Slug,
		sessionID:     session.ID,
		expectedError: errParticipationStatusNotFound,
	}}
	for i, testCase := range testCases {
		t.Logf("Test RemoveParticipationSession testcase: %d", i)
		var err helios.Error
		var sessionSaved auth.Session
		err = RemoveParticipationSession(testCase.user, testCase.eventSlug, testCase.sessionID)
		helios.DB.Where("id = ?", session.ID).First(&sessionSaved)
		if testCase.expectedError == nil {
			assert.Nil(t, err)
			assert.Empty(t, sessionSaved.ID)
		} else {
			assert.Equal(t, testCase.expectedError, err)
		}
	}
}

func TestGetSynchronizationData(t *testing.T) {
	helios.App.BeforeTest()

	var userLocal auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal})
	var venue Venue = VenueFactorySaved(Venue{})
	var event1 Event = EventFactorySaved(Event{})
	var event2 Event = EventFactorySaved(Event{})
	QuestionFactorySaved(Question{Event: &event1})
	QuestionFactorySaved(Question{Event: &event1})
	QuestionFactorySaved(Question{Event: &event2})
	participations := []Participation{
		ParticipationFactorySaved(Participation{Event: &event1, User: &userLocal, Venue: &venue, KeyHashedTwice: "key1"}),
		ParticipationFactorySaved(Participation{Event: &event1, Venue: &venue, KeyHashedTwice: "key2"}),
		ParticipationFactorySaved(Participation{Event: &event1, Venue: &venue, KeyHashedTwice: "key3"}),
	}
	ParticipationFactorySaved(Participation{Event: &event1})
	ParticipationFactorySaved(Participation{Event: &event1})
	ParticipationFactorySaved(Participation{Event: &event2})
	expectedUsersKey := make(map[string]string)
	for _, participation := range participations {
		expectedUsersKey[participation.User.Username] = participation.KeyHashedTwice
	}
	type getSynchronizationDataTestCase struct {
		user                   auth.User
		eventSlug              string
		expectedEvent          Event
		expectedVenue          Venue
		expectedQuestionLength int
		expectedUserLength     int
		expectedUsersKey       map[string]string
		expectedError          helios.Error
	}
	testCases := []getSynchronizationDataTestCase{{
		user:          auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventSlug:     event1.Slug,
		expectedError: errSynchronizationNotAuthorized,
	}, {
		user:          auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:     event1.Slug,
		expectedError: errSynchronizationNotAuthorized,
	}, {
		user:          auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant}),
		eventSlug:     event1.Slug,
		expectedError: errSynchronizationNotAuthorized,
	}, {
		user:          userLocal,
		eventSlug:     "abc",
		expectedError: errEventNotFound,
	}, {
		user:          userLocal,
		eventSlug:     event2.Slug,
		expectedError: errEventNotFound,
	}, {
		user:                   userLocal,
		eventSlug:              event1.Slug,
		expectedEvent:          event1,
		expectedVenue:          venue,
		expectedQuestionLength: 2,
		expectedUserLength:     3,
		expectedUsersKey:       expectedUsersKey,
	}}
	for i, testCase := range testCases {
		t.Logf("Test GetSynchronizationData testcase: %d", i)
		var event *Event
		var venue *Venue
		var questions []Question
		var users []auth.User
		var usersKey map[string]string
		var err helios.Error
		event, venue, questions, users, usersKey, err = GetSynchronizationData(testCase.user, testCase.eventSlug)
		if testCase.expectedError == nil {
			assert.Nil(t, err)
			assert.Equal(t, testCase.expectedEvent.Title, event.Title)
			assert.Equal(t, testCase.expectedVenue.Name, venue.Name)
			assert.Equal(t, testCase.expectedQuestionLength, len(questions))
			assert.Equal(t, testCase.expectedUserLength, len(users))
			assert.Equal(t, testCase.expectedUsersKey, usersKey)
		} else {
			assert.Equal(t, testCase.expectedError, err)
		}
	}
}

func TestPutSynchronizationData(t *testing.T) {
	helios.App.BeforeTest()

	var userCountBefore, eventCountBefore, venueCountBefore, questionCountBefore, participationCountBefore, userQuestionCountBefore int
	var userLocal auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal})
	var userAdmin auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin})
	var userParticipant1 auth.User = auth.UserFactory(auth.User{Role: auth.UserRoleParticipant})
	var venue Venue = VenueFactorySaved(Venue{})
	var oldEvent Event = EventFactorySaved(Event{})
	var oldQuestions []Question = []Question{
		QuestionFactorySaved(Question{Event: &oldEvent}),
		QuestionFactorySaved(Question{Event: &oldEvent}),
		QuestionFactorySaved(Question{Event: &oldEvent}),
		QuestionFactorySaved(Question{Event: &oldEvent}),
	}
	var oldParticipations []Participation = []Participation{
		ParticipationFactorySaved(Participation{Event: &oldEvent, Venue: &venue, User: &userParticipant1}),
		ParticipationFactorySaved(Participation{Event: &oldEvent, Venue: &venue}),
		ParticipationFactorySaved(Participation{Event: &oldEvent, Venue: &venue}),
		ParticipationFactorySaved(Participation{Event: &oldEvent, Venue: &venue}),
	}
	UserQuestionFactorySaved(UserQuestion{Question: &oldQuestions[0], Participation: &oldParticipations[0]})
	helios.DB.Model(&auth.User{}).Count(&userCountBefore)
	helios.DB.Model(&Event{}).Count(&eventCountBefore)
	helios.DB.Model(&Venue{}).Count(&venueCountBefore)
	helios.DB.Model(&Question{}).Count(&questionCountBefore)
	helios.DB.Model(&Participation{}).Count(&participationCountBefore)
	helios.DB.Model(&UserQuestion{}).Count(&userQuestionCountBefore)
	type putSynchronizationDataTestCase struct {
		user                       auth.User
		event                      Event
		venue                      Venue
		questions                  []Question
		users                      []auth.User
		usersKey                   map[string]string
		expectedError              helios.Error
		expectedEventCount         int
		expectedVenueCount         int
		expectedUserCount          int
		expectedQuestionCount      int
		expectedParticipationCount int
		expectedUserQuestionCount  int
	}
	testCases := []putSynchronizationDataTestCase{{
		user:                       userAdmin,
		event:                      EventFactory(Event{}),
		venue:                      VenueFactory(Venue{}),
		questions:                  []Question{},
		users:                      []auth.User{},
		usersKey:                   make(map[string]string),
		expectedError:              errSynchronizationNotAuthorized,
		expectedUserCount:          userCountBefore,
		expectedVenueCount:         venueCountBefore,
		expectedEventCount:         eventCountBefore,
		expectedQuestionCount:      questionCountBefore,
		expectedParticipationCount: participationCountBefore,
		expectedUserQuestionCount:  userQuestionCountBefore,
	}, {
		user:                       userLocal,
		event:                      EventFactory(Event{}),
		venue:                      VenueFactory(Venue{}),
		questions:                  []Question{QuestionFactory(Question{})},
		users:                      []auth.User{auth.UserFactory(auth.User{Username: "user1", Role: auth.UserRoleParticipant})},
		usersKey:                   map[string]string{"user1": "key_user_1"},
		expectedUserCount:          userCountBefore + 1,
		expectedEventCount:         eventCountBefore + 1,
		expectedVenueCount:         venueCountBefore + 1,
		expectedQuestionCount:      questionCountBefore + 1,
		expectedParticipationCount: participationCountBefore + 1,
		expectedUserQuestionCount:  userQuestionCountBefore + 1,
	}, {
		user:                       userLocal,
		event:                      oldEvent,
		venue:                      VenueFactory(Venue{}),
		questions:                  []Question{QuestionFactory(Question{})},
		users:                      []auth.User{auth.UserFactory(auth.User{Username: "user2", Role: auth.UserRoleParticipant}), userParticipant1},
		usersKey:                   map[string]string{"user2": "key_user_2"},
		expectedUserCount:          userCountBefore + 2,
		expectedVenueCount:         venueCountBefore + 2,
		expectedEventCount:         eventCountBefore + 1,
		expectedQuestionCount:      questionCountBefore + 1 - len(oldQuestions) + 1,
		expectedParticipationCount: participationCountBefore + 1 - len(oldParticipations) + 2,
		expectedUserQuestionCount:  userQuestionCountBefore + 1 - 1 + 2,
	}}
	for i, testCase := range testCases {
		t.Logf("Test PutSynchronizationData testcase: %d", i)
		var err helios.Error
		var userCount, eventCount, venueCount, questionCount, participationCount, userQuestionCount int
		err = PutSynchronizationData(testCase.user, testCase.event, testCase.venue, testCase.questions, testCase.users, testCase.usersKey)
		helios.DB.Model(&auth.User{}).Count(&userCount)
		helios.DB.Model(&Event{}).Count(&eventCount)
		helios.DB.Model(&Venue{}).Count(&venueCount)
		helios.DB.Model(&Question{}).Count(&questionCount)
		helios.DB.Model(&Participation{}).Count(&participationCount)
		helios.DB.Model(&UserQuestion{}).Count(&userQuestionCount)
		assert.Equal(t, testCase.expectedUserCount, userCount)
		assert.Equal(t, testCase.expectedEventCount, eventCount)
		assert.Equal(t, testCase.expectedVenueCount, venueCount)
		assert.Equal(t, testCase.expectedQuestionCount, questionCount)
		assert.Equal(t, testCase.expectedParticipationCount, participationCount)
		assert.Equal(t, testCase.expectedUserQuestionCount, userQuestionCount)
		if testCase.expectedError == nil {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, testCase.expectedError, err)
		}
	}
}

func TestDecryptEventData(t *testing.T) {
	var userLocal auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal})
	var event3 Event = EventFactorySaved(Event{})
	var event2 Event = EventFactorySaved(Event{})
	var event1 Event = EventFactorySaved(Event{})
	var simKey string = event1.SimKey
	event1.DecryptedAt = time.Time{}
	event1.SimKey = ""
	event1.PrvKey = ""
	helios.DB.Save(&event1)
	var questions []Question = []Question{
		QuestionFactorySaved(Question{Event: &event1, Content: "content"}),
		QuestionFactorySaved(Question{Event: &event1, Content: "content"}),
		QuestionFactorySaved(Question{Event: &event1, Content: "content"}),
	}
	var err error = encryptQuestions(questions, simKey)
	assert.Nil(t, err)
	ParticipationFactorySaved(Participation{User: &userLocal, Event: &event1})
	ParticipationFactorySaved(Participation{User: &userLocal, Event: &event3})
	for _, question := range questions {
		t.Log(question.Content)
		helios.DB.Save(&question)
	}
	type decryptEventDataTestCase struct {
		user          auth.User
		eventSlug     string
		simKey        string
		expectedError helios.Error
	}
	testCases := []decryptEventDataTestCase{{
		user:          auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventSlug:     event1.Slug,
		simKey:        simKey,
		expectedError: errDecryptEventForbidden,
	}, {
		user:          userLocal,
		eventSlug:     event2.Slug,
		simKey:        simKey,
		expectedError: errEventNotFound,
	}, {
		user:          userLocal,
		eventSlug:     event1.Slug,
		simKey:        "wrong_key",
		expectedError: errDecryptEventFailed,
	}, {
		user:      userLocal,
		eventSlug: event1.Slug,
		simKey:    simKey,
	}, {
		user:      userLocal,
		eventSlug: event1.Slug,
		simKey:    simKey,
	}}
	for i, testCase := range testCases {
		t.Logf("Test DecryptEventData testcase: %d", i)
		var err helios.Error
		err = DecryptEventData(testCase.user, testCase.eventSlug, testCase.simKey)
		if testCase.expectedError == nil {
			var eventSaved Event
			helios.DB.Where("id = ?", event1.ID).First(&eventSaved)
			assert.Nil(t, err)
			assert.NotEmpty(t, eventSaved.DecryptedAt)
			assert.NotEmpty(t, eventSaved.SimKey)
			for _, question := range questions {
				var questionSaved Question
				helios.DB.Where("id = ?", question.ID).First(&questionSaved)
				assert.Equal(t, "content", questionSaved.Content)
			}
		} else {
			assert.Equal(t, testCase.expectedError, err)
		}
	}
}

func TestEncryption(t *testing.T) {
	type encryptionTestCase struct {
		plaintext  []byte
		ciphertext []byte // ciphertext encrypted by https://gchq.github.io/CyberChef/ to ensure there is no mgiac behind algorithm
		key        []byte
	}
	testCases := []encryptionTestCase{{
		plaintext: []byte("16 chars secret!"),
		ciphertext: []byte{
			0x2d, 0x98, 0x9d, 0xf2, 0x33, 0x90, 0xeb, 0xc1, 0x6c, 0x27, 0xf8, 0x71, 0xc2, 0x95, 0x1b, 0x48, // iv
			0x4c, 0x84, 0x23, 0xe0, 0x4e, 0x29, 0x68, 0x66, 0x0a, 0xd8, 0x56, 0x98, 0x9b, 0x35, 0x82, 0xc1, // ciphertext
		},
		key: []byte("32 characters super secret key!!"),
	}, {
		plaintext: []byte("secret text"),
		ciphertext: []byte{
			0xfa, 0x91, 0x76, 0x94, 0xaa, 0x73, 0x85, 0xf2, 0x37, 0xa3, 0xa2, 0x4a, 0xb6, 0xab, 0x7d, 0xc8, // iv
			0x98, 0xda, 0xce, 0xf2, 0xf6, 0x1c, 0xfe, 0xca, 0x9b, 0xdc, 0x08, // ciphertext
		},
		key: []byte("32 characters super secret key!!"),
	}, {
		plaintext: []byte("secret text longer than 16 character to ensure cfb mode"),
		ciphertext: []byte{
			0x91, 0x50, 0x71, 0xbb, 0xab, 0x57, 0xa6, 0xa3, 0x87, 0x92, 0xb0, 0x2c, 0x3b, 0xb9, 0xfc, 0xa0, // iv
			0xf7, 0x4d, 0x68, 0x73, 0xd9, 0x94, 0x51, 0x0f, 0x4a, 0xd8, 0xb2, 0x72, 0xef, 0x8d, 0xe1, 0x31, // ciphertext
			0xa6, 0x81, 0xd9, 0x9f, 0x09, 0xc7, 0xfe, 0x7c, 0x2d, 0x90, 0x16, 0x8d, 0xff, 0xf3, 0x94, 0xcc,
			0x85, 0x9d, 0xd3, 0x3c, 0x51, 0x88, 0xa6, 0xaf, 0x1c, 0xeb, 0x84, 0xae, 0x98, 0x95, 0x88, 0x41,
			0xbf, 0xd3, 0x9a, 0x52, 0xc6, 0x32, 0xca,
		},
		key: []byte("another 32 char super secret key"),
	}}
	for i, testCase := range testCases {
		t.Logf("Test encrypt/decrypt testcase: %d", i)
		var err error
		var encrypted, decrypted, result []byte
		encrypted, err = encrypt(testCase.key, testCase.plaintext)
		assert.Nil(t, err)
		decrypted, err = decrypt(testCase.key, encrypted)
		assert.Nil(t, err)
		assert.Equal(t, testCase.plaintext, decrypted)
		result, err = decrypt(testCase.key, testCase.ciphertext)
		assert.Nil(t, err)
		assert.Equal(t, testCase.plaintext, result, encrypted)
	}
}

func TestEncryptQuestions(t *testing.T) {
	var questions []Question = []Question{
		QuestionFactory(Question{Content: "content1", Choices: "choice1.1|choice1.2"}),
		QuestionFactory(Question{Content: "content2"}),
	}
	var key string = "32 characters super secret key!!"
	var err error
	err = encryptQuestions(questions, key)
	assert.Nil(t, err)
	assert.NotEqual(t, "content1", questions[0].Content)
	assert.NotEqual(t, "content2", questions[1].Content)
	assert.NotEqual(t, "choice1.1|choice1.2", questions[0].Choices)
	decryptQuestions(questions, key)
	assert.Nil(t, err)
	assert.Equal(t, "content1", questions[0].Content)
	assert.Equal(t, "content2", questions[1].Content)
	assert.Equal(t, "choice1.1|choice1.2", questions[0].Choices)
}
