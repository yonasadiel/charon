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
	testCases := []loginTestCase{{
		user:     UserFactorySaved(User{Username: "user1", Password: "def"}),
		username: "user1",
		password: "def",
	}, {
		user:          UserFactorySaved(User{Username: "user2", Password: "def"}),
		username:      "def",
		password:      "def",
		expectedError: errWrongUsernamePassword,
	}, {
		user:          UserFactorySaved(User{Username: "user3", Password: "def"}),
		username:      "user3",
		password:      "abc",
		expectedError: errWrongUsernamePassword,
	}}
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

func TestGetAllUser(t *testing.T) {
	helios.App.BeforeTest()
	var userAdmin User = UserFactorySaved(User{Role: UserRoleAdmin})
	var userOrganizer User = UserFactorySaved(User{Role: UserRoleOrganizer})
	var userLocal User = UserFactorySaved(User{Role: UserRoleLocal})
	var userParticipant User = UserFactorySaved(User{Role: UserRoleParticipant})
	type getAllUserTestCase struct {
		user           User
		expectedLength int
	}
	testCases := []getAllUserTestCase{{
		user:           userAdmin,
		expectedLength: 3,
	}, {
		user:           userOrganizer,
		expectedLength: 2,
	}, {
		user:           userLocal,
		expectedLength: 1,
	}, {
		user:           userParticipant,
		expectedLength: 0,
	}}
	for i, testCase := range testCases {
		t.Logf("Test GetAllUser testcase: %d", i)
		var users []User
		users = GetAllUser(testCase.user)
		assert.Equal(t, testCase.expectedLength, len(users))
	}
}

func TestUpsertUser(t *testing.T) {
	helios.App.BeforeTest()

	var userParticipant User = UserFactorySaved(User{Role: UserRoleParticipant})
	var userLocal User = UserFactorySaved(User{Role: UserRoleLocal})

	type upsertUserTestCase struct {
		user              User
		newUser           User
		password          string
		expectedError     helios.Error
		expectedUserCount int
	}
	testCases := []upsertUserTestCase{{
		user:              userParticipant,
		newUser:           UserFactory(User{Role: UserRoleParticipant}),
		password:          "pass1",
		expectedError:     errUserRoleTooHigh,
		expectedUserCount: 2,
	}, {
		user:              userLocal,
		newUser:           UserFactory(User{Role: UserRoleParticipant}),
		password:          "pass2",
		expectedUserCount: 3,
	}, {
		user:              userLocal,
		newUser:           UserFactory(User{ID: userLocal.ID, Name: "abc"}),
		password:          "pass3",
		expectedUserCount: 3,
	}}
	for i, testCase := range testCases {
		var newUserCount int
		var newUserSaved User
		t.Logf("Test UpsertUser testcase: %d", i)
		testCase.newUser.Password = testCase.password
		err := UpsertUser(testCase.user, &testCase.newUser)
		helios.DB.Model(User{}).Count(&newUserCount)
		helios.DB.Where("id = ?", testCase.newUser.ID).First(&newUserSaved)
		assert.Equal(t, testCase.expectedUserCount, newUserCount)
		if testCase.expectedError == nil {
			assert.Nil(t, err)
			assert.Equal(t, testCase.newUser.Name, newUserSaved.Name, "If the newUser has already existed, it should be updated")
			assert.NotEqual(t, testCase.newUser.Password, testCase.password)
			assert.True(t, checkPasswordHash(testCase.password, testCase.newUser.Password), "password should be hashed")
		} else {
			assert.Equal(t, testCase.expectedError, err)
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
