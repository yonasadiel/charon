package exam

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yonasadiel/charon/backend/auth"
	"github.com/yonasadiel/helios"
)

func TestSerializeVenue(t *testing.T) {
	var venue Venue = VenueFactory(Venue{
		ID:   3,
		Name: "venue name",
	})
	var expectedJSON string = `{"id":3,"name":"venue name"}`
	var serialized VenueData = SerializeVenue(venue)
	var serializedJSON []byte
	var errMarshalling error
	serializedJSON, errMarshalling = json.Marshal(serialized)
	assert.Nil(t, errMarshalling)
	assert.Equal(t, expectedJSON, string(serializedJSON))
}

func TestDeserializeVenue(t *testing.T) {
	type deserializeVenueTestCase struct {
		venueDataJSON string
		expectedVenue Venue
		expectedError string
	}
	testCases := []deserializeVenueTestCase{{
		venueDataJSON: `{"name":"Venue 1"}`,
		expectedVenue: Venue{
			ID:   0,
			Name: "Venue 1",
		},
	}, {
		venueDataJSON: `{"id":3,"name":"Venue 2"}`,
		expectedVenue: Venue{
			ID:   3,
			Name: "Venue 2",
		},
	}, {
		venueDataJSON: `{}`,
		expectedError: `{"code":"form_error","message":{"_error":[],"name":["Name can't be empty"]}}`,
	}}
	for i, testCase := range testCases {
		t.Logf("Test DeserializeVenue testcase: %d", i)
		var venue Venue
		var venueData VenueData
		var errUnmarshalling error
		var errDeserialization helios.Error
		errUnmarshalling = json.Unmarshal([]byte(testCase.venueDataJSON), &venueData)
		errDeserialization = DeserializeVenue(venueData, &venue)
		assert.Nil(t, errUnmarshalling)
		if testCase.expectedError == "" {
			assert.Nil(t, errDeserialization)
			assert.Equal(t, testCase.expectedVenue.ID, venue.ID, "Empty id on json will give 0")
			assert.Equal(t, testCase.expectedVenue.Name, venue.Name)
		} else {
			var errDeserializationJSON []byte
			var errMarshalling error
			errDeserializationJSON, errMarshalling = json.Marshal(errDeserialization.GetMessage())
			assert.Nil(t, errMarshalling)
			assert.NotNil(t, errDeserialization)
			assert.Equal(t, testCase.expectedError, string(errDeserializationJSON))
		}
	}
}

func TestSerializeEvent(t *testing.T) {
	var event Event = EventFactory(Event{
		ID:          3,
		Slug:        "math-final-exam",
		Title:       "Math Final Exam",
		Description: "desc",
		StartsAt:    time.Date(2020, 8, 12, 9, 30, 10, 0, time.FixedZone("Asia/Jakarta", int((7*time.Hour).Seconds()))),
		EndsAt:      time.Date(2020, 8, 12, 4, 30, 10, 0, time.FixedZone("UTC", 0)),
	})
	var expectedJSON string = `{"id":3,"slug":"math-final-exam","title":"Math Final Exam","description":"desc","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"2020-08-12T11:30:10+07:00"}`
	var serialized EventData = SerializeEvent(event)
	var serializedJSON []byte
	var errMarshalling error
	serializedJSON, errMarshalling = json.Marshal(serialized)
	assert.Nil(t, errMarshalling)
	assert.Equal(t, expectedJSON, string(serializedJSON))
}

