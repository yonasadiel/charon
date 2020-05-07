package exam

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
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
	if !user.IsOrganizer() && !user.IsAdmin() && !user.IsLocal() {
		return errEventChangeNotAuthorized
	}

	if event.ID == 0 {
		if user.IsAdmin() || user.IsOrganizer() {
			var err error
			var prvKey *rsa.PrivateKey
			var simKeySign []byte
			prvKey, err = rsa.GenerateKey(rand.Reader, 1024)
			if err != nil {
				return helios.ErrInternalServerError
			}
			event.SimKey = generateRandomToken(32)
			event.PrvKey = base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(prvKey))
			event.PubKey = base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(&prvKey.PublicKey))
			simKeyHashed := sha256.Sum256([]byte(event.SimKey))
			simKeySign, err = rsa.SignPSS(rand.Reader, prvKey, crypto.SHA256, simKeyHashed[:], nil)
			if err != nil {
				return helios.ErrInternalServerError
			}
			event.SimKeySign = base64.StdEncoding.EncodeToString(simKeySign)
		}
		helios.DB.Omit("last_synchronization").Create(event)
		if user.IsLocal() {
			helios.DB.Create(&Participation{
				User:  &user,
				Event: event,
			})
		}
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
		Where("(users.role < ? or users.id = ?)", user.Role, user.ID).
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
	var participationSaved Participation
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

	helios.DB.Where("user_id = ?", participationUser.ID).Where("event_id = ?", event.ID).First(&participationSaved)
	participation.ID = participationSaved.ID
	participation.User = &participationUser
	participation.UserID = participationUser.ID
	participation.EventID = event.ID
	participation.Event = &event
	participation.VenueID = venue.ID
	participation.Venue = &venue
	participation.KeyHashedOnce = fmt.Sprintf("%x", sha256.Sum256([]byte(participation.KeyPlain)))
	participation.KeyHashedTwice = fmt.Sprintf("%x", sha256.Sum256([]byte(participation.KeyHashedOnce)))
	if participation.ID == 0 {
		helios.DB.Create(&participation)
	} else {
		helios.DB.Save(&participation)
	}

	return nil
}

// VerifyParticipation checks if the hashedOnce equal to the participation key
func VerifyParticipation(user auth.User, eventSlug string, hashedOnce string) helios.Error {
	var event Event
	var participation Participation
	var errGetEvent helios.Error
	event, errGetEvent = GetEventOfUser(user, eventSlug)
	if errGetEvent != nil {
		return errGetEvent
	}
	helios.DB.Where("user_id = ?", user.ID).Where("event_id = ?", event.ID).First(&participation)
	hashedTwice := fmt.Sprintf("%x", sha256.Sum256([]byte(hashedOnce)))
	if participation.KeyHashedTwice == hashedTwice {
		participation.KeyHashedOnce = hashedOnce
		helios.DB.Save(&participation)
		return nil
	}
	return errParticipationWrongKey
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
		helios.DB.Where("event_id = ?", event.ID).Order("questions.id asc").Find(&questions)
	} else {
		helios.DB.
			Select("questions.*, user_questions.answer as user_answer").
			Table("questions").
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
	question.Event = &event
	// TODO: make sure all choices have the same length
	if question.ID == 0 {
		tx.Create(question)
	} else {
		tx.Save(question)
	}
	tx.Commit()

	return nil
}

