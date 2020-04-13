package exam

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yonasadiel/helios"
)

func TestSerializeVenye(t *testing.T) {
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
		Title:       "Math Final Exam",
		Description: "desc",
		StartsAt:    time.Date(2020, 8, 12, 9, 30, 10, 0, time.FixedZone("Asia/Jakarta", int((7*time.Hour).Seconds()))),
		EndsAt:      time.Date(2020, 8, 12, 4, 30, 10, 0, time.FixedZone("UTC", 0)),
	})
	var expectedJSON string = `{"id":3,"title":"Math Final Exam","description":"desc","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"2020-08-12T11:30:10+07:00"}`
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
		eventDataJSON: `{"title":"Math Final Exam","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"2020-08-12T02:30:10Z","description":"desc"}`,
		expectedEvent: Event{
			ID:          0,
			Title:       "Math Final Exam",
			Description: "desc",
			StartsAt:    time.Date(2020, 8, 12, 9, 30, 10, 0, time.FixedZone("Asia/Jakarta", int((7*time.Hour).Seconds()))),
			EndsAt:      time.Date(2020, 8, 12, 2, 30, 10, 0, time.FixedZone("UTC", 0)),
		},
	}, {
		// endsAt is before startsAt
		eventDataJSON: `{"title":"Math Final Exam","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"2020-08-12T02:30:09Z"}`,
		expectedError: `{"code":"form_error","message":{"_error":[],"endsAt":["End time should be after start time"]}}`,
	}, {
		// wrong format endsAt and startsAt
		eventDataJSON: `{"title":"Math Final Exam","startsAt":"bad_format","endsAt":"bad_format"}`,
		expectedError: `{"code":"form_error","message":{"_error":[],"endsAt":["Failed to parse time"],"startsAt":["Failed to parse time"]}}`,
	}, {
		// empty fields
		eventDataJSON: `{}`,
		expectedError: `{"code":"form_error","message":{"_error":[],"endsAt":["End time must be provided"],"startsAt":["Start time must be provided"],"title":["Title can't be empty"]}}`,
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
