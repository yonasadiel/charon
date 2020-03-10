package exam

import (
	"net/http"

	"github.com/yonasadiel/charon/auth"
	"github.com/yonasadiel/helios"
)

// EventListView send list of questions
func EventListView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.StatusCode)
		return
	}

	var events []Event = GetAllEventOfUser(user)
	serializedEvents := make([]EventResponse, 0)
	for _, event := range events {
		serializedEvents = append(serializedEvents, SerializeEvent(event))
	}
	req.SendJSON(serializedEvents, http.StatusOK)
}

// QuestionListView send list of questions
func QuestionListView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.StatusCode)
		return
	}

	eventID, errParseQuestionID := req.GetURLParamUint("eventID")
	if errParseQuestionID != nil {
		req.SendJSON(errQuestionNotFound.GetMessage(), errQuestionNotFound.StatusCode)
		return
	}

	var questions []Question
	var err *helios.APIError
	questions, err = GetAllQuestionOfEventAndUser(user, eventID)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.StatusCode)
		return
	}

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

	eventID, errParseEventID := req.GetURLParamUint("eventID")
	if errParseEventID != nil {
		req.SendJSON(errEventNotFound.GetMessage(), errEventNotFound.StatusCode)
		return
	}

	questionID, errParseQuestionID := req.GetURLParamUint("questionID")
	if errParseQuestionID != nil {
		req.SendJSON(errQuestionNotFound.GetMessage(), errQuestionNotFound.StatusCode)
		return
	}

	var question *Question
	var err *helios.APIError
	question, err = GetQuestionOfUser(user, eventID, questionID)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.StatusCode)
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

	eventID, errParseEventID := req.GetURLParamUint("eventID")
	if errParseEventID != nil {
		req.SendJSON(errEventNotFound.GetMessage(), errEventNotFound.StatusCode)
		return
	}

	questionID, errParseQuestionID := req.GetURLParamUint("questionID")
	if errParseQuestionID != nil {
		req.SendJSON(errQuestionNotFound.GetMessage(), errQuestionNotFound.StatusCode)
		return
	}

	var submitSubmissionRequest SubmitSubmissionRequest
	var errDeserialization *helios.APIError = req.DeserializeRequestData(&submitSubmissionRequest)
	if errDeserialization != nil {
		req.SendJSON(errDeserialization.GetMessage(), errDeserialization.StatusCode)
		return
	}

	var submission *Submission
	var err *helios.APIError
	submission, err = SubmitSubmission(user, eventID, questionID, submitSubmissionRequest.Answer)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.StatusCode)
	} else {
		var questionResponse QuestionResponse = SerializeQuestion(*submission.Question)
		req.SendJSON(questionResponse, http.StatusCreated)
	}
}
