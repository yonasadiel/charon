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
		requestData:        `{"title":"Math Final Exam","slug":"math-final-exam","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"2020-08-12T04:30:10Z"}`,
		expectedStatusCode: http.StatusCreated,
		expectedEventCount: eventCountBefore + 1,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		requestData:        `{"title":"Math Final Exam","slug":"math-final-exam","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"INVALID_END_TIME"}`,
		expectedStatusCode: http.StatusBadRequest,
		expectedErrorCode:  "form_error",
		expectedEventCount: eventCountBefore + 1,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		requestData:        `{"title":"Math Final Exam","slug":"math-final-exam","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"2020-08-12T02:30:09Z"}`,
		expectedStatusCode: http.StatusBadRequest,
		expectedErrorCode:  "form_error",
		expectedEventCount: eventCountBefore + 1,
	}, {
		user:               "bad_user",
		requestData:        `{"title":"Math Final Exam","slug":"math-final-exam","startsAt":"2020-08-12T09:30:10+07:00","endsAt":"2020-08-12T04:30:10Z"}`,
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
		assert.Equal(t, testCase.expectedStatusCode, req.StatusCode)
		assert.Equal(t, testCase.expectedEventCount, eventCount)
		if testCase.expectedErrorCode != "" {
			var err map[string]interface{}
			json.Unmarshal(req.JSONResponse, &err)
			assert.Equal(t, testCase.expectedErrorCode, err["code"])
		}
	}
}

func TestParticipationListView(t *testing.T) {
	helios.App.BeforeTest()

	var event1 Event = EventFactorySaved(Event{})
	ParticipationFactorySaved(Participation{Event: &event1})
	ParticipationFactorySaved(Participation{Event: &event1})
	type questionListTestCase struct {
		user               interface{}
		eventSlug          string
		expectedStatusCode int
		expectedErrorCode  string
	}
	testCases := []questionListTestCase{{
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:          event1.Slug,
		expectedStatusCode: http.StatusOK,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:          "random",
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errEventNotFound.Code,
	}, {
		user:               "bad_user",
		eventSlug:          event1.Slug,
		expectedStatusCode: http.StatusInternalServerError,
		expectedErrorCode:  helios.ErrInternalServerError.Code,
	}}
	for i, testCase := range testCases {
		t.Logf("Test ParticipationListView testcase: %d", i)
		req := helios.NewMockRequest()
		req.SetContextData(auth.UserContextKey, testCase.user)
		req.URLParam["eventSlug"] = testCase.eventSlug

		ParticipationListView(&req)

		assert.Equal(t, testCase.expectedStatusCode, req.StatusCode)
		if testCase.expectedErrorCode != "" {
			var err map[string]interface{}
			json.Unmarshal(req.JSONResponse, &err)
			assert.Equal(t, testCase.expectedErrorCode, err["code"])
		}
	}
}

