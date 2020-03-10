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
}

func TestQuestionDetailView(t *testing.T) {
	beforeTest(true)

	req1 := helios.NewMockRequest()
	req1.SetContextData(auth.UserContextKey, user1)
	req1.URLParam["eventID"] = strconv.Itoa(int(event1.ID))
	req1.URLParam["questionID"] = strconv.Itoa(int(questionSimple.ID))

	QuestionDetailView(&req1)

	assert.Equal(t, http.StatusOK, req1.StatusCode, "Unexpected status code")

	req2 := helios.NewMockRequest()
	req2.SetContextData(auth.UserContextKey, user1)
	req2.URLParam["eventID"] = strconv.Itoa(int(event1.ID))
	req2.URLParam["questionID"] = "malformed"

	QuestionDetailView(&req2)

	var err2 helios.APIError
	json.Unmarshal(req2.JSONResponse, &err2)
	assert.Equal(t, http.StatusNotFound, req2.StatusCode, "Unexpected status code")
	assert.Equal(t, errQuestionNotFound.Code, err2.Code, "Different error code")

	req3 := helios.NewMockRequest()
	req3.SetContextData(auth.UserContextKey, user1)
	req3.URLParam["eventID"] = strconv.Itoa(int(event1.ID))
	req3.URLParam["questionID"] = "879654"

	QuestionDetailView(&req3)

	var err3 helios.APIError
	json.Unmarshal(req2.JSONResponse, &err3)
	assert.Equal(t, http.StatusNotFound, req3.StatusCode, "Unexpected status code")
	assert.Equal(t, errQuestionNotFound.Code, err3.Code, "Different error code")

	req4 := helios.NewMockRequest()
	req4.SetContextData(auth.UserContextKey, user1)
	req4.URLParam["eventID"] = "4567890"
	req4.URLParam["questionID"] = "56796"

	QuestionDetailView(&req4)

	var err4 helios.APIError
	json.Unmarshal(req2.JSONResponse, &err4)
	assert.Equal(t, http.StatusNotFound, req4.StatusCode, "Unexpected status code")
	assert.Equal(t, errQuestionNotFound.Code, err4.Code, "Different error code")

	req5 := helios.NewMockRequest()
	req5.SetContextData(auth.UserContextKey, user1)
	req5.URLParam["eventID"] = "malformed"
	req5.URLParam["questionID"] = "malformed"

	QuestionDetailView(&req5)

	var err5 helios.APIError
	json.Unmarshal(req2.JSONResponse, &err5)
	assert.Equal(t, http.StatusNotFound, req5.StatusCode, "Unexpected status code")
	assert.Equal(t, errQuestionNotFound.Code, err5.Code, "Different error code")
}

func TestSubmissionCreateView(t *testing.T) {
	beforeTest(true)

	var submissionCountBefore, submissionCountAfter int
	helios.DB.Model(&Submission{}).Count(&submissionCountBefore)

	req1 := helios.NewMockRequest()
	req1.SetContextData(auth.UserContextKey, user1)
	req1.URLParam["eventID"] = strconv.Itoa(int(event1.ID))
	req1.URLParam["questionID"] = strconv.Itoa(int(questionSimple.ID))
	req1.RequestData = SubmitSubmissionRequest{Answer: "answer1"}

	SubmissionCreateView(&req1)

	helios.DB.Model(&Submission{}).Count(&submissionCountAfter)

	assert.Equal(t, http.StatusCreated, req1.StatusCode, "Unexpected status code")
	assert.Equal(t, submissionCountBefore+1, submissionCountAfter, "Submission should be made")
}

