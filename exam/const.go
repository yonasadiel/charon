package exam

import (
	"net/http"

	"github.com/yonasadiel/helios"
)

var errEventNotFound = helios.APIError{
	StatusCode: http.StatusNotFound,
	Code:       "event_not_found",
	Message:    "No event with given ID",
}

var errEventChangeNotAuthorized = helios.APIError{
	StatusCode: http.StatusUnauthorized,
	Code:       "not_authorized_edit_event",
	Message:    "User is not authorized to make changes on event",
}

var errQuestionChangeNotAuthorized = helios.APIError{
	StatusCode: http.StatusUnauthorized,
	Code:       "not_authorized_edit_question",
	Message:    "User is not authorized to make changes on question",
}

var errQuestionNotFound = helios.APIError{
	StatusCode: http.StatusNotFound,
	Code:       "question_not_found",
	Message:    "No question with given ID",
}

var errAnswerNotValid = helios.APIError{
	StatusCode: http.StatusBadRequest,
	Code:       "invalid_answer",
	Message:    "Answer is not valid",
}
