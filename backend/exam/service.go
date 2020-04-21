package exam

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/yonasadiel/charon/backend/auth"
	"github.com/yonasadiel/helios"
)

// GetAllVenue returns all venues.
// Only admin and organizer have permission for this use case.
func GetAllVenue(user auth.User) ([]Venue, helios.Error) {
	if !user.IsAdmin() && !user.IsOrganizer() {
		return nil, errVenueAccessNotAuthorized
	}

	var venues []Venue
	helios.DB.Find(&venues)
	return venues, nil
}

// UpsertVenue creates or updates a venue. It creates if
// ID = 0, or updates otherwise. Only user with organizer or
// admin role that can creates / updates venue.
// If it is create, then venue.ID will be changed.
func UpsertVenue(user auth.User, venue *Venue) helios.Error {
	if !user.IsAdmin() && !user.IsOrganizer() {
		return errVenueAccessNotAuthorized
	}

	if venue.ID == 0 {
		helios.DB.Create(venue)
	} else {
		helios.DB.Save(venue)
	}
	return nil
}

// DeleteVenue deletes a venue with given id
// and returns the deleted venue. Only organizer and
// admin that can do deletion. If there is an event
// organized on the venue, it will fail
func DeleteVenue(user auth.User, venueID uint) (*Venue, helios.Error) {
	if !user.IsOrganizer() && !user.IsAdmin() {
		return nil, errVenueAccessNotAuthorized
	}

	var venue Venue
	var participationCount int

	helios.DB.Where("id = ?", venueID).First(&venue)
	if venue.ID == 0 {
		return nil, errVenueNotFound
	}

	helios.DB.Model(&Participation{}).Where("venue_id = ?", venue.ID).Count(&participationCount)
	if participationCount > 0 {
		return nil, errVenueCantDeletedEventExists
	}

	helios.DB.Delete(&venue)
	return &venue, nil
}

// GetAllEventOfUser returns all events that is participated by user.
// If the user is admin or organizer, then return all events that are exist.
func GetAllEventOfUser(user auth.User) []Event {
	var events []Event

	if user.IsAdmin() || user.IsOrganizer() {
		helios.DB.
			Order("events.starts_at asc").
			Find(&events)
	} else { // user is local or participant
		helios.DB.
			Select("events.*").
			Table("events").
			Joins("inner join participations on participations.event_id = events.id").
			Where("user_id = ?", user.ID).
			Order("events.starts_at asc").
			Find(&events)
	}

	return events
}

// GetEventOfUser returns the event if exist
// If the user is local or participant, the permission will be checked.
func GetEventOfUser(user auth.User, eventSlug string) (Event, helios.Error) {
	var event Event

	if user.IsAdmin() || user.IsOrganizer() {
		helios.DB.Where("slug = ?", eventSlug).First(&event)
	} else {
		helios.DB.
			Table("events").
			Select("events.*").
			Joins("inner join participations on participations.event_id = events.id").
			Where("participations.user_id = ?", user.ID).
			Where("events.slug = ?", eventSlug).
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
		event.SimKey = generateRandomToken(32)
		helios.DB.Omit("last_synchronization").Create(event)
	} else {
		helios.DB.Omit("last_synchronization", "key").Save(event)
	}

	return nil
}

// GetAllParticipationOfUserAndEvent returns all participations of the event.
func GetAllParticipationOfUserAndEvent(user auth.User, eventSlug string) ([]Participation, helios.Error) {
	var event Event
	var participations []Participation
	var errGetEvent helios.Error

	event, errGetEvent = GetEventOfUser(user, eventSlug)
	if errGetEvent != nil {
		return nil, errGetEvent
	}

	helios.DB.
		Table("participations").
		Joins("inner join users on participations.user_id = users.id").
		Preload("User").
		Preload("Venue").
		Where("event_id = ?", event.ID).
		Where("users.role < ?", user.Role).
		Find(&participations)

	return participations, nil
}

