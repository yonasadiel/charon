package problem

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yonasadiel/charon/auth"
	"github.com/yonasadiel/helios"
)

func TestQuestionListView(t *testing.T) {
	beforeTest(true)

	req := helios.NewMockRequest()
	req.SetContextData(auth.UserContextKey, user1)

	QuestionListView(&req)

	assert.Equal(t, http.StatusOK, req.StatusCode, "Unexpected status code")
}

func TestQuestionDetailView(t *testing.T) {
	beforeTest(true)

	req1 := helios.NewMockRequest()
	req1.SetContextData(auth.UserContextKey, user1)
	req1.URLParam["questionId"] = strconv.Itoa(int(questionSimple.ID))

	QuestionDetailView(&req1)

	assert.Equal(t, http.StatusOK, req1.StatusCode, "Unexpected status code")

	req2 := helios.NewMockRequest()
	req2.SetContextData(auth.UserContextKey, user1)
	req2.URLParam["questionId"] = "malformed"

	QuestionDetailView(&req2)

	assert.Equal(t, http.StatusNotFound, req2.StatusCode, "Unexpected status code")

	req3 := helios.NewMockRequest()
	req3.SetContextData(auth.UserContextKey, user1)
	req3.URLParam["questionId"] = "879654"

	QuestionDetailView(&req3)

	assert.Equal(t, http.StatusNotFound, req3.StatusCode, "Unexpected status code")
}

func TestSubmissionCreateView(t *testing.T) {
	beforeTest(true)

	var submissionCountBefore, submissionCountAfter int
	helios.DB.Model(&Submission{}).Count(&submissionCountBefore)

	req := helios.NewMockRequest()
	req.SetContextData(auth.UserContextKey, user1)
	req.URLParam["questionId"] = strconv.Itoa(int(questionSimple.ID))
	req.RequestData = SubmitSubmissionRequest{Answer: "answer1"}

	SubmissionCreateView(&req)

	helios.DB.Model(&Submission{}).Count(&submissionCountAfter)

	assert.Equal(t, http.StatusCreated, req.StatusCode, "Unexpected status code")
	assert.Equal(t, submissionCountBefore+1, submissionCountAfter, "Submission should be made")
}

func TestSubmissionCreateViewMalformedQuestionID(t *testing.T) {
	beforeTest(true)

	var submissionCountBefore, submissionCountAfter int
	helios.DB.Model(&Submission{}).Count(&submissionCountBefore)

	req := helios.NewMockRequest()
	req.SetContextData(auth.UserContextKey, user1)
	req.URLParam["questionId"] = "abc"
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
	req.URLParam["questionId"] = "976857463"
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
	req.URLParam["questionId"] = "976857463"
	req.RequestData = nil

	SubmissionCreateView(&req)

	helios.DB.Model(&Submission{}).Count(&submissionCountAfter)

	assert.Equal(t, http.StatusUnsupportedMediaType, req.StatusCode, "Unexpected status code")
	assert.Equal(t, submissionCountBefore, submissionCountAfter, "Submission should not be made")
}
