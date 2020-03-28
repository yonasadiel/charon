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

	user1.SetAsOrganizer()
	var eventCountBefore, eventCountAfter int
	helios.DB.Model(Event{}).Count(&eventCountBefore)

	req := helios.NewMockRequest()
	req.SetContextData(auth.UserContextKey, user1)
	req.RequestData = `{"title":"Math Final Exam","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"2020-08-12T04:30:10Z"}`
	EventCreateView(&req)
	helios.DB.Model(Event{}).Count(&eventCountAfter)

	assert.Equal(t, eventCountBefore+1, eventCountAfter, "Event should be added to database")
	assert.Equal(t, http.StatusCreated, req.StatusCode, "Unexpected status code")

	req.RequestData = `{"title":"Math Final Exam","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"INVALID_END_TIME"}`
	EventCreateView(&req)
	helios.DB.Model(Event{}).Count(&eventCountAfter)

	assert.Equal(t, eventCountBefore+1, eventCountAfter, "Event should not be added to database")
	assert.Equal(t, http.StatusBadRequest, req.StatusCode, "Unexpected status code")

	req.RequestData = `{"title":"Math Final Exam","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"2020-08-12T02:30:09Z"}`
	EventCreateView(&req)
	helios.DB.Model(Event{}).Count(&eventCountAfter)

	assert.Equal(t, eventCountBefore+1, eventCountAfter, "Event should note be added to database")
	assert.Equal(t, http.StatusBadRequest, req.StatusCode, "Unexpected status code")

	req.SetContextData(auth.UserContextKey, "bad_user")
	EventCreateView(&req)
	helios.DB.Model(Event{}).Count(&eventCountAfter)

	assert.Equal(t, eventCountBefore+1, eventCountAfter, "Event should note be added to database")
	assert.Equal(t, http.StatusInternalServerError, req.StatusCode, "Unexpected status code")
}

func TestQuestionListView(t *testing.T) {
	beforeTest(true)

	req1 := helios.NewMockRequest()
	req1.SetContextData(auth.UserContextKey, user1)
	req1.URLParam["eventID"] = strconv.Itoa(int(event1.ID))

	QuestionListView(&req1)

	req2 := helios.NewMockRequest()
	req2.SetContextData(auth.UserContextKey, user1)
	req2.URLParam["eventID"] = "abcdef"

	QuestionListView(&req2)

	assert.Equal(t, http.StatusNotFound, req2.StatusCode, "eventID is not configured correctly")

	req3 := helios.NewMockRequest()
	req3.SetContextData(auth.UserContextKey, user1)
	req3.URLParam["eventID"] = "8900"

	QuestionListView(&req3)

	assert.Equal(t, http.StatusNotFound, req3.StatusCode, "eventID is not exist on database")

	req4 := helios.NewMockRequest()
	req4.SetContextData(auth.UserContextKey, "bad_user")
	req4.URLParam["eventID"] = strconv.Itoa(int(event1.ID))

	QuestionListView(&req4)

	assert.Equal(t, http.StatusInternalServerError, req4.StatusCode, "eventID is not exist on database")
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
		},
		questionCreateTestCase{
			user:                  userOrganizer,
			eventID:               "999",
			requestData:           `{"id":1,"content":"content2","choices":[],"answer":"abc","eventId":2}`,
			expectedQuestionCount: questionCountBefore + 1,
			expectedStatusCode:    http.StatusNotFound,
		},
		questionCreateTestCase{
			user:                  userParticipant,
			eventID:               strconv.Itoa(int(event1.ID)),
			requestData:           `{"id":1,"content":"content3","choices":[],"answer":"abc","eventId":2}`,
			expectedQuestionCount: questionCountBefore + 1,
			expectedStatusCode:    http.StatusUnauthorized,
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
			user:                  "bad_user_data",
			eventID:               strconv.Itoa(int(event1.ID)),
			requestData:           `{"content":"content6","choices":[],"answer":"abc","eventId":2}`,
			expectedQuestionCount: questionCountBefore + 2,
			expectedStatusCode:    http.StatusInternalServerError,
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
			assert.Equal(t, testCase.expectedErrorCode, testCase.expectedErrorCode, "Unexpected error code")
		}
	}
}
