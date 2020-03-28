package exam

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yonasadiel/charon/auth"
	"github.com/yonasadiel/helios"
)

func TestEventListView(t *testing.T) {
	beforeTest(true)

	req := helios.NewMockRequest()
	req.SetContextData(auth.UserContextKey, user1)
	EventListView(&req)
	assert.Equal(t, http.StatusOK, req.StatusCode, "Unexpected status code")

	req.SetContextData(auth.UserContextKey, "bad_user")
	EventListView(&req)
	assert.Equal(t, http.StatusInternalServerError, req.StatusCode, "Unexpected status code")
}

func TestEventCreateView(t *testing.T) {
	beforeTest(false)
	var eventCountBefore int
	var eventCount int
	helios.DB.Model(Event{}).Count(&eventCountBefore)

	type eventCreateTestCase struct {
		user               interface{}
		requestData        string
		expectedStatusCode int
		expectedErrorCode  string
		expectedEventCount int
	}
	testCases := []eventCreateTestCase{
		eventCreateTestCase{
			user:               userOrganizer,
			requestData:        `{"title":"Math Final Exam","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"2020-08-12T04:30:10Z"}`,
			expectedStatusCode: http.StatusCreated,
			expectedEventCount: eventCountBefore + 1,
		},
		eventCreateTestCase{
			user:               userOrganizer,
			requestData:        `{"title":"Math Final Exam","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"INVALID_END_TIME"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  helios.ErrJSONParseFailed.Code,
			expectedEventCount: eventCountBefore + 1,
		},
		eventCreateTestCase{
			user:               userOrganizer,
			requestData:        `{"title":"Math Final Exam","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"2020-08-12T02:30:09Z"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  "form_error",
			expectedEventCount: eventCountBefore + 1,
		},
		eventCreateTestCase{
			user:               "bad_user",
			requestData:        `{"title":"Math Final Exam","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"2020-08-12T04:30:10Z"}`,
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  helios.ErrInternalServerError.Code,
			expectedEventCount: eventCountBefore + 1,
		},
	}
	for i, testCase := range testCases {
		t.Logf("Test EventCreate testcase: %d", i)
		req := helios.NewMockRequest()
		req.SetContextData(auth.UserContextKey, testCase.user)
		req.RequestData = testCase.requestData

		EventCreateView(&req)

		helios.DB.Model(Event{}).Count(&eventCount)
		assert.Equal(t, testCase.expectedStatusCode, req.StatusCode, "Unexpected status code")
		assert.Equal(t, testCase.expectedEventCount, eventCount, "Unexpected event count")
		if testCase.expectedErrorCode != "" {
			var err map[string]interface{}
			json.Unmarshal(req.JSONResponse, &err)
			assert.Equal(t, testCase.expectedErrorCode, err["code"], "Different error code")
		}
	}
}

func TestQuestionListView(t *testing.T) {
	beforeTest(true)

	type questionListTestCase struct {
		user               interface{}
		eventID            string
		expectedStatusCode int
		expectedErrorCode  string
	}
	testCases := []questionListTestCase{
		questionListTestCase{
			user:               user1,
			eventID:            strconv.Itoa(int(event1.ID)),
			expectedStatusCode: http.StatusOK,
		},
		questionListTestCase{
			user:               user1,
			eventID:            "malformed",
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  errEventNotFound.Code,
		},
		questionListTestCase{
			user:               user1,
			eventID:            "79697",
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  errEventNotFound.Code,
		},
		questionListTestCase{
			user:               "bad_user",
			eventID:            strconv.Itoa(int(event1.ID)),
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  helios.ErrInternalServerError.Code,
		},
	}
	for i, testCase := range testCases {
		t.Logf("Test QuestionList testcase: %d", i)
		req := helios.NewMockRequest()
		req.SetContextData(auth.UserContextKey, testCase.user)
		req.URLParam["eventID"] = testCase.eventID

		QuestionListView(&req)

		assert.Equal(t, testCase.expectedStatusCode, req.StatusCode, "Unexpected status code")
		if testCase.expectedErrorCode != "" {
			var err map[string]interface{}
			json.Unmarshal(req.JSONResponse, &err)
			assert.Equal(t, testCase.expectedErrorCode, err["code"], "Different error code")
		}
	}
}

func TestQuestionCreateView(t *testing.T) {
	beforeTest(true)

	var questionCountBefore int
	var questionCount int
	helios.DB.Model(Question{}).Count(&questionCountBefore)

	type questionCreateTestCase struct {
		user                  interface{}
		eventID               string
		requestData           string
		expectedQuestionCount int
		expectedStatusCode    int
		expectedErrorCode     string
	}

	testCases := []questionCreateTestCase{
		questionCreateTestCase{
			user:                  userOrganizer,
			eventID:               strconv.Itoa(int(event1.ID)),
			requestData:           `{"id":1,"content":"content1","choices":[],"answer":"abc","eventId":2}`,
			expectedQuestionCount: questionCountBefore + 1,
			expectedStatusCode:    http.StatusCreated,
		},
		questionCreateTestCase{
			user:                  userOrganizer,
			eventID:               "abcdef",
			requestData:           `{"id":1,"content":"content2","choices":[],"answer":"abc","eventId":2}`,
			expectedQuestionCount: questionCountBefore + 1,
			expectedStatusCode:    http.StatusNotFound,
			expectedErrorCode:     errEventNotFound.Code,
		},
		questionCreateTestCase{
			user:                  userAdmin,
			eventID:               "999",
			requestData:           `{"id":1,"content":"content2","choices":[],"answer":"abc","eventId":2}`,
			expectedQuestionCount: questionCountBefore + 1,
			expectedStatusCode:    http.StatusNotFound,
			expectedErrorCode:     errEventNotFound.Code,
		},
		questionCreateTestCase{
			user:                  userParticipant,
			eventID:               strconv.Itoa(int(event1.ID)),
			requestData:           `{"id":1,"content":"content3","choices":[],"answer":"abc","eventId":2}`,
			expectedQuestionCount: questionCountBefore + 1,
			expectedStatusCode:    http.StatusUnauthorized,
			expectedErrorCode:     errQuestionChangeNotAuthorized.Code,
		},
		questionCreateTestCase{
			user:                  userLocal,
			eventID:               strconv.Itoa(int(event1.ID)),
			requestData:           `{"id":1,"content":"content3","choices":[],"answer":"abc","eventId":2}`,
			expectedQuestionCount: questionCountBefore + 1,
			expectedStatusCode:    http.StatusUnauthorized,
			expectedErrorCode:     errQuestionChangeNotAuthorized.Code,
		},
		questionCreateTestCase{
			user:                  userOrganizer,
			eventID:               strconv.Itoa(int(event1.ID)),
			requestData:           `bad_format`,
			expectedQuestionCount: questionCountBefore + 1,
			expectedStatusCode:    http.StatusBadRequest,
		},
		questionCreateTestCase{
			user:                  userOrganizer,
			eventID:               strconv.Itoa(int(event1.ID)),
			requestData:           `{"content":"content5","choices":["content5_1","content5_2"],"answer":"abc","eventId":2}`,
			expectedQuestionCount: questionCountBefore + 2,
			expectedStatusCode:    http.StatusCreated,
		},
		questionCreateTestCase{
			user:                  "bad_user",
			eventID:               strconv.Itoa(int(event1.ID)),
			requestData:           `{"content":"content6","choices":[],"answer":"abc","eventId":2}`,
			expectedQuestionCount: questionCountBefore + 2,
			expectedStatusCode:    http.StatusInternalServerError,
			expectedErrorCode:     helios.ErrInternalServerError.Code,
		},
	}

	for i, testCase := range testCases {
		t.Logf("Test QuestionCreate testcase: %d", i)
		req := helios.NewMockRequest()
		req.SetContextData(auth.UserContextKey, testCase.user)
		req.URLParam["eventID"] = testCase.eventID
		req.RequestData = testCase.requestData

		QuestionCreateView(&req)

		helios.DB.Model(Question{}).Count(&questionCount)
		assert.Equal(t, testCase.expectedQuestionCount, questionCount, "Different question count")
		assert.Equal(t, testCase.expectedStatusCode, req.StatusCode, "Unexpected status code")
		if testCase.expectedErrorCode != "" {
			var err map[string]interface{}
			json.Unmarshal(req.JSONResponse, &err)
			assert.Equal(t, testCase.expectedErrorCode, err["code"], "Different error code")
		}
	}
}

func TestQuestionDetailView(t *testing.T) {
	beforeTest(true)

	type questionDetailTestCase struct {
		user               interface{}
		eventID            string
		questionID         string
		expectedStatusCode int
		expectedErrorCode  string
	}

	testCases := []questionDetailTestCase{
		questionDetailTestCase{
			user:               user1,
			eventID:            strconv.Itoa(int(event1.ID)),
			questionID:         strconv.Itoa(int(questionSimple.ID)),
			expectedStatusCode: http.StatusOK,
		},
		questionDetailTestCase{
			user:               user1,
			eventID:            strconv.Itoa(int(event1.ID)),
			questionID:         "bad_question_id",
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  errQuestionNotFound.Code,
		},
		questionDetailTestCase{
			user:               user1,
			eventID:            strconv.Itoa(int(event1.ID)),
			questionID:         "879654",
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  errQuestionNotFound.Code,
		},
		questionDetailTestCase{
			user:               user1,
			eventID:            "4567890",
			questionID:         "879654",
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  errEventNotFound.Code,
		},
		questionDetailTestCase{
			user:               user1,
			eventID:            "malformed",
			questionID:         "malformed",
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  errEventNotFound.Code,
		},
		questionDetailTestCase{
			user:               "bad_user",
			eventID:            strconv.Itoa(int(event1.ID)),
			questionID:         strconv.Itoa(int(questionSimple.ID)),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for i, testCase := range testCases {
		t.Logf("Test QuestionDetail testcase: %d", i)
		req := helios.NewMockRequest()
		req.SetContextData(auth.UserContextKey, testCase.user)
		req.URLParam["eventID"] = testCase.eventID
		req.URLParam["questionID"] = testCase.questionID

		QuestionDetailView(&req)

		assert.Equal(t, testCase.expectedStatusCode, req.StatusCode, "Unexpected status code")
		if testCase.expectedErrorCode != "" {
			var err map[string]interface{}
			json.Unmarshal(req.JSONResponse, &err)
			assert.Equal(t, testCase.expectedErrorCode, err["code"], "Different error code")
		}
	}
}

func TestSubmissionCreateView(t *testing.T) {
	beforeTest(true)
	var submissionCountBefore int
	var submissionCount int
	helios.DB.Model(&Submission{}).Count(&submissionCountBefore)

	type submissionCreateTestCase struct {
		user                    interface{}
		eventID                 string
		questionID              string
		requestData             string
		expectedStatusCode      int
		expectedErrorCode       string
		expectedSubmissionCount int
	}
	testCases := []submissionCreateTestCase{
		submissionCreateTestCase{
			user:                    user1,
			eventID:                 strconv.Itoa(int(event1.ID)),
			questionID:              strconv.Itoa(int(questionSimple.ID)),
			requestData:             `{"answer":"answer1"}`,
			expectedStatusCode:      http.StatusCreated,
			expectedSubmissionCount: submissionCountBefore + 1,
		},
		submissionCreateTestCase{
			user:                    user1,
			eventID:                 "malformed",
			questionID:              "malformed",
			requestData:             `{"answer":"answer2"}`,
			expectedStatusCode:      http.StatusNotFound,
			expectedErrorCode:       errEventNotFound.Code,
			expectedSubmissionCount: submissionCountBefore + 1,
		},
		submissionCreateTestCase{
			user:                    user1,
			eventID:                 "345678",
			questionID:              "4567987",
			requestData:             `{"answer":"answer3"}`,
			expectedStatusCode:      http.StatusNotFound,
			expectedErrorCode:       errEventNotFound.Code,
			expectedSubmissionCount: submissionCountBefore + 1,
		},
		submissionCreateTestCase{
			user:                    user1,
			eventID:                 strconv.Itoa(int(event1.ID)),
			questionID:              "malformed",
			requestData:             `{"answer":"answer4"}`,
			expectedStatusCode:      http.StatusNotFound,
			expectedErrorCode:       errQuestionNotFound.Code,
			expectedSubmissionCount: submissionCountBefore + 1,
		},
		submissionCreateTestCase{
			user:                    user1,
			eventID:                 strconv.Itoa(int(event1.ID)),
			questionID:              "876789",
			requestData:             `{"answer":"answer5"}`,
			expectedStatusCode:      http.StatusNotFound,
			expectedErrorCode:       errQuestionNotFound.Code,
			expectedSubmissionCount: submissionCountBefore + 1,
		},
		submissionCreateTestCase{
			user:                    user1,
			eventID:                 strconv.Itoa(int(event1.ID)),
			questionID:              strconv.Itoa(int(questionSimple.ID)),
			requestData:             `malformed`,
			expectedStatusCode:      http.StatusBadRequest,
			expectedSubmissionCount: submissionCountBefore + 1,
		},
		submissionCreateTestCase{
			user:                    "bad_user",
			eventID:                 strconv.Itoa(int(event1.ID)),
			questionID:              strconv.Itoa(int(questionSimple.ID)),
			requestData:             `{"answer":"answer7"}`,
			expectedStatusCode:      http.StatusInternalServerError,
			expectedSubmissionCount: submissionCountBefore + 1,
		},
	}

	for i, testCase := range testCases {
		t.Logf("Test SubmissionCreate testcase: %d", i)
		req := helios.NewMockRequest()
		req.SetContextData(auth.UserContextKey, testCase.user)
		req.URLParam["eventID"] = testCase.eventID
		req.URLParam["questionID"] = testCase.questionID
		req.RequestData = testCase.requestData

		SubmissionCreateView(&req)

		helios.DB.Model(&Submission{}).Count(&submissionCount)
		assert.Equal(t, testCase.expectedStatusCode, req.StatusCode, "Unexpected status code")
		assert.Equal(t, testCase.expectedSubmissionCount, submissionCount, "Unexpected submission count")
		if testCase.expectedErrorCode != "" {
			var err map[string]interface{}
			json.Unmarshal(req.JSONResponse, &err)
			assert.Equal(t, testCase.expectedErrorCode, err["code"], "Unexpected error code")
		}
	}
}