// UpsertParticipation creates or updates a participation. Only available to
// user with higher role. If the user is not participate to the event, create new
// participation on the venue. If it has already existed, update the venue.
func UpsertParticipation(user auth.User, eventSlug string, userUsername string, participation *Participation) helios.Error {
	var event Event
	var participationUser auth.User
	var venue Venue
	var errGetEvent helios.Error
	event, errGetEvent = GetEventOfUser(user, eventSlug)
	if errGetEvent != nil {
		return errGetEvent
	}

	helios.DB.Where("username = ?", userUsername).First(&participationUser)
	if participationUser.ID == 0 {
		return errUserNotFound
	} else if participationUser.Role >= user.Role {
		return errParticipationChangeNotAuthorized
	}

	helios.DB.Where("id = ?", participation.VenueID).First(&venue)
	if venue.ID == 0 {
		return errVenueNotFound
	}

	helios.DB.Where("user_id = ?", participationUser.ID).Where("event_id = ?", event.ID).First(&participation)
	participation.User = &participationUser
	participation.UserID = participationUser.ID
	participation.EventID = event.ID
	participation.Event = &event
	participation.VenueID = venue.ID
	participation.Venue = &venue
	if participation.ID == 0 {
		helios.DB.Create(&participation)
	} else {
		helios.DB.Save(&participation)
	}

	return nil
}

// DeleteParticipation deletes a participation with given id
// and returns the deleted participation. Only available to
// user with higher role.
func DeleteParticipation(user auth.User, eventSlug string, participationID uint) (*Participation, helios.Error) {
	var event Event
	var participation Participation
	var errGetEvent helios.Error
	event, errGetEvent = GetEventOfUser(user, eventSlug)
	if errGetEvent != nil {
		return nil, errGetEvent
	}

	helios.DB.Preload("User").Preload("Venue").Where("id = ?", participationID).Where("event_id = ?", event.ID).First(&participation)
	if participation.ID == 0 {
		return nil, errParticipationNotFound
	} else if participation.User.Role >= user.Role {
		return nil, errParticipationChangeNotAuthorized
	}

	tx := helios.DB.Begin()
	tx.Where("participation_id = ?", participationID).Delete(UserQuestion{})
	tx.Delete(&participation)
	tx.Commit()
	return &participation, nil
}

// GetAllQuestionOfUserAndEvent returns all questions in database
// that exists on an event and belongs to an user.
// Current submission of the user will be attached.
func GetAllQuestionOfUserAndEvent(user auth.User, eventSlug string) ([]Question, helios.Error) {
	var event Event
	var questions []Question
	var errGetEvent helios.Error

	event, errGetEvent = GetEventOfUser(user, eventSlug)
	if errGetEvent != nil {
		return nil, errGetEvent
	}

	if !user.IsAdmin() && !user.IsOrganizer() && event.StartsAt.After(time.Now()) {
		return nil, errEventIsNotYetStarted
	}

	// Querying for user questions and user submissions
	if user.IsAdmin() || user.IsOrganizer() || user.IsLocal() {
		helios.DB.Preload("Choices").Where("event_id = ?", event.ID).Find(&questions)
	} else {
		helios.DB.
			Select("questions.*, user_questions.answer as user_answer").
			Table("questions").
			Preload("Choices").
			Joins("inner join user_questions on user_questions.question_id = questions.id").
			Joins("inner join participations on participations.id = user_questions.participation_id").
			Where("questions.event_id = ?", event.ID).
			Where("participations.user_id = ?", user.ID).
			Order("user_questions.ordering asc").
			Find(&questions)
	}

	return questions, nil
}

// UpsertQuestion creates or updates a question. Only available to
// admin and organizer. Notice that the EventID may be changed, so
// this function may move a question to other event.
// If it is updating, all choices will be deleted then recreated.
func UpsertQuestion(user auth.User, eventSlug string, question *Question) helios.Error {
	if !user.IsOrganizer() && !user.IsAdmin() {
		return errQuestionChangeNotAuthorized
	}

	var event Event
	var errGetEvent helios.Error
	event, errGetEvent = GetEventOfUser(user, eventSlug)
	if errGetEvent != nil {
		return errGetEvent
	}

	tx := helios.DB.Begin()
	if question.ID == 0 {
		choices := question.Choices
		question.Choices = []QuestionChoice{}
		question.Event = &event
		tx.Create(question)
		for _, choice := range choices {
			choice.ID = 0
			choice.QuestionID = question.ID
			tx.Create(&choice)
		}
		question.Choices = choices
	} else {
		choices := question.Choices
		question.Choices = []QuestionChoice{}
		question.Event = &event
		tx.Delete(QuestionChoice{}, "question_id = ?", question.ID)
		tx.Save(question)
		for _, choice := range choices {
			choice.ID = 0
			choice.QuestionID = question.ID
			tx.Create(&choice)
		}
		question.Choices = choices
	}
	tx.Commit()

	return nil
}