// GetQuestionOfEventAndUser returns a question with given id, but first check
// if the user has rights to the question
func GetQuestionOfEventAndUser(user auth.User, eventSlug string, questionNumber uint) (*Question, helios.Error) {
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
		helios.DB.
			Where("event_id = ?", event.ID).
			Order("questions.id asc").
			Offset(questionNumber - 1).
			First(&question)
	} else {
		helios.DB.
			Select("questions.*, user_questions.answer as user_answer").
			Table("questions").
			Joins("inner join user_questions on user_questions.question_id = questions.id").
			Joins("inner join participations on participations.id = user_questions.participation_id").
			Where("questions.event_id = ?", event.ID).
			Where("participations.user_id = ?", user.ID).
			Where("participations.event_id = ?", event.ID).
			Order("user_questions.ordering asc").
			Offset(questionNumber - 1).
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
func DeleteQuestion(user auth.User, eventSlug string, questionNumber uint) (*Question, helios.Error) {
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

	helios.DB.
		Where("event_id = ?", event.ID).
		Order("questions.id asc").
		Offset(questionNumber - 1).
		First(&question)
	if question.ID == 0 {
		return nil, errQuestionNotFound
	}
	tx := helios.DB.Begin()
	tx.Where("question_id = ?", question.ID).Delete(UserQuestion{})
	tx.Delete(&question)
	tx.Commit()
	return &question, nil
}

// SubmitSubmission submit a submission from user to a question.
func SubmitSubmission(user auth.User, eventSlug string, questionNumber uint, answer string) (*Question, helios.Error) {
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
		Joins("inner join questions on questions.id = user_questions.question_id").
		Joins("inner join participations on participations.id = user_questions.participation_id").
		Where("questions.event_id = ?", event.ID).
		Where("participations.user_id = ?", user.ID).
		Where("participations.event_id = ?", event.ID).
		Order("user_questions.ordering asc").
		Offset(questionNumber - 1).
		First(&userQuestion)

	if userQuestion.ID == 0 {
		return nil, errQuestionNotFound
	}

	userQuestion.Answer = answer
	userQuestion.Question.UserAnswer = answer
	helios.DB.Save(&userQuestion)
	userQuestion.Question.ID = questionNumber
	return userQuestion.Question, nil
}

// GetParticipationStatus returns status of all participants
func GetParticipationStatus(user auth.User, eventSlug string) ([]ParticipationStatus, helios.Error) {
	if !user.IsLocal() {
		return nil, errParticipationStatusAccessNotAuthorized
	}
	var event Event
	var errGetEvent helios.Error
	event, errGetEvent = GetEventOfUser(user, eventSlug)
	if errGetEvent != nil {
		return nil, errGetEvent
	}

	var status []ParticipationStatus
	helios.DB.
		Select("users.username as user_username, sessions.ip_address, sessions.created_at as login_at, sessions.id as session_id, users.session_locked as user_session_locked").
		Table("participations").
		Joins("left join users on (users.id = participations.user_id and users.deleted_at is null)").
		Joins("left join sessions on (sessions.user_id = users.id and sessions.deleted_at is null)").
		Where("event_id = ?", event.ID).
		Where("users.role = ?", auth.UserRoleParticipant).
		Where("participations.deleted_at is null").
		Find(&status)
	return status, nil
}

// RemoveParticipationSession removes session to force user logout
func RemoveParticipationSession(user auth.User, eventSlug string, sessionID uint) helios.Error {
	if !user.IsLocal() {
		return errParticipationStatusAccessNotAuthorized
	}
	var event Event
	var errGetEvent helios.Error
	event, errGetEvent = GetEventOfUser(user, eventSlug)
	if errGetEvent != nil {
		return errGetEvent
	}
	var session auth.Session
	helios.DB.
		Select("sessions.*").
		Table("participations").
		Joins("left join users on (users.id = participations.user_id and users.deleted_at is null)").
		Joins("left join sessions on (sessions.user_id = users.id and sessions.deleted_at is null)").
		Where("event_id = ?", event.ID).
		Where("users.role = ?", auth.UserRoleParticipant).
		Where("sessions.id = ?", sessionID).
		Where("sessions.deleted_at is null").
		Where("participations.deleted_at is null").
		First(&session)
	if session.ID == 0 {
		return errParticipationStatusNotFound
	}
	helios.DB.Delete(auth.Session{}, "id = ?", session.ID)
	return nil
}

// GetSynchronizationData gets the synchronization data of event.
// Only local user has the permission
func GetSynchronizationData(user auth.User, eventSlug string) (*Event, *Venue, []Question, []auth.User, map[string]string, helios.Error) {
	if !user.IsLocal() {
		return nil, nil, nil, nil, nil, errSynchronizationNotAuthorized
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
		return nil, nil, nil, nil, nil, errEventNotFound
	}

	var event Event
	var questions []Question
	var users []auth.User
	var participations []Participation
	var usersKey map[string]string
	helios.DB.Where("id = ?", participation.EventID).First(&event)
	helios.DB.Where("event_id = ?", event.ID).Find(&questions)
	helios.DB.
		Select("users.*").
		Joins("inner join participations on participations.user_id = users.id").
		Where("participations.event_id = ?", event.ID).
		Where("participations.venue_id = ?", participation.Venue.ID).
		Find(&users)
	helios.DB.
		Select("participations.*").
		Preload("User").
		Where("participations.event_id = ?", event.ID).
		Where("participations.venue_id = ?", participation.Venue.ID).
		Find(&participations)

	err := encryptQuestions(questions, event.SimKey)
	if err != nil {
		return nil, nil, nil, nil, nil, helios.ErrInternalServerError
	}

	usersKey = make(map[string]string)
	for _, participation := range participations {
		usersKey[participation.User.Username] = participation.KeyHashedTwice
	}
	event.SimKey = ""

	return &event, participation.Venue, questions, users, usersKey, nil
}

// PutSynchronizationData puts the synchronization data of event.
// Only local user has the permission
func PutSynchronizationData(user auth.User, event Event, venue Venue, questions []Question, users []auth.User, usersKey map[string]string) helios.Error {
	if !user.IsLocal() {
		return errSynchronizationNotAuthorized
	}

	var eventSaved Event
	var userParticipation Participation

	tx := helios.DB.Begin()
	venue.ID = 0
	tx.Create(&venue)

	// Update or create event and user participation
	tx.Where("slug = ?", event.Slug).First(&eventSaved)
	event.LastSynchronization = time.Now()
	if eventSaved.ID == 0 {
		event.ID = 0
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
	for i := range users {
		var userSaved auth.User
		tx.Where("username = ?", users[i].Username).First(&userSaved)
		if userSaved.ID == 0 {
			users[i].ID = 0
			tx.Create(&users[i])
		} else {
			users[i].ID = userSaved.ID
			tx.Save(&users[i])
		}
	}
	// reset all questions and particpations
	tx.Exec(`UPDATE user_questions
		SET deleted_at = time('now')
		WHERE id IN (
			SELECT user_questions.id
			FROM user_questions
			INNER JOIN questions ON questions.id = user_questions.question_id
			WHERE questions.event_id = ?
		)`, event.ID)
	tx.Delete(Question{}, "event_id = ?", event.ID)
	tx.Delete(Participation{}, "event_id = ?", event.ID)
	// create all questions and participations
	for i := range questions {
		questions[i].ID = 0
		questions[i].Event = &event
		questions[i].EventID = event.ID
		tx.Create(&questions[i])
	}
	for i := range users {
		var participation Participation = Participation{
			UserID:         users[i].ID,
			VenueID:        venue.ID,
			EventID:        event.ID,
			KeyHashedTwice: usersKey[users[i].Username],
			// TODO: if the key is malformed and missing user
		}
		tx.Create(&participation)
		for j := range questions {
			// TODO: send userquestion instead of all question
			tx.Create(&UserQuestion{
				ParticipationID: participation.ID,
				QuestionID:      questions[j].ID,
				Ordering:        uint((j + 1) * 10),
			})
		}
	}
	tx.Commit()
	return nil
}

// DecryptEventData decrypts all event data that is encrypted on synchronization data
func DecryptEventData(user auth.User, eventSlug string, simKey string) helios.Error {
	if !user.IsLocal() {
		return errDecryptEventForbidden
	}
	var event Event
	var errGetEvent helios.Error
	event, errGetEvent = GetEventOfUser(user, eventSlug)
	if errGetEvent != nil {
		return errGetEvent
	}
	if !event.DecryptedAt.IsZero() {
		// Already decrypted
		return nil
	}

	var simKeySign, pubKeyMarshalled []byte
	var err error
	var pubKey *rsa.PublicKey
	simKeySign, err = base64.StdEncoding.DecodeString(event.SimKeySign)
	if err != nil {
		return helios.ErrInternalServerError
	}
	simKeyHashed := sha256.Sum256([]byte(simKey))
	pubKeyMarshalled, err = base64.StdEncoding.DecodeString(event.PubKey)
	if err != nil {
		return helios.ErrInternalServerError
	}
	pubKey, err = x509.ParsePKCS1PublicKey(pubKeyMarshalled)
	if err != nil {
		return helios.ErrInternalServerError
	}
	err = rsa.VerifyPSS(pubKey, crypto.SHA256, simKeyHashed[:], simKeySign, nil)
	if err != nil {
		return errDecryptEventFailed
	}

	var questions []Question
	tx := helios.DB.Begin()
	event.DecryptedAt = time.Now()
	event.SimKey = simKey
	tx.Save(&event)
	tx.Where("event_id = ?", event.ID).Find(&questions)
	err = decryptQuestions(questions, simKey)
	if err != nil {
		tx.Rollback()
		return helios.ErrInternalServerError
	}
	for _, question := range questions {
		tx.Save(&question)
	}
	tx.Commit()
	return nil
}

func encryptQuestions(questions []Question, encryptionKey string) error {
	for i := range questions {
		var encryptedContent, encryptedChoices string
		var err error
		encryptedContent, err = encryptToBase64(encryptionKey, questions[i].Content)
		if err != nil {
			return err
		}
		questions[i].Content = encryptedContent
		encryptedChoices, err = encryptToBase64(encryptionKey, questions[i].Choices)
		if err != nil {
			return err
		}
		questions[i].Choices = encryptedChoices
	}
	return nil
}

func decryptQuestions(questions []Question, decryptionKey string) error {
	for i := range questions {
		var decryptedContent, decryptedChoices string
		var err error
		decryptedContent, err = decryptFromBase64(decryptionKey, questions[i].Content)
		if err != nil {
			return err
		}
		questions[i].Content = decryptedContent
		decryptedChoices, err = decryptFromBase64(decryptionKey, questions[i].Choices)
		if err != nil {
			return err
		}
		questions[i].Choices = decryptedChoices
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