func TestDeserializeEvent(t *testing.T) {
	type deserializeEventTestCase struct {
		eventDataJSON string
		expectedEvent Event
		expectedError string
	}
	testCases := []deserializeEventTestCase{{
		eventDataJSON: `{"title":"Math Final Exam","slug":"math-final-exam","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"2020-08-12T02:30:10Z","description":"desc"}`,
		expectedEvent: Event{
			ID:          0,
			Slug:        "math-final-exam",
			Title:       "Math Final Exam",
			Description: "desc",
			StartsAt:    time.Date(2020, 8, 12, 9, 30, 10, 0, time.FixedZone("Asia/Jakarta", int((7*time.Hour).Seconds()))),
			EndsAt:      time.Date(2020, 8, 12, 2, 30, 10, 0, time.FixedZone("UTC", 0)),
		},
	}, {
		// endsAt is before startsAt
		eventDataJSON: `{"title":"Math Final Exam","slug":"math-final-exam","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"2020-08-12T02:30:09Z"}`,
		expectedError: `{"code":"form_error","message":{"_error":[],"endsAt":["End time should be after start time"]}}`,
	}, {
		// wrong format endsAt and startsAt
		eventDataJSON: `{"title":"Math Final Exam","slug":"math-final-exam","startsAt":"bad_format","endsAt":"bad_format"}`,
		expectedError: `{"code":"form_error","message":{"_error":[],"endsAt":["Failed to parse time"],"startsAt":["Failed to parse time"]}}`,
	}, {
		// empty fields
		eventDataJSON: `{}`,
		expectedError: `{"code":"form_error","message":{"_error":[],"endsAt":["End time must be provided"],"slug":["Slug can't be empty"],"startsAt":["Start time must be provided"],"title":["Title can't be empty"]}}`,
	}}
	for i, testCase := range testCases {
		t.Logf("Test DeserializeEvent testcase: %d", i)
		var event Event
		var eventData EventData
		var errUnmarshalling error
		var errDeserialization helios.Error
		errUnmarshalling = json.Unmarshal([]byte(testCase.eventDataJSON), &eventData)
		errDeserialization = DeserializeEvent(eventData, &event)
		assert.Nil(t, errUnmarshalling)
		if testCase.expectedError == "" {
			assert.Nil(t, errDeserialization)
			assert.Equal(t, testCase.expectedEvent.ID, event.ID, "Empty id on json will give 0")
			assert.Equal(t, testCase.expectedEvent.Title, event.Title)
			assert.Equal(t, testCase.expectedEvent.Description, event.Description)
			assert.True(t, testCase.expectedEvent.StartsAt.Equal(event.StartsAt))
			assert.True(t, testCase.expectedEvent.EndsAt.Equal(event.EndsAt))
		} else {
			var errDeserializationJSON []byte
			var errMarshalling error
			errDeserializationJSON, errMarshalling = json.Marshal(errDeserialization.GetMessage())
			assert.Nil(t, errMarshalling)
			assert.NotNil(t, errDeserialization)
			assert.Equal(t, testCase.expectedError, string(errDeserializationJSON))
		}
	}
}

func TestSerializeParticipation(t *testing.T) {
	var user auth.User = auth.UserFactory(auth.User{Username: "abc"})
	var venue Venue = VenueFactory(Venue{ID: 5})
	var participation Participation = ParticipationFactory(Participation{
		ID:    3,
		User:  &user,
		Venue: &venue,
	})
	var expectedJSON string = `{"id":3,"userUsername":"abc","venueId":5}`
	var serialized ParticipationData = SerializeParticipation(participation)
	var serializedJSON []byte
	var errMarshalling error
	serializedJSON, errMarshalling = json.Marshal(serialized)
	assert.Nil(t, errMarshalling)
	assert.Equal(t, expectedJSON, string(serializedJSON))
}