// GetQuestionOfEventAndUser returns a question with given id, but first check
// if the user has rights to the question
func GetQuestionOfEventAndUser(user auth.User, eventSlug string, questionID uint) (*Question, helios.Error) {
	var event Event
	var question Question
	var errGetEvent helios.Error

	event, errGetEvent = GetEventOfUser(user, eventSlug)
	if errGetEvent != nil {
		return nil, errGetEvent
	}

	if !user.IsAdmin() && !user.IsOrganizer() && event.StartsAt.After(time.Now()) {
		return nil, errEventIsNotYetStarted
	}

	if user.IsAdmin() || user.IsOrganizer() || user.IsLocal() {
		helios.DB.Where("id = ?", questionID).Where("event_id = ?", event.ID).First(&question)
	} else {
		helios.DB.
			Select("questions.*, user_questions.answer as user_answer").
			Table("questions").
			Preload("Choices").
			Joins("inner join user_questions on user_questions.question_id = questions.id").
			Joins("inner join participations on participations.id = user_questions.participation_id").
			Where("questions.event_id = ?", event.ID).
			Where("participations.user_id = ?", user.ID).
			Where("participations.event_id = ?", event.ID).
			Where("questions.id = ?", questionID).
			First(&question)
	}
	if question.ID == 0 {
		return nil, errQuestionNotFound
	}
	return &question, nil
}

// DeleteQuestion deletes a question with given id
// and returns the deleted question. Only organizer and
// admin that can do deletion
func DeleteQuestion(user auth.User, eventSlug string, questionID uint) (*Question, helios.Error) {
	if !user.IsOrganizer() && !user.IsAdmin() {
		return nil, errQuestionChangeNotAuthorized
	}

	var event Event
	var question Question
	var errGetEvent helios.Error

	event, errGetEvent = GetEventOfUser(user, eventSlug)
	if errGetEvent != nil {
		return nil, errGetEvent
	}

	helios.DB.Where("event_id = ?", event.ID).Where("id = ?", questionID).First(&question)
	if question.ID == 0 {
		return nil, errQuestionNotFound
	}
	tx := helios.DB.Begin()
	tx.Where("question_id = ?", questionID).Delete(UserQuestion{})
	tx.Where("question_id = ?", questionID).Delete(QuestionChoice{})
	tx.Delete(&question)
	tx.Commit()
	return &question, nil
}

// SubmitSubmission submit a submission from user to a question.
func SubmitSubmission(user auth.User, eventSlug string, questionID uint, answer string) (*Question, helios.Error) {
	if !user.IsParticipant() {
		return nil, errSubmissionNotAuthorized
	}

	var event Event
	var userQuestion UserQuestion
	var errGetEvent helios.Error

	event, errGetEvent = GetEventOfUser(user, eventSlug)
	if errGetEvent != nil {
		return nil, errGetEvent
	}

	if !user.IsAdmin() && !user.IsOrganizer() && event.StartsAt.After(time.Now()) {
		return nil, errEventIsNotYetStarted
	}

	helios.DB.
		Select("user_questions.*").
		Table("user_questions").
		Preload("Question").
		Preload("Question.Choices").
		Joins("inner join questions on questions.id = user_questions.question_id").
		Joins("inner join participations on participations.id = user_questions.participation_id").
		Where("questions.event_id = ?", event.ID).
		Where("participations.user_id = ?", user.ID).
		Where("participations.event_id = ?", event.ID).
		Where("questions.id = ?", questionID).
		First(&userQuestion)

	if userQuestion.ID == 0 {
		return nil, errQuestionNotFound
	}

	if !isAnswerValidChoice(answer, userQuestion.Question.Choices) {
		return nil, errAnswerNotValid
	}

	userQuestion.Answer = answer
	userQuestion.Question.UserAnswer = answer
	helios.DB.Save(&userQuestion)
	return userQuestion.Question, nil
}

// GetSynchronizationData gets the synchronization data of event.
// Only local user has the permission
func GetSynchronizationData(user auth.User, eventSlug string) (*Event, *Venue, []Question, []auth.User, helios.Error) {
	if !user.IsLocal() {
		return nil, nil, nil, nil, errSynchronizationNotAuthorized
	}

	var participation Participation
	helios.DB.
		Table("participations").
		Select("participations.*").
		Preload("Venue").
		Joins("inner join events on events.id = participations.event_id").
		Where("participations.user_id = ?", user.ID).
		Where("events.slug = ?", eventSlug).
		First(&participation)
	if participation.ID == 0 {
		return nil, nil, nil, nil, errEventNotFound
	}

	var event Event
	var questions []Question
	var users []auth.User
	helios.DB.Where("id = ?", participation.EventID).First(&event)
	helios.DB.Preload("Choices").Where("event_id = ?", event.ID).Find(&questions)
	helios.DB.
		Select("users.*").
		Joins("inner join participations on participations.user_id = users.id").
		Where("participations.event_id = ?", event.ID).
		Where("participations.venue_id = ?", participation.Venue.ID).
		Find(&users)

	err := encryptQuestions(questions, event.SimKey)
	if err != nil {
		return nil, nil, nil, nil, helios.ErrInternalServerError
	}

	return &event, participation.Venue, questions, users, nil
}

