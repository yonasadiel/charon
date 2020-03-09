package exam

import (
	"net/http"
	"strconv"

	"github.com/yonasadiel/charon/auth"
	"github.com/yonasadiel/helios"
)

// QuestionListView send list of questions
func QuestionListView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.StatusCode)
		return
	}

	var questions []Question = GetAllQuestionOfUser(user)
	serializedQuestions := make([]QuestionResponse, 0)
	for _, question := range questions {
		serializedQuestions = append(serializedQuestions, SerializeQuestion(question))
	}
	req.SendJSON(serializedQuestions, http.StatusOK)
}

// QuestionDetailView send list of questions
func QuestionDetailView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.StatusCode)
		return
	}

	questionIDStr := req.GetURLParam("questionId")
	questionID64, errParseQuestionID := strconv.ParseUint(questionIDStr, 10, 32)
	if errParseQuestionID != nil {
		req.SendJSON(errQuestionNotFound.GetMessage(), errQuestionNotFound.StatusCode)
		return
	}
	questionID := uint(questionID64)

	var question *Question = GetQuestionOfUser(questionID, user)
	if question == nil {
		req.SendJSON(errQuestionNotFound.GetMessage(), errQuestionNotFound.StatusCode)
		return
	}
	var serializedQuestion QuestionResponse = SerializeQuestion(*question)
	req.SendJSON(serializedQuestion, http.StatusOK)
}

// SubmissionCreateView create a submission of a question
func SubmissionCreateView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.StatusCode)
		return
	}
	questionIDStr := req.GetURLParam("questionId")
	questionID64, errParseQuestionID := strconv.ParseUint(questionIDStr, 10, 32)
	if errParseQuestionID != nil {
		req.SendJSON(errQuestionNotFound.GetMessage(), errQuestionNotFound.StatusCode)
		return
	}
	questionID := uint(questionID64)

	var submitSubmissionRequest SubmitSubmissionRequest
	var errDeserialization *helios.APIError = req.DeserializeRequestData(&submitSubmissionRequest)
	if errDeserialization != nil {
		req.SendJSON(errDeserialization.GetMessage(), errDeserialization.StatusCode)
		return
	}

	var submission *Submission
	var err *helios.APIError
	submission, err = SubmitSubmission(questionID, user, submitSubmissionRequest.Answer)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.StatusCode)
	} else {
		var questionResponse QuestionResponse = SerializeQuestion(*submission.Question)
		req.SendJSON(questionResponse, http.StatusCreated)
	}
}
