package exam

import (
	"github.com/yonasadiel/charon/auth"
	"github.com/yonasadiel/helios"
)

// GetAllEventOfUser returns all events that is participated by user.
// If the user is admin or organizer, then return all events that are exist.
func GetAllEventOfUser(user auth.User) []Event {
	var events []Event

	if user.IsAdmin() || user.IsOrganizer() {
		helios.DB.Find(&events)
	} else { // user is local or participant
		helios.DB.
			Table("user_events").
			Select("events.*").
			Joins("left join events on events.id = user_events.event_id").
			Where("user_id = ?", user.ID).
			Order("events.starts_at").
			Find(&events)
	}

	return events
}

// GetEventOfUser returns the event if exist
// If the user is local or participant, the permission will be checked.
func GetEventOfUser(user auth.User, eventID uint) (Event, helios.Error) {
	var event Event

	if user.IsAdmin() || user.IsOrganizer() {
		helios.DB.Where("id = ?", eventID).First(&event)
	} else {
		helios.DB.
			Table("user_events").
			Select("events.*").
			Joins("left join events on events.id = user_events.event_id").
			Where("user_id = ?", user.ID).
			Where("events.id = ?", eventID).
			Order("events.starts_at").
			First(&event)
	}
	if event.ID == 0 {
		return event, errEventNotFound
	}
	return event, nil
}

// UpsertEvent creates or updates an exam event. It creates if
// ID = 0, or updates otherwise. Only user with organizer or
// admin role that can creates / updates event.
// If it is create, then event.ID will be changed.
func UpsertEvent(user auth.User, event *Event) helios.Error {
	if !user.IsOrganizer() && !user.IsAdmin() {
		return errEventChangeNotAuthorized
	}

	if event.ID == 0 {
		helios.DB.Create(event)
	} else {
		helios.DB.Save(event)
	}

	return nil
}

// GetAllQuestionOfEventAndUser returns all questions in database
// that exists on an event and belongs to an user.
// Current submission of the user will be attached.
func GetAllQuestionOfEventAndUser(user auth.User, eventID uint) ([]Question, helios.Error) {
	var event Event
	var questions []Question
	var userSubmissions []Submission
	var userSubmissionByQuestionID = make(map[uint]Submission)
	var errGetEvent helios.Error

	event, errGetEvent = GetEventOfUser(user, eventID)
	if errGetEvent != nil {
		return nil, errGetEvent
	}

	// Querying for user questions and user submissions
	if user.IsAdmin() || user.IsOrganizer() {
		helios.DB.Preload("Choices").Where("event_id = ?", eventID).Find(&questions)
	} else {
		helios.DB.
			Table("user_questions").
			Preload("Choices").
			Select("questions.*").
			Joins("left join questions on questions.id = user_questions.question_id").
			Where("user_id = ?", user.ID).
			Where("questions.event_id = ?", event.ID).
			Order("user_questions.ordering asc").
			Find(&questions)
		helios.DB.Where("user_id = ?", user.ID).Order("created_at asc").Find(&userSubmissions)
	}

	// Mapping user submission by the question id
	for _, userSubmission := range userSubmissions {
		// Here, we safely assume that if the userSubmission is latter
		// in array, then it has latest creation time
		userSubmissionByQuestionID[userSubmission.QuestionID] = userSubmission
	}

	// Set the question answer to the user submission
	for i := range questions {
		userSubmission := userSubmissionByQuestionID[questions[i].ID]
		if userSubmission.ID != 0 {
			questions[i].UserAnswer = userSubmission.Answer
		}
	}

	return questions, nil
}

