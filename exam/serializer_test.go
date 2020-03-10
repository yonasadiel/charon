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
	expectedJSON := `{"id":3,"title":"Math Final Exam","starts_at":"2020-08-12T09:30:10+07:00","ends_at":"2020-08-12T04:30:10Z"}`
	ser := SerializeEvent(event)
	serJSON, err := json.Marshal(ser)
	if err != nil {
		t.Errorf("Error marshaling json: %s", err)
	}
	assert.Equal(t, expectedJSON, string(serJSON), "Unequal JSON")
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
	var questionData QuestionResponse
	err := json.Unmarshal([]byte(originalJSON), &questionData)
	if err != nil {
		t.Errorf("Error marshaling json: %s", err)
	}
	question := DeserializeQuestion(questionData)

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
	var questionData QuestionResponse
	err := json.Unmarshal([]byte(originalJSON), &questionData)
	if err != nil {
		t.Errorf("Error marshaling json: %s", err)
	}
	question := DeserializeQuestion(questionData)

	assert.Equal(t, expectedQuestion.ID, question.ID, "Unequal question ID")
	assert.Equal(t, expectedQuestion.Content, question.Content, "Unequal question content")
	assert.Equal(t, len(expectedQuestion.Choices), len(question.Choices), "Unequal question num of choices")
	for i := range expectedQuestion.Choices {
		assert.Equal(t, expectedQuestion.Choices[i].Text, question.Choices[i].Text, "Unequal question num of choices")
	}
}
