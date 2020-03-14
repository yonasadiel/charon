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
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.GetStatusCode())
		return
	}

	var events []Event = GetAllEventOfUser(user)
	serializedEvents := make([]EventData, 0)
	for _, event := range events {
		serializedEvents = append(serializedEvents, SerializeEvent(event))
	}
	req.SendJSON(serializedEvents, http.StatusOK)
}

// EventCreateView creates the event
func EventCreateView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.GetStatusCode())
		return
	}

	var eventData EventData
	var event Event
	var err helios.Error
	err = req.DeserializeRequestData(&eventData)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}
	err = DeserializeEvent(eventData, &event)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}
	UpsertEvent(user, &event)
	req.SendJSON(SerializeEvent(event), http.StatusCreated)
}

// QuestionListView send list of questions
func QuestionListView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.GetStatusCode())
		return
	}

	eventID, errParseQuestionID := req.GetURLParamUint("eventID")
	if errParseQuestionID != nil {
		req.SendJSON(errQuestionNotFound.GetMessage(), errQuestionNotFound.GetStatusCode())
		return
	}

	var questions []Question
	var err helios.Error
	questions, err = GetAllQuestionOfEventAndUser(user, eventID)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}

	serializedQuestions := make([]QuestionData, 0)
	for _, question := range questions {
		serializedQuestions = append(serializedQuestions, SerializeQuestion(question))
	}
	req.SendJSON(serializedQuestions, http.StatusOK)
}

// QuestionDetailView send list of questions
func QuestionDetailView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.GetStatusCode())
		return
	}

	eventID, errParseEventID := req.GetURLParamUint("eventID")
	if errParseEventID != nil {
		req.SendJSON(errEventNotFound.GetMessage(), errEventNotFound.GetStatusCode())
		return
	}

	questionID, errParseQuestionID := req.GetURLParamUint("questionID")
	if errParseQuestionID != nil {
		req.SendJSON(errQuestionNotFound.GetMessage(), errQuestionNotFound.GetStatusCode())
		return
	}

	var question *Question
	var err helios.Error
	question, err = GetQuestionOfUser(user, eventID, questionID)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}
	var serializedQuestion QuestionData = SerializeQuestion(*question)
	req.SendJSON(serializedQuestion, http.StatusOK)
}

// SubmissionCreateView create a submission of a question
func SubmissionCreateView(req helios.Request) {
	user, ok := req.GetContextData(auth.UserContextKey).(auth.User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.GetStatusCode())
		return
	}

	eventID, errParseEventID := req.GetURLParamUint("eventID")
	if errParseEventID != nil {
		req.SendJSON(errEventNotFound.GetMessage(), errEventNotFound.GetStatusCode())
		return
	}

	questionID, errParseQuestionID := req.GetURLParamUint("questionID")
	if errParseQuestionID != nil {
		req.SendJSON(errQuestionNotFound.GetMessage(), errQuestionNotFound.GetStatusCode())
		return
	}

	var submitSubmissionRequest SubmitSubmissionRequest
	var errDeserialization helios.Error = req.DeserializeRequestData(&submitSubmissionRequest)
	if errDeserialization != nil {
		req.SendJSON(errDeserialization.GetMessage(), errDeserialization.GetStatusCode())
		return
	}

	var submission *Submission
	var err helios.Error
	submission, err = SubmitSubmission(user, eventID, questionID, submitSubmissionRequest.Answer)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
	} else {
		var questionData QuestionData = SerializeQuestion(*submission.Question)
		req.SendJSON(questionData, http.StatusCreated)
	}
}
