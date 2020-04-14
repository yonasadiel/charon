package exam

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yonasadiel/charon/backend/auth"
	"github.com/yonasadiel/helios"
)

func TestVenueListView(t *testing.T) {
	helios.App.BeforeTest()

	VenueFactorySaved(Venue{})
	VenueFactorySaved(Venue{})
	type venueListTestCase struct {
		user               interface{}
		expectedStatusCode int
	}
	testCases := []venueListTestCase{{
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		expectedStatusCode: http.StatusOK,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal}),
		expectedStatusCode: http.StatusOK,
	}, {
		user:               "bad_format",
		expectedStatusCode: http.StatusInternalServerError,
	}}
	for i, testCase := range testCases {
		t.Logf("Test VenueListView testcase: %d", i)
		var req helios.MockRequest = helios.NewMockRequest()
		req.SetContextData(auth.UserContextKey, testCase.user)
		VenueListView(&req)
		assert.Equal(t, testCase.expectedStatusCode, req.StatusCode)
	}
}

func TestVenueCreateView(t *testing.T) {
	helios.App.BeforeTest()

	type venueCreateTestCase struct {
		user               interface{}
		requestData        string
		expectedStatusCode int
		expectedErrorCode  string
		expectedVenueCount int
	}
	testCases := []venueCreateTestCase{{
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		requestData:        `{"id":2,"name":"Venue 1"}`,
		expectedStatusCode: http.StatusCreated,
		expectedVenueCount: 1,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		requestData:        `{"name":""}`,
		expectedStatusCode: http.StatusBadRequest,
		expectedErrorCode:  "form_error",
		expectedVenueCount: 1,
	}, {
		user:               "bad_user",
		requestData:        `{"name":""}`,
		expectedStatusCode: http.StatusInternalServerError,
		expectedErrorCode:  helios.ErrInternalServerError.Code,
		expectedVenueCount: 1,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		requestData:        `bad_request_data`,
		expectedStatusCode: http.StatusBadRequest,
		expectedErrorCode:  helios.ErrJSONParseFailed.Code,
		expectedVenueCount: 1,
	}}
	for i, testCase := range testCases {
		t.Logf("Test VenueCreateView testcase: %d", i)
		var eventCount int
		var req helios.MockRequest
		req = helios.NewMockRequest()
		req.SetContextData(auth.UserContextKey, testCase.user)
		req.RequestData = testCase.requestData

		VenueCreateView(&req)

		helios.DB.Model(Venue{}).Count(&eventCount)
		assert.Equal(t, testCase.expectedStatusCode, req.StatusCode)
		assert.Equal(t, testCase.expectedVenueCount, eventCount)
		if testCase.expectedErrorCode != "" {
			var err map[string]interface{}
			json.Unmarshal(req.JSONResponse, &err)
			assert.Equal(t, testCase.expectedErrorCode, err["code"])
		}
	}
}

func TestVenueDeleteView(t *testing.T) {
	helios.App.BeforeTest()

	var venue1 Venue = VenueFactorySaved(Venue{})
	var venue2 Venue = VenueFactorySaved(Venue{})
	ParticipationFactorySaved(Participation{Venue: &venue2})

	type venueDeleteViewTestCase struct {
		user               interface{}
		venueID            string
		expectedStatusCode int
		expectedErrorCode  string
	}
	testCases := []venueDeleteViewTestCase{{
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant}),
		venueID:            strconv.Itoa(int(venue1.ID)),
		expectedStatusCode: http.StatusForbidden,
		expectedErrorCode:  errVenueAccessNotAuthorized.Code,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal}),
		venueID:            strconv.Itoa(int(venue1.ID)),
		expectedStatusCode: http.StatusForbidden,
		expectedErrorCode:  errVenueAccessNotAuthorized.Code,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		venueID:            "bad_venue_id",
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errVenueNotFound.Code,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		venueID:            "879654",
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errVenueNotFound.Code,
	}, {
		user:               "bad_user",
		venueID:            strconv.Itoa(int(venue1.ID)),
		expectedStatusCode: http.StatusInternalServerError,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		venueID:            strconv.Itoa(int(venue1.ID)),
		expectedStatusCode: http.StatusOK,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		venueID:            strconv.Itoa(int(venue2.ID)),
		expectedStatusCode: errVenueCantDeletedEventExists.StatusCode,
		expectedErrorCode:  errVenueCantDeletedEventExists.Code,
	}}

	for i, testCase := range testCases {
		t.Logf("Test VenueDelete testcase: %d", i)
		req := helios.NewMockRequest()
		req.SetContextData(auth.UserContextKey, testCase.user)
		req.URLParam["venueID"] = testCase.venueID

		VenueDeleteView(&req)

		assert.Equal(t, testCase.expectedStatusCode, req.StatusCode, "Unexpected status code")
		if testCase.expectedErrorCode != "" {
			var err map[string]interface{}
			var errUnmarshalling error
			errUnmarshalling = json.Unmarshal(req.JSONResponse, &err)
			assert.Nil(t, errUnmarshalling)
			assert.Equal(t, testCase.expectedErrorCode, err["code"], "Different error code")
		}
	}
}