func TestDeserializeParticipation(t *testing.T) {
	type deserializeParticipationTestCase struct {
		participationDataJSON string
		expectedParticipation Participation
		expectedError         string
	}
	testCases := []deserializeParticipationTestCase{{
		participationDataJSON: `{"id":2,"eventId":3,"venueId":4,"userId":5,"userUsername":"abc"}`,
		expectedParticipation: Participation{
			ID:      2,
			VenueID: 4,
		},
	}, {
		participationDataJSON: `{"venueId":4,"userUsername":"abc"}`,
		expectedParticipation: Participation{
			VenueID: 4,
		},
	}, {
		participationDataJSON: `{}`,
		expectedError:         `{"code":"form_error","message":{"_error":[],"userUsername":["Username can't be empty"],"venueId":["Venue can't be empty"]}}`,
	}}
	for i, testCase := range testCases {
		t.Logf("Test DeserializeParticipation testcase: %d", i)
		var participationData ParticipationData
		var participation Participation
		var errUnmarshalling error
		var errDeserialization helios.Error
		errUnmarshalling = json.Unmarshal([]byte(testCase.participationDataJSON), &participationData)
		errDeserialization = DeserializeParticipation(participationData, &participation)
		assert.Nil(t, errUnmarshalling)
		if testCase.expectedError == "" {
			assert.Nil(t, errDeserialization)
			assert.Equal(t, testCase.expectedParticipation.ID, participation.ID)
			assert.Equal(t, testCase.expectedParticipation.VenueID, participation.VenueID)
			assert.Equal(t, testCase.expectedParticipation.EventID, participation.EventID)
			assert.Equal(t, testCase.expectedParticipation.UserID, participation.UserID)
			assert.Nil(t, participation.Event)
			assert.Nil(t, participation.User)
			assert.Nil(t, participation.Venue)
		} else {
			var errDeserializationJSON []byte
			var errMarshalling error
			errDeserializationJSON, errMarshalling = json.Marshal(errDeserialization.GetMessage())
			assert.Nil(t, errMarshalling)
			assert.NotNil(t, errDeserialization)
			assert.Equal(t, testCase.expectedError, string(errDeserializationJSON))
		}
	}
}

func TestSerializeQuestion(t *testing.T) {
	type serializeQuestionTestCase struct {
		question     Question
		expectedJSON string
	}
	testCases := []serializeQuestionTestCase{{
		question: Question{
			ID:      2,
			Content: "Question Content",
			Choices: []QuestionChoice{},
		},
		expectedJSON: `{"id":2,"content":"Question Content","choices":[],"answer":""}`,
	}, {
		question: Question{
			ID:         2,
			Content:    "Question Content",
			Choices:    []QuestionChoice{{Text: "a"}, {Text: "b"}, {Text: "c"}},
			UserAnswer: "answer2",
		},
		expectedJSON: `{"id":2,"content":"Question Content","choices":["a","b","c"],"answer":"answer2"}`,
	}}
	for i, testCase := range testCases {
		t.Logf("Test SerializeQuestion testcase: %d", i)
		var serialized QuestionData
		var serializedJSON []byte
		var errMarshalling error
		serialized = SerializeQuestion(testCase.question)
		serializedJSON, errMarshalling = json.Marshal(serialized)
		assert.Nil(t, errMarshalling)
		assert.Equal(t, testCase.expectedJSON, string(serializedJSON))
	}
}

func TestDeserializeQuestion(t *testing.T) {
	type deserializeQuestionTestCase struct {
		questionDataJSON string
		expectedQuestion Question
		expectedError    string
	}
	testCases := []deserializeQuestionTestCase{{
		questionDataJSON: `{"id":2,"content":"Question Content","choices":[],"answer":""}`,
		expectedQuestion: Question{
			ID:      2,
			Content: "Question Content",
			Choices: []QuestionChoice(nil),
		},
	}, {
		questionDataJSON: `{"id":2,"content":"Question Content","choices":["a","","b","","c", ""],"answer":""}`,
		expectedQuestion: Question{
			ID:      2,
			Content: "Question Content",
			Choices: []QuestionChoice{{Text: "a"}, {Text: "b"}, {Text: "c"}},
		},
	}, {
		questionDataJSON: `{"id":2,"content":"","choices":[],"answer":""}`,
		expectedError:    `{"code":"form_error","message":{"_error":[],"content":["Content can't be empty"]}}`,
	}}
	for i, testCase := range testCases {
		t.Logf("Test DeserializeQuestion testcase: %d", i)
		var questionData QuestionData
		var question Question
		var errUnmarshalling error
		var errDeserialization helios.Error
		errUnmarshalling = json.Unmarshal([]byte(testCase.questionDataJSON), &questionData)
		errDeserialization = DeserializeQuestion(questionData, &question)
		assert.Nil(t, errUnmarshalling)
		if testCase.expectedError == "" {
			assert.Nil(t, errDeserialization)
			assert.Equal(t, testCase.expectedQuestion.ID, question.ID)
			assert.Equal(t, testCase.expectedQuestion.Content, question.Content)
			assert.Equal(t, len(testCase.expectedQuestion.Choices), len(question.Choices))
			for i := range testCase.expectedQuestion.Choices {
				assert.Equal(t, testCase.expectedQuestion.Choices[i].Text, question.Choices[i].Text)
			}
		} else {
			var errDeserializationJSON []byte
			var errMarshalling error
			errDeserializationJSON, errMarshalling = json.Marshal(errDeserialization.GetMessage())
			assert.Nil(t, errMarshalling)
			assert.NotNil(t, errDeserialization)
			assert.Equal(t, testCase.expectedError, string(errDeserializationJSON))
		}
	}
}