func TestParticipationCreateView(t *testing.T) {
	helios.App.BeforeTest()

	var participationCountBefore int
	auth.UserFactorySaved(auth.User{Username: "participant", Role: auth.UserRoleParticipant})
	var userLocal auth.User = auth.UserFactorySaved(auth.User{Username: "local", Role: auth.UserRoleLocal})
	var event1 Event = EventFactorySaved(Event{})
	var venue1 Venue = VenueFactorySaved(Venue{})
	var venue2 Venue = VenueFactorySaved(Venue{})
	ParticipationFactorySaved(Participation{EventID: event1.ID, Event: &event1, UserID: userLocal.ID, User: &userLocal})
	helios.DB.Model(Participation{}).Count(&participationCountBefore)

	type participationCreateTestCase struct {
		user                       interface{}
		eventSlug                  string
		requestData                string
		expectedParticipationCount int
		expectedStatusCode         int
		expectedErrorCode          string
	}

	testCases := []participationCreateTestCase{{
		user:                       userLocal,
		eventSlug:                  event1.Slug,
		requestData:                fmt.Sprintf(`{"id":1,"userUsername":"participant","venueId":%d}`, venue1.ID),
		expectedParticipationCount: participationCountBefore + 1,
		expectedStatusCode:         http.StatusOK,
	}, {
		user:                       auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:                  "abcdef",
		requestData:                fmt.Sprintf(`{"id":1,"userUsername":"participant","venueId":%d}`, venue1.ID),
		expectedParticipationCount: participationCountBefore + 1,
		expectedStatusCode:         http.StatusNotFound,
		expectedErrorCode:          errEventNotFound.Code,
	}, {
		user:                       userLocal,
		eventSlug:                  event1.Slug,
		requestData:                fmt.Sprintf(`{"id":1,"userUsername":"local","venueId":%d}`, venue1.ID),
		expectedParticipationCount: participationCountBefore + 1,
		expectedStatusCode:         http.StatusForbidden,
		expectedErrorCode:          errParticipationChangeNotAuthorized.Code,
	}, {
		user:                       auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:                  event1.Slug,
		requestData:                `bad_format`,
		expectedParticipationCount: participationCountBefore + 1,
		expectedStatusCode:         http.StatusBadRequest,
	}, {
		user:                       auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:                  event1.Slug,
		requestData:                `{}`,
		expectedParticipationCount: participationCountBefore + 1,
		expectedStatusCode:         http.StatusBadRequest,
		expectedErrorCode:          "form_error",
	}, {
		user:                       auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:                  event1.Slug,
		requestData:                fmt.Sprintf(`{"id":1,"userUsername":"participant","venueId":%d}`, venue2.ID),
		expectedParticipationCount: participationCountBefore + 1,
		expectedStatusCode:         http.StatusOK,
	}, {
		user:                       "bad_user",
		eventSlug:                  event1.Slug,
		requestData:                fmt.Sprintf(`{"id":1,"userUsername":"participant","venueId":%d}`, venue1.ID),
		expectedParticipationCount: participationCountBefore + 1,
		expectedStatusCode:         http.StatusInternalServerError,
		expectedErrorCode:          helios.ErrInternalServerError.Code,
	}}

	for i, testCase := range testCases {
		t.Logf("Test ParticipationCreate testcase: %d", i)
		var participationCount int
		var req helios.MockRequest
		req = helios.NewMockRequest()
		req.SetContextData(auth.UserContextKey, testCase.user)
		req.URLParam["eventSlug"] = testCase.eventSlug
		req.RequestData = testCase.requestData

		ParticipationCreateView(&req)

		helios.DB.Model(Participation{}).Count(&participationCount)
		assert.Equal(t, testCase.expectedParticipationCount, participationCount)
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

func TestParticipationDeleteView(t *testing.T) {
	helios.App.BeforeTest()

	var userParticipant auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant})
	var userLocal auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal})
	var event1 Event = EventFactorySaved(Event{})
	var participation1 Participation = ParticipationFactorySaved(Participation{User: &userParticipant, Event: &event1})
	var participation2 Participation = ParticipationFactorySaved(Participation{User: &userLocal, Event: &event1})

	type questionDeleteTestCase struct {
		user               interface{}
		eventSlug          string
		participationID    string
		expectedStatusCode int
		expectedErrorCode  string
	}
	testCases := []questionDeleteTestCase{{
		user:               "bad_user",
		eventSlug:          event1.Slug,
		participationID:    strconv.Itoa(int(participation1.ID)),
		expectedStatusCode: http.StatusInternalServerError,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventSlug:          "slug",
		participationID:    "malformed",
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errParticipationNotFound.Code,
	}, {
		user:               userLocal,
		eventSlug:          event1.Slug,
		participationID:    strconv.Itoa(int(participation2.ID)),
		expectedStatusCode: errParticipationChangeNotAuthorized.StatusCode,
		expectedErrorCode:  errParticipationChangeNotAuthorized.Code,
	}, {
		user:               userLocal,
		eventSlug:          event1.Slug,
		participationID:    strconv.Itoa(int(participation1.ID)),
		expectedStatusCode: http.StatusOK,
	}}

	for i, testCase := range testCases {
		t.Logf("Test ParticipationDelete testcase: %d", i)
		req := helios.NewMockRequest()
		req.SetContextData(auth.UserContextKey, testCase.user)
		req.URLParam["eventSlug"] = testCase.eventSlug
		req.URLParam["participationID"] = testCase.participationID

		ParticipationDeleteView(&req)

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

func TestQuestionListView(t *testing.T) {
	helios.App.BeforeTest()

	var event1 Event = EventFactorySaved(Event{})
	QuestionFactorySaved(Question{Event: &event1})
	QuestionFactorySaved(Question{Event: &event1})
	type questionListTestCase struct {
		user               interface{}
		eventSlug          string
		expectedStatusCode int
		expectedErrorCode  string
	}
	testCases := []questionListTestCase{{
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:          event1.Slug,
		expectedStatusCode: http.StatusOK,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:          "random",
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errEventNotFound.Code,
	}, {
		user:               "bad_user",
		eventSlug:          event1.Slug,
		expectedStatusCode: http.StatusInternalServerError,
		expectedErrorCode:  helios.ErrInternalServerError.Code,
	}}
	for i, testCase := range testCases {
		t.Logf("Test QuestionListView testcase: %d", i)
		req := helios.NewMockRequest()
		req.SetContextData(auth.UserContextKey, testCase.user)
		req.URLParam["eventSlug"] = testCase.eventSlug

		QuestionListView(&req)

		assert.Equal(t, testCase.expectedStatusCode, req.StatusCode)
		if testCase.expectedErrorCode != "" {
			var err map[string]interface{}
			json.Unmarshal(req.JSONResponse, &err)
			assert.Equal(t, testCase.expectedErrorCode, err["code"])
		}
	}
}

func TestQuestionCreateView(t *testing.T) {
	helios.App.BeforeTest()

	var event1 Event = EventFactorySaved(Event{})

	type questionCreateTestCase struct {
		user                  interface{}
		eventSlug             string
		requestData           string
		expectedQuestionCount int
		expectedStatusCode    int
		expectedErrorCode     string
	}

	testCases := []questionCreateTestCase{{
		user:                  auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:             event1.Slug,
		requestData:           `{"id":1,"content":"content1","choices":[],"answer":"abc","eventId":2}`,
		expectedQuestionCount: 1,
		expectedStatusCode:    http.StatusCreated,
	}, {
		user:                  auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:             "abcdef",
		requestData:           `{"id":1,"content":"content2","choices":[],"answer":"abc","eventId":2}`,
		expectedQuestionCount: 1,
		expectedStatusCode:    http.StatusNotFound,
		expectedErrorCode:     errEventNotFound.Code,
	}, {
		user:                  auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant}),
		eventSlug:             event1.Slug,
		requestData:           `{"id":1,"content":"content3","choices":[],"answer":"abc","eventId":2}`,
		expectedQuestionCount: 1,
		expectedStatusCode:    http.StatusForbidden,
		expectedErrorCode:     errQuestionChangeNotAuthorized.Code,
	}, {
		user:                  auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal}),
		eventSlug:             event1.Slug,
		requestData:           `{"id":1,"content":"content3","choices":[],"answer":"abc","eventId":2}`,
		expectedQuestionCount: 1,
		expectedStatusCode:    http.StatusForbidden,
		expectedErrorCode:     errQuestionChangeNotAuthorized.Code,
	}, {
		user:                  auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:             event1.Slug,
		requestData:           `bad_format`,
		expectedQuestionCount: 1,
		expectedStatusCode:    http.StatusBadRequest,
	}, {
		user:                  auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:             event1.Slug,
		requestData:           `{"choices":["content5_1","content5_2"],"answer":"abc","eventId":2}`,
		expectedQuestionCount: 1,
		expectedStatusCode:    http.StatusBadRequest,
		expectedErrorCode:     "form_error",
	}, {
		user:                  auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:             event1.Slug,
		requestData:           `{"content":"content5","choices":["content5_1","content5_2"],"answer":"abc","eventId":2}`,
		expectedQuestionCount: 2,
		expectedStatusCode:    http.StatusCreated,
	}, {
		user:                  "bad_user",
		eventSlug:             event1.Slug,
		requestData:           `{"content":"content6","choices":[],"answer":"abc","eventId":2}`,
		expectedQuestionCount: 2,
		expectedStatusCode:    http.StatusInternalServerError,
		expectedErrorCode:     helios.ErrInternalServerError.Code,
	}}

	for i, testCase := range testCases {
		t.Logf("Test QuestionCreate testcase: %d", i)
		var questionCount int
		var req helios.MockRequest
		req = helios.NewMockRequest()
		req.SetContextData(auth.UserContextKey, testCase.user)
		req.URLParam["eventSlug"] = testCase.eventSlug
		req.RequestData = testCase.requestData

		QuestionCreateView(&req)

		helios.DB.Model(Question{}).Count(&questionCount)
		assert.Equal(t, testCase.expectedQuestionCount, questionCount)
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
		eventSlug          string
		questionID         string
		expectedStatusCode int
		expectedErrorCode  string
	}
	testCases := []questionDetailTestCase{{
		user:               userParticipant,
		eventSlug:          event1.Slug,
		questionID:         strconv.Itoa(int(question1.ID)),
		expectedStatusCode: http.StatusOK,
	}, {
		user:               userParticipant,
		eventSlug:          event1.Slug,
		questionID:         strconv.Itoa(int(question2.ID)),
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errQuestionNotFound.Code,
	}, {
		user:               userParticipant,
		eventSlug:          event1.Slug,
		questionID:         "bad_question_id",
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errQuestionNotFound.Code,
	}, {
		user:               userParticipant,
		eventSlug:          event1.Slug,
		questionID:         "879654",
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errQuestionNotFound.Code,
	}, {
		user:               userParticipant,
		eventSlug:          "random",
		questionID:         "malformed",
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errQuestionNotFound.Code,
	}, {
		user:               "bad_user",
		eventSlug:          event1.Slug,
		questionID:         strconv.Itoa(int(question1.ID)),
		expectedStatusCode: http.StatusInternalServerError,
	}}

	for i, testCase := range testCases {
		t.Logf("Test QuestionDetail testcase: %d", i)
		var req helios.MockRequest
		req = helios.NewMockRequest()
		req.SetContextData(auth.UserContextKey, testCase.user)
		req.URLParam["eventSlug"] = testCase.eventSlug
		req.URLParam["questionID"] = testCase.questionID

		QuestionDetailView(&req)

		assert.Equal(t, testCase.expectedStatusCode, req.StatusCode)
		if testCase.expectedErrorCode != "" {
			var err map[string]interface{}
			json.Unmarshal(req.JSONResponse, &err)
			assert.Equal(t, testCase.expectedErrorCode, err["code"])
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
		eventSlug          string
		questionID         string
		expectedStatusCode int
		expectedErrorCode  string
	}
	testCases := []questionDeleteTestCase{{
		user:               userParticipant,
		eventSlug:          event1.Slug,
		questionID:         strconv.Itoa(int(question1.ID)),
		expectedStatusCode: http.StatusForbidden,
		expectedErrorCode:  errQuestionChangeNotAuthorized.Code,
	}, {
		user:               userLocal,
		eventSlug:          event1.Slug,
		questionID:         strconv.Itoa(int(question1.ID)),
		expectedStatusCode: http.StatusForbidden,
		expectedErrorCode:  errQuestionChangeNotAuthorized.Code,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleOrganizer}),
		eventSlug:          event1.Slug,
		questionID:         "bad_question_id",
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errQuestionNotFound.Code,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventSlug:          event1.Slug,
		questionID:         "879654",
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errQuestionNotFound.Code,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventSlug:          "4567890",
		questionID:         "malformed",
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errQuestionNotFound.Code,
	}, {
		user:               "bad_user",
		eventSlug:          event1.Slug,
		questionID:         strconv.Itoa(int(question1.ID)),
		expectedStatusCode: http.StatusInternalServerError,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventSlug:          event1.Slug,
		questionID:         strconv.Itoa(int(question1.ID)),
		expectedStatusCode: http.StatusOK,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventSlug:          event1.Slug,
		questionID:         strconv.Itoa(int(question1.ID)),
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errQuestionNotFound.Code,
	}}

	for i, testCase := range testCases {
		t.Logf("Test QuestionDelete testcase: %d", i)
		req := helios.NewMockRequest()
		req.SetContextData(auth.UserContextKey, testCase.user)
		req.URLParam["eventSlug"] = testCase.eventSlug
		req.URLParam["questionID"] = testCase.questionID

		QuestionDeleteView(&req)

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

func TestSubmissionCreateView(t *testing.T) {
	helios.App.BeforeTest()

	var userQuestion UserQuestion = UserQuestionFactorySaved(UserQuestion{})
	var userParticipant auth.User = *userQuestion.Participation.User
	var event1 Event = *userQuestion.Question.Event
	var question1 Question = *userQuestion.Question
	type submissionCreateTestCase struct {
		user               interface{}
		eventSlug          string
		questionID         string
		requestData        string
		expectedStatusCode int
		expectedErrorCode  string
	}
	testCases := []submissionCreateTestCase{{
		user:               userParticipant,
		eventSlug:          event1.Slug,
		questionID:         strconv.Itoa(int(question1.ID)),
		requestData:        fmt.Sprintf(`{"answer":"%s"}`, question1.Choices[0].Text),
		expectedStatusCode: http.StatusCreated,
	}, {
		user:               userParticipant,
		eventSlug:          "random",
		questionID:         "malformed",
		requestData:        `{"answer":"answer2"}`,
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errQuestionNotFound.Code,
	}, {
		user:               userParticipant,
		eventSlug:          event1.Slug,
		questionID:         "malformed",
		requestData:        `{"answer":"answer4"}`,
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errQuestionNotFound.Code,
	}, {
		user:               userParticipant,
		eventSlug:          event1.Slug,
		questionID:         "876789",
		requestData:        `{"answer":"answer5"}`,
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errQuestionNotFound.Code,
	}, {
		user:               userParticipant,
		eventSlug:          event1.Slug,
		questionID:         strconv.Itoa(int(question1.ID)),
		requestData:        `malformed`,
		expectedStatusCode: http.StatusBadRequest,
	}, {
		user:               "bad_user",
		eventSlug:          event1.Slug,
		questionID:         strconv.Itoa(int(question1.ID)),
		requestData:        `{"answer":"answer7"}`,
		expectedStatusCode: http.StatusInternalServerError,
	}}

	for i, testCase := range testCases {
		t.Logf("Test SubmissionCreate testcase: %d", i)
		var req helios.MockRequest = helios.NewMockRequest()
		req.SetContextData(auth.UserContextKey, testCase.user)
		req.URLParam["eventSlug"] = testCase.eventSlug
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

func TestSynchronizationDataView(t *testing.T) {
	helios.App.BeforeTest()

	var userLocal auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleLocal})
	var userParticipant auth.User = auth.UserFactorySaved(auth.User{Role: auth.UserRoleParticipant})
	var event1 Event = EventFactorySaved(Event{})
	var event2 Event = EventFactorySaved(Event{})
	ParticipationFactorySaved(Participation{User: &userLocal, Event: &event1})
	ParticipationFactorySaved(Participation{User: &userParticipant, Event: &event1})
	QuestionFactorySaved(Question{Event: &event1})
	type synchronizationDataViewTestCase struct {
		user               interface{}
		eventSlug          string
		expectedStatusCode int
		expectedErrorCode  string
	}
	testCases := []synchronizationDataViewTestCase{{
		user:               userLocal,
		eventSlug:          event1.Slug,
		expectedStatusCode: http.StatusOK,
	}, {
		user:               auth.UserFactorySaved(auth.User{Role: auth.UserRoleAdmin}),
		eventSlug:          event1.Slug,
		expectedStatusCode: http.StatusForbidden,
		expectedErrorCode:  errSynchronizationNotAuthorized.Code,
	}, {
		user:               userLocal,
		eventSlug:          event2.Slug,
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errEventNotFound.Code,
	}, {
		user:               userLocal,
		eventSlug:          "random",
		expectedStatusCode: http.StatusNotFound,
		expectedErrorCode:  errEventNotFound.Code,
	}, {
		user:               "bad_user",
		eventSlug:          event1.Slug,
		expectedStatusCode: http.StatusInternalServerError,
	}}

	for i, testCase := range testCases {
		t.Logf("Test SynchronizationDataView testcase: %d", i)
		var req helios.MockRequest
		req = helios.NewMockRequest()
		req.SetContextData(auth.UserContextKey, testCase.user)
		req.URLParam["eventSlug"] = testCase.eventSlug

		SynchronizationDataView(&req)

		assert.Equal(t, testCase.expectedStatusCode, req.StatusCode)
		if testCase.expectedErrorCode != "" {
			var err map[string]interface{}
			json.Unmarshal(req.JSONResponse, &err)
			assert.Equal(t, testCase.expectedErrorCode, err["code"])
		}
	}
}