func TestEventListView(t *testing.T) {
	helios.App.BeforeTest()

	EventFactorySaved(Event{})
	EventFactorySaved(Event{})
	type eventListTestCase struct {
		user               interface{}
		expectedStatusCode int
	}
	testCases := []eventListTestCase{{
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		expectedStatusCode: http.StatusOK,
	}, {
		user:               "bad_format",
		expectedStatusCode: http.StatusInternalServerError,
	}}
	for i, testCase := range testCases {
		t.Logf("Test EventListView testcase: %d", i)
		req := helios.NewMockRequest()
		req.SetContextData(auth.UserContextKey, testCase.user)
		EventListView(&req)
		assert.Equal(t, testCase.expectedStatusCode, req.StatusCode)
	}
}

func TestEventCreateView(t *testing.T) {
	helios.App.BeforeTest()
	var eventCountBefore int
	helios.DB.Model(Event{}).Count(&eventCountBefore)

	type eventCreateTestCase struct {
		user               interface{}
		requestData        string
		expectedStatusCode int
		expectedErrorCode  string
		expectedEventCount int
	}
	testCases := []eventCreateTestCase{{
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		requestData:        `{"title":"Math Final Exam","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"2020-08-12T04:30:10Z"}`,
		expectedStatusCode: http.StatusCreated,
		expectedEventCount: eventCountBefore + 1,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		requestData:        `{"title":"Math Final Exam","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"INVALID_END_TIME"}`,
		expectedStatusCode: http.StatusBadRequest,
		expectedErrorCode:  "form_error",
		expectedEventCount: eventCountBefore + 1,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		requestData:        `{"title":"Math Final Exam","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"2020-08-12T02:30:09Z"}`,
		expectedStatusCode: http.StatusBadRequest,
		expectedErrorCode:  "form_error",
		expectedEventCount: eventCountBefore + 1,
	}, {
		user:               "bad_user",
		requestData:        `{"title":"Math Final Exam","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"2020-08-12T04:30:10Z"}`,
		expectedStatusCode: http.StatusInternalServerError,
		expectedErrorCode:  helios.ErrInternalServerError.Code,
		expectedEventCount: eventCountBefore + 1,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		requestData:        `bad_request_data`,
		expectedStatusCode: http.StatusBadRequest,
		expectedErrorCode:  helios.ErrJSONParseFailed.Code,
		expectedEventCount: eventCountBefore + 1,
	}}
	for i, testCase := range testCases {
		t.Logf("Test EventCreate testcase: %d", i)
		var eventCount int
		var req helios.MockRequest
		req = helios.NewMockRequest()
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
	helios.App.BeforeTest()

	var event1 Event = EventFactorySaved(Event{})
	QuestionFactorySaved(Question{Event: &event1})
	QuestionFactorySaved(Question{Event: &event1})
	type questionListTestCase struct {
		user               interface{}
		eventID            string
		expectedStatusCode int
		expectedErrorCode  string
	}
	testCases := []questionListTestCase{{
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventID:            strconv.Itoa(int(event1.ID)),
		expectedStatusCode: http.StatusOK,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventID:            "malformed",
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errEventNotFound.Code,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventID:            "79697",
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errEventNotFound.Code,
	}, {
		user:               "bad_user",
		eventID:            strconv.Itoa(int(event1.ID)),
		expectedStatusCode: http.StatusInternalServerError,
		expectedErrorCode:  helios.ErrInternalServerError.Code,
	}}
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
	helios.App.BeforeTest()

	var questionCountBefore int
	var event1 Event = EventFactorySaved(Event{})
	helios.DB.Model(Question{}).Count(&questionCountBefore)

	type questionCreateTestCase struct {
		user                  interface{}
		eventID               string
		requestData           string
		expectedQuestionCount int
		expectedStatusCode    int
		expectedErrorCode     string
	}

	testCases := []questionCreateTestCase{{
		user:                  auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventID:               strconv.Itoa(int(event1.ID)),
		requestData:           `{"id":1,"content":"content1","choices":[],"answer":"abc","eventId":2}`,
		expectedQuestionCount: questionCountBefore + 1,
		expectedStatusCode:    http.StatusCreated,
	}, {
		user:                  auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventID:               "abcdef",
		requestData:           `{"id":1,"content":"content2","choices":[],"answer":"abc","eventId":2}`,
		expectedQuestionCount: questionCountBefore + 1,
		expectedStatusCode:    http.StatusNotFound,
		expectedErrorCode:     errEventNotFound.Code,
	}, {
		user:                  auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventID:               "999",
		requestData:           `{"id":1,"content":"content2","choices":[],"answer":"abc","eventId":2}`,
		expectedQuestionCount: questionCountBefore + 1,
		expectedStatusCode:    http.StatusNotFound,
		expectedErrorCode:     errEventNotFound.Code,
	}, {
		user:                  auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant}),
		eventID:               strconv.Itoa(int(event1.ID)),
		requestData:           `{"id":1,"content":"content3","choices":[],"answer":"abc","eventId":2}`,
		expectedQuestionCount: questionCountBefore + 1,
		expectedStatusCode:    http.StatusForbidden,
		expectedErrorCode:     errQuestionChangeNotAuthorized.Code,
	}, {
		user:                  auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal}),
		eventID:               strconv.Itoa(int(event1.ID)),
		requestData:           `{"id":1,"content":"content3","choices":[],"answer":"abc","eventId":2}`,
		expectedQuestionCount: questionCountBefore + 1,
		expectedStatusCode:    http.StatusForbidden,
		expectedErrorCode:     errQuestionChangeNotAuthorized.Code,
	}, {
		user:                  auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventID:               strconv.Itoa(int(event1.ID)),
		requestData:           `bad_format`,
		expectedQuestionCount: questionCountBefore + 1,
		expectedStatusCode:    http.StatusBadRequest,
	}, {
		user:                  auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventID:               strconv.Itoa(int(event1.ID)),
		requestData:           `{"content":"content5","choices":["content5_1","content5_2"],"answer":"abc","eventId":2}`,
		expectedQuestionCount: questionCountBefore + 2,
		expectedStatusCode:    http.StatusCreated,
	}, {
		user:                  "bad_user",
		eventID:               strconv.Itoa(int(event1.ID)),
		requestData:           `{"content":"content6","choices":[],"answer":"abc","eventId":2}`,
		expectedQuestionCount: questionCountBefore + 2,
		expectedStatusCode:    http.StatusInternalServerError,
		expectedErrorCode:     helios.ErrInternalServerError.Code,
	}}

	for i, testCase := range testCases {
		t.Logf("Test QuestionCreate testcase: %d", i)
		var questionCount int
		var req helios.MockRequest
		req = helios.NewMockRequest()
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
	helios.App.BeforeTest()

	var userParticipant auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant})
	var event1 Event = EventFactorySaved(Event{})
	var question1 Question = QuestionFactorySaved(Question{Event: &event1})
	var question2 Question = QuestionFactorySaved(Question{Event: &event1})
	var participation1 Participation = ParticipationFactorySaved(Participation{User: &userParticipant, Event: &event1})
	UserQuestionFactorySaved(UserQuestion{Participation: &participation1, Question: &question1})
	type questionDetailTestCase struct {
		user               interface{}
		eventID            string
		questionID         string
		expectedStatusCode int
		expectedErrorCode  string
	}
	testCases := []questionDetailTestCase{{
		user:               userParticipant,
		eventID:            strconv.Itoa(int(event1.ID)),
		questionID:         strconv.Itoa(int(question1.ID)),
		expectedStatusCode: http.StatusOK,
	}, {
		user:               userParticipant,
		eventID:            strconv.Itoa(int(event1.ID)),
		questionID:         strconv.Itoa(int(question2.ID)),
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errQuestionNotFound.Code,
	}, {
		user:               userParticipant,
		eventID:            strconv.Itoa(int(event1.ID)),
		questionID:         "bad_question_id",
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errQuestionNotFound.Code,
	}, {
		user:               userParticipant,
		eventID:            strconv.Itoa(int(event1.ID)),
		questionID:         "879654",
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errQuestionNotFound.Code,
	}, {
		user:               userParticipant,
		eventID:            "4567890",
		questionID:         "879654",
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errEventNotFound.Code,
	}, {
		user:               userParticipant,
		eventID:            "malformed",
		questionID:         "malformed",
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errEventNotFound.Code,
	}, {
		user:               "bad_user",
		eventID:            strconv.Itoa(int(event1.ID)),
		questionID:         strconv.Itoa(int(question1.ID)),
		expectedStatusCode: http.StatusInternalServerError,
	}}

	for i, testCase := range testCases {
		t.Logf("Test QuestionDetail testcase: %d", i)
		var req helios.MockRequest
		req = helios.NewMockRequest()
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

func TestQuestionDeleteView(t *testing.T) {
	helios.App.BeforeTest()

	var userParticipant auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant})
	var userLocal auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal})
	var event1 Event = EventFactorySaved(Event{})
	var question1 Question = QuestionFactorySaved(Question{Event: &event1})
	var participation1 Participation = ParticipationFactorySaved(Participation{User: &userParticipant, Event: &event1})
	ParticipationFactorySaved(Participation{User: &userLocal, Event: &event1})
	UserQuestionFactorySaved(UserQuestion{Participation: &participation1, Question: &question1})

	type questionDeleteTestCase struct {
		user               interface{}
		eventID            string
		questionID         string
		expectedStatusCode int
		expectedErrorCode  string
	}
	testCases := []questionDeleteTestCase{{
		user:               userParticipant,
		eventID:            strconv.Itoa(int(event1.ID)),
		questionID:         strconv.Itoa(int(question1.ID)),
		expectedStatusCode: http.StatusForbidden,
		expectedErrorCode:  errQuestionChangeNotAuthorized.Code,
	}, {
		user:               userLocal,
		eventID:            strconv.Itoa(int(event1.ID)),
		questionID:         strconv.Itoa(int(question1.ID)),
		expectedStatusCode: http.StatusForbidden,
		expectedErrorCode:  errQuestionChangeNotAuthorized.Code,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventID:            strconv.Itoa(int(event1.ID)),
		questionID:         "bad_question_id",
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errQuestionNotFound.Code,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventID:            strconv.Itoa(int(event1.ID)),
		questionID:         "879654",
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errQuestionNotFound.Code,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventID:            "4567890",
		questionID:         "879654",
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errEventNotFound.Code,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventID:            "malformed",
		questionID:         "malformed",
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errEventNotFound.Code,
	}, {
		user:               "bad_user",
		eventID:            strconv.Itoa(int(event1.ID)),
		questionID:         strconv.Itoa(int(question1.ID)),
		expectedStatusCode: http.StatusInternalServerError,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventID:            strconv.Itoa(int(event1.ID)),
		questionID:         strconv.Itoa(int(question1.ID)),
		expectedStatusCode: http.StatusOK,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventID:            strconv.Itoa(int(event1.ID)),
		questionID:         strconv.Itoa(int(question1.ID)),
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errQuestionNotFound.Code,
	}}

	for i, testCase := range testCases {
		t.Logf("Test QuestionDelete testcase: %d", i)
		req := helios.NewMockRequest()
		req.SetContextData(auth.UserContextKey, testCase.user)
		req.URLParam["eventID"] = testCase.eventID
		req.URLParam["questionID"] = testCase.questionID

		QuestionDeleteView(&req)

		assert.Equal(t, testCase.expectedStatusCode, req.StatusCode, "Unexpected status code")
		if testCase.expectedErrorCode != "" {
			var err map[string]interface{}
			var errUnmarshalling error
			errUnmarshalling = json.Unmarshal(req.JSONResponse, &err)
			assert.Nil(t, errUnmarshalling)
			assert.Equal(t, testCase.expectedErrorCode, err["code"], "Different error code")
		}
	}
}

