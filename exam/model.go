package exam

import (
	"time"

	"github.com/yonasadiel/charon/auth"
	"github.com/yonasadiel/helios"
)

// Question that will be served to participant
// UserAnswer is the user current answer to related question
type Question struct {
	ID         uint   `gorm:"primary_key"`
	Content    string `gorm:"type:text"`
	UserAnswer string `gorm:"-"`

	Choices []QuestionChoice `gorm:"foreignkey:QuestionID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// QuestionChoice is choice of question with multiple choices type
type QuestionChoice struct {
	ID         uint `gorm:"primray_key"`
	Text       string
	QuestionID uint

	Question *Question `gorm:"foreignkey:QuestionID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// UserQuestion is the many-to-many relation between user and question
// because user can't see al questions, but only theirs.
type UserQuestion struct {
	ID         uint `gorm:"primary_key"`
	QuestionID uint
	UserID     uint
	Ordering   uint

	Question *Question  `gorm:"foreignkey:QuestionID"`
	User     *auth.User `gorm:"foreignkey:UserID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Submission of user to the question. This model should be insert only,
// never deleted, for log purpose.
type Submission struct {
	ID         uint   `gorm:"primary_key"`
	Answer     string `gorm:"type:text"`
	QuestionID uint
	UserID     uint

	Question *Question  `gorm:"foreignkey:QuestionID"`
	User     *auth.User `gorm:"foreignkey:UserID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func init() {
	helios.App.RegisterModel(QuestionChoice{})
	helios.App.RegisterModel(Question{})
	helios.App.RegisterModel(UserQuestion{})
	helios.App.RegisterModel(Submission{})
}
