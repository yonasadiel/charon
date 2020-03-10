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