// PutSynchronizationData puts the synchronization data of event.
// Only local user has the permission
func PutSynchronizationData(user auth.User, event Event, venue Venue, questions []Question, users []auth.User) helios.Error {
	if !user.IsLocal() {
		return errSynchronizationNotAuthorized
	}

	var eventSaved Event
	var userParticipation Participation

	tx := helios.DB.Begin()
	tx.Create(&venue)

	// Update or create event and user participation
	tx.Where("slug = ?", event.Slug).First(&eventSaved)
	event.LastSynchronization = time.Now()
	if eventSaved.ID == 0 {
		tx.Create(&event)
	} else {
		event.ID = eventSaved.ID
		tx.Save(&event)
	}

	userParticipation = Participation{
		UserID:  user.ID,
		VenueID: venue.ID,
		EventID: event.ID,
	}
	tx.Create(&userParticipation)

	// update or create user
	for _, user := range users {
		var userSaved auth.User
		tx.Where("username = ?", user.Username).First(&userSaved)
		user.Role = auth.UserRoleParticipant
		if userSaved.ID == 0 {
			tx.Create(&user)
		} else {
			user.ID = userSaved.ID
			tx.Save(&user)
		}
	}
	// reset all questions and particpations
	tx.Delete(Question{}, "event_id = ?", event.ID)
	tx.Delete(Participation{}, "event_id = ?", event.ID)
	// create all questions and participations
	for i := range questions {
		questions[i].ID = 0
		questions[i].Event = &event
		questions[i].EventID = event.ID
		questions[i].Choices = []QuestionChoice{}
		tx.Create(&questions[i])
	}
	for i := range users {
		var participation Participation = Participation{
			UserID:  users[i].ID,
			VenueID: venue.ID,
			EventID: event.ID,
		}
		tx.Create(&participation)
	}
	tx.Commit()
	return nil
}

func encryptQuestions(questions []Question, encryptionKey string) error {
	for i := range questions {
		encryptedContent, err := encryptToBase64(encryptionKey, questions[i].Content)
		if err != nil {
			return err
		}
		questions[i].Content = encryptedContent
		for j := range questions[i].Choices {
			encryptedChoice, err := encryptToBase64(encryptionKey, questions[i].Choices[j].Text)
			if err != nil {
				return err
			}
			questions[i].Choices[j].Text = encryptedChoice
		}
	}
	return nil
}

func decryptQuestions(questions []Question, decryptionKey string) error {
	fmt.Println(len(questions))
	for i := range questions {
		decryptedContent, err := decryptFromBase64(decryptionKey, questions[i].Content)
		if err != nil {
			return err
		}
		questions[i].Content = decryptedContent
		for j := range questions[i].Choices {
			decryptedChoice, err := decryptFromBase64(decryptionKey, questions[i].Choices[j].Text)
			if err != nil {
				return err
			}
			questions[i].Choices[j].Text = decryptedChoice
		}
	}
	return nil
}

// https://golang.org/pkg/crypto/cipher/#example_NewCFBEncrypter
func encrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func encryptToBase64(key, plaintext string) (string, error) {
	var encryptedBytes []byte
	var err error
	encryptedBytes, err = encrypt([]byte(key), []byte(plaintext))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

// https://golang.org/pkg/crypto/cipher/#example_NewCFBDecrypter
func decrypt(key []byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	plaintext := make([]byte, len(ciphertext)-aes.BlockSize)
	iv := ciphertext[:aes.BlockSize]
	copy(plaintext, ciphertext[aes.BlockSize:])

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(plaintext, plaintext)
	return plaintext, nil
}

func decryptFromBase64(key, ciphertext string) (string, error) {
	var encryptedBytes, decryptedBytes []byte
	var err error
	encryptedBytes, err = base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	decryptedBytes, err = decrypt([]byte(key), encryptedBytes)
	if err != nil {
		return "", err
	}
	return string(decryptedBytes), nil
}