func TestSerializeSynchronizationData(t *testing.T) {
	type serializeSynchronizationDataTestCase struct {
		event          Event
		questions      []Question
		participations []Participation
		users          []auth.User
		expectedJSON   string
	}
	testCases := []serializeSynchronizationDataTestCase{{
		event: Event{
			ID:          3,
			Slug:        "math-final-exam",
			Title:       "Math Final Exam",
			Description: "desc",
			StartsAt:    time.Date(2020, 8, 12, 9, 30, 10, 0, time.FixedZone("Asia/Jakarta", int((7*time.Hour).Seconds()))),
			EndsAt:      time.Date(2020, 8, 12, 4, 30, 10, 0, time.FixedZone("UTC", 0)),
		},
		questions: []Question{{
			ID:         2,
			Content:    "Question Content",
			Choices:    []QuestionChoice{{Text: "a"}, {Text: "b"}, {Text: "c"}},
			UserAnswer: "answer2",
		}, {}},
		participations: []Participation{{
			ID:    3,
			User:  &auth.User{Username: "abc"},
			Venue: &Venue{ID: 3},
		}},
		users: []auth.User{{
			ID:       4,
			Username: "def",
			Role:     auth.UserRoleAdmin,
			Name:     "abc",
		}},
		expectedJSON: `{` +
			`"event":{"id":3,"slug":"math-final-exam","title":"Math Final Exam","description":"desc","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"2020-08-12T11:30:10+07:00"},` +
			`"questions":[{"id":2,"content":"Question Content","choices":["a","b","c"],"answer":"answer2"},{"id":0,"content":"","choices":[],"answer":""}],` +
			`"participations":[{"id":3,"userUsername":"abc","venueId":3}],` +
			`"users":[{"name":"abc","username":"def","role":"admin"}]` +
			`}`,
	}, {
		event:          Event{},
		questions:      []Question{},
		participations: []Participation{},
		users:          []auth.User{},
		expectedJSON: `{` +
			`"event":{"id":0,"slug":"","title":"","description":"","startsAt":"0001-01-01T07:07:12+07:07","endsAt":"0001-01-01T07:07:12+07:07"},` +
			`"questions":[],` +
			`"participations":[],` +
			`"users":[]` +
			`}`,
	}}
	for i, testCase := range testCases {
		t.Logf("Test SerializeSynchronizationData testcase: %d", i)
		var serialized SynchronizationData
		var serializedJSON []byte
		var errMarshalling error
		serialized = SerializeSynchronizationData(testCase.event, testCase.questions, testCase.participations, testCase.users)
		serializedJSON, errMarshalling = json.Marshal(serialized)
		assert.Nil(t, errMarshalling)
		assert.Equal(t, testCase.expectedJSON, string(serializedJSON))
	}
}

