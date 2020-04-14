package exam

import (
	"net/http"

	"github.com/yonasadiel/helios"
)

var errVenueAccessNotAuthorized = helios.APIError{
	StatusCode: http.StatusForbidden,
	Code:       "venue_access_forbidden",
	Message:    "You don't have permission to access venue",
}

var errVenueNotFound = helios.APIError{
	StatusCode: http.StatusNotFound,
	Code:       "venue_not_found",
	Message:    "No venue with given ID",
}

var errVenueCantDeletedEventExists = helios.APIError{
	StatusCode: http.StatusNotFound,
	Code:       "venue_cant_deleted_event_exists",
	Message:    "The venue can't be deleted because there is event existed on the venue",
}

var errEventNotFound = helios.APIError{
	StatusCode: http.StatusNotFound,
	Code:       "event_not_found",
	Message:    "No event with given ID",
}

var errEventChangeNotAuthorized = helios.APIError{
	StatusCode: http.StatusForbidden,
	Code:       "not_authorized_edit_event",
	Message:    "User is not authorized to make changes on event",
}

var errEventIsNotYetStarted = helios.APIError{
	StatusCode: http.StatusForbidden,
	Code:       "event_is_not_yet_started",
	Message:    "The event is not yet started",
}

var errParticipationChangeNotAuthorized = helios.APIError{
	StatusCode: http.StatusForbidden,
	Code:       "not_authorized_edit_participation",
	Message:    "User is not authorized to make changes on participation",
}

var errQuestionChangeNotAuthorized = helios.APIError{
	StatusCode: http.StatusForbidden,
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

var errSubmissionNotAuthorized = helios.APIError{
	StatusCode: http.StatusForbidden,
	Code:       "cannot_submit_submission",
	Message:    "You are not allowed to submit to this question",
}
