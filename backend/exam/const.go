package exam

import (
	"net/http"

	"github.com/yonasadiel/helios"
)

var errVenueAccessNotAuthorized = helios.ErrorAPI{
	StatusCode: http.StatusForbidden,
	Code:       "venue_access_forbidden",
	Message:    "You don't have permission to access venue",
}

var errVenueNotFound = helios.ErrorAPI{
	StatusCode: http.StatusNotFound,
	Code:       "venue_not_found",
	Message:    "No venue with given ID",
}

var errVenueCantDeletedEventExists = helios.ErrorAPI{
	StatusCode: http.StatusNotFound,
	Code:       "venue_cant_deleted_event_exists",
	Message:    "The venue can't be deleted because there is event existed on the venue",
}

var errEventNotFound = helios.ErrorAPI{
	StatusCode: http.StatusNotFound,
	Code:       "event_not_found",
	Message:    "No event with given slug",
}

var errEventChangeNotAuthorized = helios.ErrorAPI{
	StatusCode: http.StatusForbidden,
	Code:       "not_authorized_edit_event",
	Message:    "User is not authorized to make changes on event",
}

var errEventIsNotYetStarted = helios.ErrorAPI{
	StatusCode: http.StatusForbidden,
	Code:       "event_is_not_yet_started",
	Message:    "The event is not yet started",
}
var errEventIsEncrypted = helios.ErrorAPI{
	StatusCode: http.StatusBadRequest,
	Code:       "event_is_encrypted",
	Message:    "Event is currently encrypted",
}

var errParticipationChangeNotAuthorized = helios.ErrorAPI{
	StatusCode: http.StatusForbidden,
	Code:       "not_authorized_edit_participation",
	Message:    "User is not authorized to make changes on participation",
}

var errParticipationNotFound = helios.ErrorAPI{
	StatusCode: http.StatusNotFound,
	Code:       "participation_not_found",
	Message:    "No participation with given id",
}

var errUserNotFound = helios.ErrorAPI{
	StatusCode: http.StatusNotFound,
	Code:       "user_not_found",
	Message:    "No user with given username",
}

var errQuestionChangeNotAuthorized = helios.ErrorAPI{
	StatusCode: http.StatusForbidden,
	Code:       "not_authorized_edit_question",
	Message:    "User is not authorized to make changes on question",
}

var errQuestionNotFound = helios.ErrorAPI{
	StatusCode: http.StatusNotFound,
	Code:       "question_not_found",
	Message:    "No question with given ID",
}

var errAnswerNotValid = helios.ErrorAPI{
	StatusCode: http.StatusBadRequest,
	Code:       "invalid_answer",
	Message:    "Answer is not valid",
}

var errSubmissionNotAuthorized = helios.ErrorAPI{
	StatusCode: http.StatusForbidden,
	Code:       "cannot_submit_submission",
	Message:    "You are not allowed to submit to this question",
}

var errSynchronizationNotAuthorized = helios.ErrorAPI{
	StatusCode: http.StatusForbidden,
	Code:       "cannot_get_synchronization",
	Message:    "You are not allowed to get the synchronziation data",
}

var errDecryptEventForbidden = helios.ErrorAPI{
	StatusCode: http.StatusForbidden,
	Code:       "decrypt_forbidden",
	Message:    "You are not allowed to decrypt the exam",
}

var errDecryptEventFailed = helios.ErrorAPI{
	StatusCode: http.StatusBadRequest,
	Code:       "decrypt_failed",
	Message:    "Failed to decrypt event. Make sure the key given is correct.",
}