func TestSubmissionCreateView(t *testing.T) {
	helios.App.BeforeTest()

	var userQuestion UserQuestion = UserQuestionFactorySaved(UserQuestion{})
	var userParticipant auth.User = *userQuestion.Participation.User
	var event1 Event = *userQuestion.Question.Event
	var question1 Question = *userQuestion.Question
	type submissionCreateTestCase struct {
		user               interface{}
		eventID            string
		questionID         string
		requestData        string
		expectedStatusCode int
		expectedErrorCode  string
	}
	testCases := []submissionCreateTestCase{{
		user:               userParticipant,
		eventID:            strconv.Itoa(int(event1.ID)),
		questionID:         strconv.Itoa(int(question1.ID)),
		requestData:        fmt.Sprintf(`{"answer":"%s"}`, question1.Choices[0].Text),
		expectedStatusCode: http.StatusCreated,
	}, {
		user:               userParticipant,
		eventID:            "malformed",
		questionID:         "malformed",
		requestData:        `{"answer":"answer2"}`,
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errEventNotFound.Code,
	}, {
		user:               userParticipant,
		eventID:            "345678",
		questionID:         "4567987",
		requestData:        `{"answer":"answer3"}`,
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errEventNotFound.Code,
	}, {
		user:               userParticipant,
		eventID:            strconv.Itoa(int(event1.ID)),
		questionID:         "malformed",
		requestData:        `{"answer":"answer4"}`,
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errQuestionNotFound.Code,
	}, {
		user:               userParticipant,
		eventID:            strconv.Itoa(int(event1.ID)),
		questionID:         "876789",
		requestData:        `{"answer":"answer5"}`,
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errQuestionNotFound.Code,
	}, {
		user:               userParticipant,
		eventID:            strconv.Itoa(int(event1.ID)),
		questionID:         strconv.Itoa(int(question1.ID)),
		requestData:        `malformed`,
		expectedStatusCode: http.StatusBadRequest,
	}, {
		user:               "bad_user",
		eventID:            strconv.Itoa(int(event1.ID)),
		questionID:         strconv.Itoa(int(question1.ID)),
		requestData:        `{"answer":"answer7"}`,
		expectedStatusCode: http.StatusInternalServerError,
	}}

	for i, testCase := range testCases {
		t.Logf("Test SubmissionCreate testcase: %d", i)
		var req helios.MockRequest = helios.NewMockRequest()
		req.SetContextData(auth.UserContextKey, testCase.user)
		req.URLParam["eventID"] = testCase.eventID
		req.URLParam["questionID"] = testCase.questionID
		req.RequestData = testCase.requestData

		SubmissionCreateView(&req)

		assert.Equal(t, testCase.expectedStatusCode, req.StatusCode)
		if testCase.expectedErrorCode != "" {
			var err map[string]interface{}
			var errUnmarshalling error
			errUnmarshalling = json.Unmarshal(req.JSONResponse, &err)
			assert.Nil(t, errUnmarshalling)
			assert.Equal(t, testCase.expectedErrorCode, err["code"])
		}
	}
}
