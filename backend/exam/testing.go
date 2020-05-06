package exam

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/yonasadiel/charon/backend/auth"
	"github.com/yonasadiel/helios"
)

var eventSeq uint = 0
var venueSeq uint = 0
var participationSeq uint = 0
var questionSeq uint = 0
var questionChoiceSeq uint = 0
var userQuestionSeq uint = 0

// EventFactory creates an event for testing. The given argument will be
// completed if the attribute is empty.
func EventFactory(event Event) Event {
	eventSeq = eventSeq + 1
	if event.Slug == "" {
		event.Slug = fmt.Sprintf("slug-event-%d", eventSeq)
	}
	if event.Description == "" {
		event.Description = fmt.Sprintf("even desc %d", eventSeq)
	}
	if event.Title == "" {
		event.Title = fmt.Sprintf("Title Event %d", eventSeq)
	}
	if event.StartsAt.IsZero() {
		event.StartsAt = time.Now().Add(-1 * time.Hour)
	}
	if event.EndsAt.IsZero() {
		event.EndsAt = event.StartsAt.Add(2 * time.Hour)
	}
	if event.SimKey == "" && event.PubKey == "" && event.PrvKey == "" && event.SimKeySign == "" {
		var prvKey *rsa.PrivateKey
		var simKeySign []byte
		event.SimKey = generateRandomToken(32)
		prvKey, _ = rsa.GenerateKey(rand.Reader, 1024)
		event.PrvKey = base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(prvKey))
		event.PubKey = base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(&prvKey.PublicKey))
		simKeyHashed := sha256.Sum256([]byte(event.SimKey))
		simKeySign, _ = rsa.SignPSS(rand.Reader, prvKey, crypto.SHA256, simKeyHashed[:], nil)
		event.SimKeySign = base64.StdEncoding.EncodeToString(simKeySign)
	}
	if event.DecryptedAt.IsZero() {
		event.DecryptedAt = event.StartsAt.Add(-1 * time.Hour)
	}
	return event
}

// EventFactorySaved do exactly like EventFactory but the result
// will be saved to database
func EventFactorySaved(event Event) Event {
	if event.ID == 0 {
		event = EventFactory(event)
		helios.DB.Create(&event)
	}
	return event
}

// VenueFactory creates a venue for testing. The given argument will be
// completed if the attribute is empty.
func VenueFactory(venue Venue) Venue {
	venueSeq = venueSeq + 1
	if venue.Name == "" {
		venue.Name = fmt.Sprintf("Venue %d", venueSeq)
	}
	return venue
}

// VenueFactorySaved do exactly like VenueFactory but the result
// will be saved to database
func VenueFactorySaved(venue Venue) Venue {
	if venue.ID == 0 {
		venue = VenueFactory(venue)
		helios.DB.Create(&venue)
	}
	return venue
}

// ParticipationFactory creates a participation for testing. The given argument will be
// completed if the attribute is empty.
func ParticipationFactory(participation Participation) Participation {
	participationSeq = participationSeq + 1
	if participation.Venue == nil && participation.VenueID == 0 {
		venue := VenueFactory(Venue{})
		participation.Venue = &venue
		participation.VenueID = venue.ID
	}
	if participation.Event == nil && participation.EventID == 0 {
		event := EventFactory(Event{})
		participation.Event = &event
	}
	if participation.User == nil && participation.UserID == 0 {
		user := auth.UserFactory(auth.User{})
		participation.User = &user
	}
	return participation
}

// ParticipationFactorySaved do exactly like ParticipationFactory but the result
// will be saved to database
func ParticipationFactorySaved(participation Participation) Participation {
	if participation.ID == 0 {
		participation = ParticipationFactory(participation)
		var venue Venue = VenueFactorySaved(*participation.Venue)
		var event Event = EventFactorySaved(*participation.Event)
		var user auth.User = auth.UserFactorySaved(*participation.User)
		participation.VenueID = venue.ID
		participation.EventID = event.ID
		participation.UserID = user.ID
		participation.Venue = nil
		participation.Event = nil
		participation.User = nil
		helios.DB.Create(&participation)
		participation.Venue = &venue
		participation.Event = &event
		participation.User = &user
	}
	return participation
}

// QuestionFactory creates a question for testing. The given argument will be
// completed if the attribute is empty.
func QuestionFactory(question Question) Question {
	questionSeq = questionSeq + 1
	if question.Content == "" {
		question.Content = fmt.Sprintf("Question content #%d", questionSeq+1)
	}
	if question.Event == nil && question.EventID == 0 {
		event := EventFactory(Event{})
		question.Event = &event
	}
	if question.Choices == "" {
		var choices []string
		for i := 0; i < 4; i++ {
			choices = append(choices, fmt.Sprintf("choice%d.%d", questionSeq, i+1))
		}
		question.Choices = strings.Join(choices, "|")
	}
	return question
}

// QuestionFactorySaved do exactly like QuestionFactory but the result
// will be saved to database
func QuestionFactorySaved(question Question) Question {
	if question.ID == 0 {
		question = QuestionFactory(question)
		helios.DB.Create(&question)
	}
	return question
}

// UserQuestionFactory creates a user question for testing. The given argument will be
// completed if the attribute is empty.
func UserQuestionFactory(userQuestion UserQuestion) UserQuestion {
	userQuestionSeq = userQuestionSeq + 1
	if userQuestion.Participation == nil && userQuestion.ParticipationID == 0 {
		participation := ParticipationFactory(Participation{})
		userQuestion.Participation = &participation
	}
	if userQuestion.Question == nil && userQuestion.QuestionID == 0 {
		question := QuestionFactory(Question{Event: userQuestion.Participation.Event})
		userQuestion.Question = &question
	}
	if userQuestion.Answer == "" {
		userQuestion.Answer = fmt.Sprintf("Answer %d", userQuestionSeq)
	}
	if userQuestion.Ordering == 0 {
		userQuestion.Ordering = userQuestionSeq
	}
	return userQuestion
}

// UserQuestionFactorySaved do exactly like UserQuestionFactory but the result
// will be saved to database
func UserQuestionFactorySaved(userQuestion UserQuestion) UserQuestion {
	if userQuestion.ID == 0 {
		userQuestion = UserQuestionFactory(userQuestion)
		var question Question = QuestionFactorySaved(*userQuestion.Question)
		var participation Participation = ParticipationFactorySaved(*userQuestion.Participation)
		userQuestion.QuestionID = question.ID
		userQuestion.ParticipationID = participation.ID
		userQuestion.Question = nil
		userQuestion.Participation = nil
		helios.DB.Create(&userQuestion)
		userQuestion.Question = &question
		userQuestion.Participation = &participation
	}
	return userQuestion
}