// UpsertQuestion creates or updates a question. Only available to
// admin and organizer. Notice that the EventID may be changed, so
// this function may move a question to other event.
// If it is updating, all choices will be deleted then recreated.
func UpsertQuestion(user auth.User, question *Question) helios.Error {
	if !user.IsOrganizer() && !user.IsAdmin() {
		return errQuestionChangeNotAuthorized
	}

	var event Event
	helios.DB.Where("id = ?", question.EventID).First(&event)
	if event.ID == 0 {
		return errEventNotFound
	}

	if question.ID == 0 {
		choices := question.Choices
		question.Choices = []QuestionChoice{}
		tx := helios.DB.Begin()
		tx.Create(question)
		for _, choice := range choices {
			choice.ID = 0
			choice.QuestionID = question.ID
			tx.Create(&choice)
		}
		question.Choices = choices
		tx.Commit()
	} else {
		choices := question.Choices
		question.Choices = []QuestionChoice{}
		tx := helios.DB.Begin()
		tx.Delete(QuestionChoice{}, "question_id = ?", question.ID)
		tx.Save(question)
		for _, choice := range choices {
			choice.ID = 0
			choice.QuestionID = question.ID
			tx.Create(&choice)
		}
		question.Choices = choices
		tx.Commit()
	}

	return nil
}

// GetQuestionOfUser returns a question with given id, but first check
// if the user has rights to the question
func GetQuestionOfUser(user auth.User, eventID uint, questionID uint) (*Question, helios.Error) {
	var event Event
	var question Question
	var userSubmission Submission
	var errGetEvent helios.Error

	event, errGetEvent = GetEventOfUser(user, eventID)
	if errGetEvent != nil {
		return nil, errGetEvent
	}

	helios.DB.
		Table("user_questions").
		Select("questions.*").
		Joins("left join questions on questions.id = user_questions.question_id").
		Where("user_id = ?", user.ID).
		Where("event_id = ?", event.ID).
		Where("question_id = ?", questionID).
		First(&question)
	if question.ID == 0 {
		return nil, errQuestionNotFound
	}

	helios.DB.Where("user_id = ?", user.ID).Where("question_id = ?", questionID).Order("created_at desc").First(&userSubmission)
	if userSubmission.ID != 0 {
		question.UserAnswer = userSubmission.Answer
	}
	return &question, nil
}

// DeleteQuestion deletes a question with given id
// and returns the deleted question. Only organizer and
// admin that can do deletion
func DeleteQuestion(user auth.User, eventID uint, questionID uint) (*Question, helios.Error) {
	if !user.IsOrganizer() && !user.IsAdmin() {
		return nil, errQuestionChangeNotAuthorized
	}

	var event Event
	var question Question
	var errGetEvent helios.Error

	event, errGetEvent = GetEventOfUser(user, eventID)
	if errGetEvent != nil {
		return nil, errGetEvent
	}

	helios.DB.Where("event_id = ?", event.ID).Where("id = ?", questionID).First(&question)
	if question.ID == 0 {
		return nil, errQuestionNotFound
	}
	tx := helios.DB.Begin()
	tx.Delete(&question)
	tx.Where("question_id = ?", questionID).Delete(Submission{})
	tx.Where("question_id = ?", questionID).Delete(QuestionChoice{})
	tx.Commit()
	return &question, nil
}

// SubmitSubmission submit a submission from user to a question.
func SubmitSubmission(user auth.User, eventID uint, questionID uint, answer string) (*Submission, helios.Error) {
	var event Event
	var question Question
	var choices []QuestionChoice
	var submission Submission
	var errGetEvent helios.Error

	event, errGetEvent = GetEventOfUser(user, eventID)
	if errGetEvent != nil {
		return nil, errGetEvent
	}
	helios.DB.Where("event_id = ?", event.ID).Where("id = ?", questionID).First(&question)
	helios.DB.Where("question_id = ?", questionID).Find(&choices)
	if question.ID == 0 {
		return nil, errQuestionNotFound
	}
	if !isAnswerValidChoice(answer, choices) {
		return nil, errAnswerNotValid
	}
	submission = Submission{
		Answer:     answer,
		QuestionID: question.ID,
		UserID:     user.ID,
		Question:   &question,
		User:       &user,
	}
	helios.DB.Create(&submission)
	return &submission, nil
}
