package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yonasadiel/helios"
)

func TestLogin(t *testing.T) {
	helios.App.BeforeTest()

	type loginTestCase struct {
		user          User
		username      string
		password      string
		ipAddr        string
		expectedError helios.Error
	}
	testCases := []loginTestCase{
		loginTestCase{
			user:     UserFactorySaved(User{Username: "user1", Password: "def"}),
			username: "user1",
			password: "def",
		},
		loginTestCase{
			user:          UserFactorySaved(User{Username: "user2", Password: "def"}),
			username:      "def",
			password:      "def",
			expectedError: errWrongUsernamePassword,
		},
		loginTestCase{
			user:          UserFactorySaved(User{Username: "user3", Password: "def"}),
			username:      "user3",
			password:      "abc",
			expectedError: errWrongUsernamePassword,
		},
	}
	for i, testCase := range testCases {
		t.Logf("Test Login testcase: %d", i)
		var userSession *Session
		var userSessionSaved Session
		var err helios.Error
		userSession, err = Login(testCase.username, testCase.password, "1.2.3.4")
		if testCase.expectedError == nil {
			helios.DB.Where("token = ?", userSession.Token).First(&userSessionSaved)
			assert.Nil(t, err)
			assert.NotNil(t, userSession)
			assert.Equal(t, userTokenLength, len(userSession.Token))
			assert.Equal(t, testCase.user.ID, userSession.UserID)
			assert.NotEqual(t, 0, userSessionSaved.ID, "Session not saved on database")
			assert.Equal(t, testCase.user.ID, userSessionSaved.UserID)
			assert.Equal(t, "1.2.3.4", userSessionSaved.IPAddress)
		} else {
			assert.Equal(t, testCase.expectedError, err)
			assert.Nil(t, userSession)
		}
	}
}

func TestHashPassword(t *testing.T) {
	passwordHashed := hashPassword("charon")
	assert.NotEmpty(t, passwordHashed, "Hashed Password is empty")
}

func TestCheckPasswordHash(t *testing.T) {
	check := checkPasswordHash("charon", "$2a$14$RgL6IqGdMZTkTibAWfuoSeOoc6OpuHezUh3PK4hBLza45pwHx4f7K")
	assert.True(t, check, "Password mismatch")
}