func TestDeserializeSynchronizationData(t *testing.T) {
	type deserializeQuestionTestCase struct {
		synchronizationDataJSON     string
		expectedEvent               Event
		expectedQuestionLength      int
		expectedParticipationLength int
		expectedUserLength          int
		expectedError               string
	}
	testCases := []deserializeQuestionTestCase{{
		synchronizationDataJSON: `{` +
			`"event":{"id":3,"slug":"math-final-exam","title":"Math Final Exam","description":"desc","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"2020-08-12T11:30:10+07:00"},` +
			`"questions":[{"id":2,"content":"Question Content","choices":["a","b","c"],"answer":"answer2"},{"id":0,"content":"a","choices":[],"answer":""}],` +
			`"participations":[{"id":3,"userUsername":"abc","venueId":3},{"id":3,"userUsername":"abc","venueId":3},{"id":3,"userUsername":"abc","venueId":3}],` +
			`"users":[{"name":"abc","username":"def","role":"admin"}]` +
			`}`,
		expectedEvent: Event{
			ID:          3,
			Slug:        "math-final-exam",
			Title:       "Math Final Exam",
			Description: "desc",
			StartsAt:    time.Date(2020, 8, 12, 9, 30, 10, 0, time.FixedZone("Asia/Jakarta", int((7*time.Hour).Seconds()))),
			EndsAt:      time.Date(2020, 8, 12, 4, 30, 10, 0, time.FixedZone("UTC", 0)),
		},
		expectedQuestionLength:      2,
		expectedParticipationLength: 3,
		expectedUserLength:          1,
	}, {
		synchronizationDataJSON: `{"event":{"endsAt":"2020-08-12T11:30:10+07:00","startsAt":"2020-08-12T09:30:10+07:00","title":"abc","slug":"abc"}}`,
		expectedEvent: Event{
			Title:    "abc",
			Slug:     "abc",
			StartsAt: time.Date(2020, 8, 12, 9, 30, 10, 0, time.FixedZone("Asia/Jakarta", int((7*time.Hour).Seconds()))),
			EndsAt:   time.Date(2020, 8, 12, 4, 30, 10, 0, time.FixedZone("UTC", 0)),
		},
		expectedQuestionLength:      0,
		expectedParticipationLength: 0,
		expectedUserLength:          0,
	}, {
		synchronizationDataJSON: `{` +
			`"event":{"endsAt":"2020-08-12T01:30:10+07:00","startsAt":"2020-08-12T09:30:10+07:00","title":"abc","slug":"abc"},` +
			`"questions":[{}],` +
			`"participations":[{"venueId":2},{"venueId":1,"userUsername":"a"}],` +
			`"users":[{"name":"abc","role":"admin","username":"abc"},{"role":"admin","username":"abc"}]` +
			`}`,
		expectedError: `{"code":"form_error","message":{` +
			`"_error":[],` +
			`"event":{"endsAt":["End time should be after start time"]},` +
			`"participations":[{"userUsername":["Username can't be empty"]},{}],` +
			`"questions":[{"content":["Content can't be empty"]}],` +
			`"users":[{},{"name":["Name can't be empty"]}]` +
			`}}`,
	}}
	for i, testCase := range testCases {
		t.Logf("Test DeserializeSynchronizationData testcase: %d", i)
		var synchronizationData SynchronizationData
		var event Event
		var questions []Question
		var participations []Participation
		var users []auth.User
		var errUnmarshalling error
		var errDeserialization helios.Error
		errUnmarshalling = json.Unmarshal([]byte(testCase.synchronizationDataJSON), &synchronizationData)
		errDeserialization = DeserializeSynchronizationData(synchronizationData, &event, &questions, &participations, &users)
		assert.Nil(t, errUnmarshalling)
		if testCase.expectedError == "" {
			assert.Nil(t, errDeserialization)
			assert.Equal(t, testCase.expectedEvent.ID, event.ID)
			assert.Equal(t, testCase.expectedEvent.Title, event.Title)
			assert.Equal(t, testCase.expectedEvent.Description, event.Description)
			assert.True(t, testCase.expectedEvent.StartsAt.Equal(event.StartsAt))
			assert.True(t, testCase.expectedEvent.EndsAt.Equal(event.EndsAt))
		} else {
			var errDeserializationJSON []byte
			var errMarshalling error
			errDeserializationJSON, errMarshalling = json.Marshal(errDeserialization.GetMessage())
			assert.Nil(t, errMarshalling)
			assert.NotNil(t, errDeserialization)
			assert.Equal(t, testCase.expectedError, string(errDeserializationJSON))
		}
	}
}
