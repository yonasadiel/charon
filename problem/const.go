package problem

import (
	"net/http"

	"github.com/yonasadiel/helios"
)

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
