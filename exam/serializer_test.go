package exam

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSerializeEvent(t *testing.T) {
	beforeTest(false)

	jakartaSeconds := int((7 * time.Hour).Seconds())
	jakartaTZ := time.FixedZone("Asia/Jakarta", jakartaSeconds)
	utcTZ := time.FixedZone("UTC", 0)

	event := Event{
		ID:       3,
		Title:    "Math Final Exam",
		StartsAt: time.Date(2020, 8, 12, 9, 30, 10, 0, jakartaTZ),
		EndsAt:   time.Date(2020, 8, 12, 4, 30, 10, 0, utcTZ),
	}
	expectedJSON := `{"id":3,"title":"Math Final Exam","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"2020-08-12T04:30:10Z"}`
	ser := SerializeEvent(event)
	serJSON, err := json.Marshal(ser)
	if err != nil {
		t.Errorf("Error marshaling json: %s", err)
	}
	assert.Equal(t, expectedJSON, string(serJSON), "Unequal JSON")
}

func TestDeserializeEvent(t *testing.T) {
	beforeTest(false)

	jakartaSeconds := int((7 * time.Hour).Seconds())
	jakartaTZ := time.FixedZone("Asia/Jakarta", jakartaSeconds)
	utcTZ := time.FixedZone("UTC", 0)

	var eventData EventData
	var event Event
	expectedEvent := Event{
		ID:       0,
		Title:    "Math Final Exam",
		StartsAt: time.Date(2020, 8, 12, 9, 30, 10, 0, jakartaTZ),
		EndsAt:   time.Date(2020, 8, 12, 2, 30, 10, 0, utcTZ),
	}
	json1 := `{"title":"Math Final Exam","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"2020-08-12T02:30:10Z"}`
	err1 := json.Unmarshal([]byte(json1), &eventData)
	if err1 != nil {
		t.Errorf("Error unmarshaling json: %s", err1)
	}
	errDeserialization1 := DeserializeEvent(eventData, &event)
	assert.Nil(t, errDeserialization1, "Failed to deserialize event")
	assert.Equal(t, uint(0), event.ID, "Empty id on json will give 0")
	assert.Equal(t, expectedEvent.Title, event.Title, "Unequal event title")
	assert.True(t, expectedEvent.StartsAt.Equal(event.StartsAt), "Unequal event start time")
	assert.True(t, expectedEvent.EndsAt.Equal(event.EndsAt), "Unequal Event end time")

	json2 := `{"title":"Math Final Exam","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"2020-08-12T02:30:09Z"}`
	eventData = EventData{}
	err2 := json.Unmarshal([]byte(json2), &eventData)
	if err2 != nil {
		t.Errorf("Error unmarshaling json: %s", err2)
	}
	errDeserialization2 := DeserializeEvent(eventData, &event)
	expectedError2 := `{"code":"form_error","message":{"_error":[],"endsAt":["End time should be after start time"]}}`
	errMessage2, _ := json.Marshal(errDeserialization2.GetMessage())
	assert.Equal(t, expectedError2, string(errMessage2), "endsAt is before startsAt, should give error")

	json3 := `{}`
	eventData = EventData{}
	err3 := json.Unmarshal([]byte(json3), &eventData)
	if err3 != nil {
		t.Errorf("Error unmarshaling json: %s", err3)
	}
	errDeserialization4 := DeserializeEvent(eventData, &event)
	expectedError4 := `{"code":"form_error","message":{"_error":[],"endsAt":["End time must be provided"],"startsAt":["Start time must be provided"],"title":["Title can't be empty"]}}`
	errMessage4, _ := json.Marshal(errDeserialization4.GetMessage())
	assert.Equal(t, expectedError4, string(errMessage4), "Error of empty json, different error message")
}

func TestSerializeQuestionEmptyChoices(t *testing.T) {
	beforeTest(false)

	question := Question{
		ID:      2,
		Content: "Question Content",
		Choices: []QuestionChoice{},
	}

	expectedJSON := `{"id":2,"content":"Question Content","choices":[],"answer":""}`
	ser := SerializeQuestion(question)
	serJSON, err := json.Marshal(ser)
	if err != nil {
		t.Errorf("Error marshaling json: %s", err)
	}

	assert.Equal(t, expectedJSON, string(serJSON), "Unequal JSON")
}

func TestSerializeQuestionWithChoices(t *testing.T) {
	beforeTest(false)

	question := Question{
		ID:      2,
		Content: "Question Content",
		Choices: []QuestionChoice{
			QuestionChoice{Text: "a"},
			QuestionChoice{Text: "b"},
			QuestionChoice{Text: "c"},
		},
		UserAnswer: "answer2",
	}

	expectedJSON := `{"id":2,"content":"Question Content","choices":["a","b","c"],"answer":"answer2"}`
	ser := SerializeQuestion(question)
	serJSON, err := json.Marshal(ser)
	if err != nil {
		t.Errorf("Error marshaling json: %s", err)
	}

	assert.Equal(t, expectedJSON, string(serJSON), "Unequal JSON")
}

func TestDeserializerQuestionEmptyChoices(t *testing.T) {
	beforeTest(false)

	expectedQuestion := Question{
		ID:      2,
		Content: "Question Content",
		Choices: []QuestionChoice(nil),
	}
	originalJSON := `{"id":2,"content":"Question Content","choices":[],"answer":""}`
	var questionData QuestionData
	var question Question
	err := json.Unmarshal([]byte(originalJSON), &questionData)
	if err != nil {
		t.Errorf("Error marshaling json: %s", err)
	}
	errDeserializeQuestion := DeserializeQuestion(questionData, &question)
	assert.Nil(t, errDeserializeQuestion, "Failed to deserialize question")
	assert.Equal(t, expectedQuestion.ID, question.ID, "Unequal question ID")
	assert.Equal(t, expectedQuestion.Content, question.Content, "Unequal question content")
	assert.Equal(t, len(expectedQuestion.Choices), len(question.Choices), "Unequal question num of choices")
}

func TestDeserializerQuestionWithChoices(t *testing.T) {
	beforeTest(false)

	expectedQuestion := Question{
		ID:      2,
		Content: "Question Content",
		Choices: []QuestionChoice{
			QuestionChoice{Text: "a"},
			QuestionChoice{Text: "b"},
			QuestionChoice{Text: "c"},
		},
	}
	originalJSON := `{"id":2,"content":"Question Content","choices":["a","b","c"],"answer":""}`
	var questionData QuestionData
	var question Question
	err := json.Unmarshal([]byte(originalJSON), &questionData)
	if err != nil {
		t.Errorf("Error marshaling json: %s", err)
	}
	errDeserializeQuestion := DeserializeQuestion(questionData, &question)

	assert.Nil(t, errDeserializeQuestion, "Error deserializting question")
	assert.Equal(t, expectedQuestion.ID, question.ID, "Unequal question ID")
	assert.Equal(t, expectedQuestion.Content, question.Content, "Unequal question content")
	assert.Equal(t, len(expectedQuestion.Choices), len(question.Choices), "Unequal question num of choices")
	for i := range expectedQuestion.Choices {
		assert.Equal(t, expectedQuestion.Choices[i].Text, question.Choices[i].Text, "Unequal question num of choices")
	}
}
