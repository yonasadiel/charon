package exam

import (
	"time"

	"github.com/yonasadiel/charon/backend/auth"
	"github.com/yonasadiel/helios"
)

// Event is the exam event
// It also stores the start and end time
type Event struct {
	ID                  uint   `gorm:"primary_key"`
	Slug                string `gorm:"size:100;unique"`
	Title               string `gorm:"size:256"`
	Description         string `gorm:"type:text"`
	SimKey              string `gorm:"size:48"`
	SimKeySign          string `gorm:"size:1024"`
	PrvKey              string `gorm:"size:1024"`
	PubKey              string `gorm:"size:1024"`
	DecryptedAt         time.Time
	LastSynchronization time.Time
	StartsAt            time.Time
	EndsAt              time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Venue is the event venue
type Venue struct {
	ID   uint `gorm:"primary_key"`
	Name string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Participation is many to many indicating an user is participating
// in a local event.
type Participation struct {
	ID              uint `gorm:"primary_key"`
	EventID         uint
	UserID          uint
	VenueID         uint
	KeyPlain        string
	KeyHashedOnce string
	KeyHashedTwice string

	Event *Event     `gorm:"foreignkey:EventID;association_autoupdate:false"`
	User  *auth.User `gorm:"foreignkey:UserID;association_autoupdate:false"`
	Venue *Venue     `gorm:"foreignkey:VenueID;association_autoupdate:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Question that will be served to participant
// UserAnswer is the user current answer to related question
type Question struct {
	ID         uint   `gorm:"primary_key"`
	Content    string `gorm:"type:text"`
	EventID    uint
	UserAnswer string `gorm:"-"`
	Choices    string // pipe (|) separated list of choices

	Event *Event `gorm:"foreignkey:EventID;association_autoupdate:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// UserQuestion is the many-to-many relation between user and question
// because user can't see al questions, but only theirs.
type UserQuestion struct {
	ID              uint `gorm:"primary_key"`
	ParticipationID uint
	QuestionID      uint
	Ordering        uint
	Answer          string `gorm:"type:text"`

	Participation *Participation `gorm:"foreignkey:ParticipationID;association_autoupdate:false"`
	Question      *Question      `gorm:"foreignkey:QuestionID;association_autoupdate:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func init() {
	helios.App.RegisterModel(Event{})
	helios.App.RegisterModel(Venue{})
	helios.App.RegisterModel(Participation{})
	helios.App.RegisterModel(Question{})
	helios.App.RegisterModel(UserQuestion{})
}