func TestSubmissionCreateViewMalformedEventID(t *testing.T) {
	beforeTest(true)

	var submissionCountBefore, submissionCountAfter int
	helios.DB.Model(&Submission{}).Count(&submissionCountBefore)

	req := helios.NewMockRequest()
	req.SetContextData(auth.UserContextKey, user1)
	req.URLParam["eventID"] = "abc"
	req.URLParam["questionID"] = "abc"
	req.RequestData = SubmitSubmissionRequest{Answer: "answer1"}

	SubmissionCreateView(&req)

	helios.DB.Model(&Submission{}).Count(&submissionCountAfter)

	assert.Equal(t, http.StatusNotFound, req.StatusCode, "Unexpected status code")
	assert.Equal(t, submissionCountBefore, submissionCountAfter, "Submission should not be made")
}

func TestSubmissionCreateViewUnknownEventID(t *testing.T) {
	beforeTest(true)

	var submissionCountBefore, submissionCountAfter int
	helios.DB.Model(&Submission{}).Count(&submissionCountBefore)

	req := helios.NewMockRequest()
	req.SetContextData(auth.UserContextKey, user1)
	req.URLParam["eventID"] = "97685746"
	req.URLParam["questionID"] = "7867"
	req.RequestData = SubmitSubmissionRequest{Answer: "answer1"}

	SubmissionCreateView(&req)

	helios.DB.Model(&Submission{}).Count(&submissionCountAfter)

	assert.Equal(t, http.StatusNotFound, req.StatusCode, "Unexpected status code")
	assert.Equal(t, submissionCountBefore, submissionCountAfter, "Submission should not be made")
}

func TestSubmissionCreateViewMalformedQuestionID(t *testing.T) {
	beforeTest(true)

	var submissionCountBefore, submissionCountAfter int
	helios.DB.Model(&Submission{}).Count(&submissionCountBefore)

	req := helios.NewMockRequest()
	req.SetContextData(auth.UserContextKey, user1)
	req.URLParam["eventID"] = strconv.Itoa(int(event1.ID))
	req.URLParam["questionID"] = "abc"
	req.RequestData = SubmitSubmissionRequest{Answer: "answer1"}

	SubmissionCreateView(&req)

	helios.DB.Model(&Submission{}).Count(&submissionCountAfter)

	assert.Equal(t, http.StatusNotFound, req.StatusCode, "Unexpected status code")
	assert.Equal(t, submissionCountBefore, submissionCountAfter, "Submission should not be made")
}

func TestSubmissionCreateViewUnknownQuestionID(t *testing.T) {
	beforeTest(true)

	var submissionCountBefore, submissionCountAfter int
	helios.DB.Model(&Submission{}).Count(&submissionCountBefore)

	req := helios.NewMockRequest()
	req.SetContextData(auth.UserContextKey, user1)
	req.URLParam["eventID"] = strconv.Itoa(int(event1.ID))
	req.URLParam["questionID"] = "976857463"
	req.RequestData = SubmitSubmissionRequest{Answer: "answer1"}

	SubmissionCreateView(&req)

	helios.DB.Model(&Submission{}).Count(&submissionCountAfter)

	assert.Equal(t, http.StatusNotFound, req.StatusCode, "Unexpected status code")
	assert.Equal(t, submissionCountBefore, submissionCountAfter, "Submission should not be made")
}

func TestSubmissionCreateViewBadRequest(t *testing.T) {
	beforeTest(true)

	var submissionCountBefore, submissionCountAfter int
	helios.DB.Model(&Submission{}).Count(&submissionCountBefore)

	req := helios.NewMockRequest()
	req.SetContextData(auth.UserContextKey, user1)
	req.URLParam["eventID"] = strconv.Itoa(int(event1.ID))
	req.URLParam["questionID"] = "976857463"
	req.RequestData = nil

	SubmissionCreateView(&req)

	helios.DB.Model(&Submission{}).Count(&submissionCountAfter)

	assert.Equal(t, http.StatusUnsupportedMediaType, req.StatusCode, "Unexpected status code")
	assert.Equal(t, submissionCountBefore, submissionCountAfter, "Submission should not be made")
}
