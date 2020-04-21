package auth

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yonasadiel/helios"
)

func TestSerializeUser(t *testing.T) {
	type serializeUserTestCase struct {
		user         User
		expectedJSON string
	}
	testCases := []serializeUserTestCase{{
		user:         UserFactory(User{Name: "User 1", Username: "user1", Password: "abcd"}),
		expectedJSON: `{"name":"User 1","username":"user1","role":"participant"}`,
	}, {
		user:         UserFactory(User{Name: "User 2", Username: "user2", Password: "abcd", Role: UserRoleLocal}),
		expectedJSON: `{"name":"User 2","username":"user2","role":"local"}`,
	}, {
		user:         UserFactory(User{Name: "User 3", Username: "user3", Password: "abcd", Role: UserRoleOrganizer}),
		expectedJSON: `{"name":"User 3","username":"user3","role":"organizer"}`,
	}, {
		user:         UserFactory(User{Name: "User 4", Username: "user4", Password: "abcd", Role: UserRoleAdmin}),
		expectedJSON: `{"name":"User 4","username":"user4","role":"admin"}`,
	}}
	for i, testCase := range testCases {
		t.Logf("Test SerializeUser testcase: %d", i)
		var serialized []byte
		var errMarshalling error

		serialized, errMarshalling = json.Marshal(SerializeUser(testCase.user))
		assert.Nil(t, errMarshalling)
		assert.Equal(t, testCase.expectedJSON, string(serialized))
	}
}

func TestDeserializeUser(t *testing.T) {
	type deserializeUserTestCase struct {
		userDataJSON  string
		expectedUser  User
		expectedError string
	}
	testCases := []deserializeUserTestCase{{
		userDataJSON: `{"name":"User 1","username":"user1","role":"participant"}`,
		expectedUser: User{
			Name:     "User 1",
			Username: "user1",
			Role:     UserRoleParticipant,
		},
	}, {
		userDataJSON: `{"id":3,"name":"User 2","username":"user2","role":"admin"}`,
		expectedUser: User{
			ID:       0,
			Name:     "User 2",
			Username: "user2",
			Role:     UserRoleAdmin,
		},
	}, {
		userDataJSON: `{"name":"User 3","username":"user3","role":"local"}`,
		expectedUser: User{
			Name:     "User 3",
			Username: "user3",
			Role:     UserRoleLocal,
		},
	}, {
		userDataJSON: `{"name":"User 4","username":"user4","role":"organizer"}`,
		expectedUser: User{
			Name:     "User 4",
			Username: "user4",
			Role:     UserRoleOrganizer,
		},
	}, {
		userDataJSON:  `{"name":"User 5","username":"user5","role":"random"}`,
		expectedError: `{"code":"form_error","message":{"_error":[],"role":["Role should be either admin, organizer, local, or participant"]}}`,
	}, {
		userDataJSON:  `{}`,
		expectedError: `{"code":"form_error","message":{"_error":[],"name":["Name can't be empty"],"role":["Role can't be empty"],"username":["Username can't be empty"]}}`,
	}}
	for i, testCase := range testCases {
		t.Logf("Test DeserializeUser testcase: %d", i)
		var user User
		var userData UserData
		var errUnmarshalling error
		var errDeserialization helios.Error
		errUnmarshalling = json.Unmarshal([]byte(testCase.userDataJSON), &userData)
		errDeserialization = DeserializeUser(userData, &user)
		assert.Nil(t, errUnmarshalling)
		if testCase.expectedError == "" {
			assert.Nil(t, errDeserialization)
			assert.Equal(t, testCase.expectedUser.ID, user.ID, "Empty id on json will give 0")
			assert.Equal(t, testCase.expectedUser.Name, user.Name)
			assert.Equal(t, testCase.expectedUser.Username, user.Username)
			assert.Equal(t, testCase.expectedUser.Role, user.Role)
		} else {
			var errDeserializationJSON []byte
			var errMarshalling error
			errDeserializationJSON, errMarshalling = json.Marshal(errDeserialization.GetMessage())
			assert.Nil(t, errMarshalling)
			assert.NotNil(t, errDeserialization)
			assert.Equal(t, testCase.expectedError, string(errDeserializationJSON))
		}
	}
}

func TestSerializeUserWithPassword(t *testing.T) {
	type serializeUserWithPasswordTestCase struct {
		user         User
		expectedJSON string
	}
	testCases := []serializeUserWithPasswordTestCase{{
		user:         User{Name: "User 1", Username: "user1", Password: "abcd", Role: UserRoleOrganizer},
		expectedJSON: `{"name":"User 1","username":"user1","role":"organizer","password":"abcd"}`,
	}}
	for i, testCase := range testCases {
		t.Logf("Test SerializeUserWithPassword testcase: %d", i)
		var serialized []byte
		var errMarshalling error

		serialized, errMarshalling = json.Marshal(SerializeUserWithPassword(testCase.user))
		assert.Nil(t, errMarshalling)
		assert.Equal(t, testCase.expectedJSON, string(serialized))
	}
}

func TestDeserializeUserWithPassword(t *testing.T) {
	type deserializeUserWithPasswordTestCase struct {
		userDataJSON  string
		expectedUser  User
		expectedError string
	}
	testCases := []deserializeUserWithPasswordTestCase{{
		userDataJSON: `{"name":"User 1","username":"user1","role":"participant","password":"abcd"}`,
		expectedUser: User{
			Name:     "User 1",
			Username: "user1",
			Role:     UserRoleParticipant,
			Password: "abcd",
		},
	}, {
		userDataJSON:  `{"name":"User 5","username":"user5","role":"random","password":"abc"}`,
		expectedError: `{"code":"form_error","message":{"_error":[],"role":["Role should be either admin, organizer, local, or participant"]}}`,
	}}
	for i, testCase := range testCases {
		t.Logf("Test DeserializeUserWithPassword testcase: %d", i)
		var user User
		var userData UserWithPasswordData
		var errUnmarshalling error
		var errDeserialization helios.Error
		errUnmarshalling = json.Unmarshal([]byte(testCase.userDataJSON), &userData)
		errDeserialization = DeserializeUserWithPassword(userData, &user)
		assert.Nil(t, errUnmarshalling)
		if testCase.expectedError == "" {
			assert.Nil(t, errDeserialization)
			assert.Equal(t, testCase.expectedUser.ID, user.ID, "Empty id on json will give 0")
			assert.Equal(t, testCase.expectedUser.Name, user.Name)
			assert.Equal(t, testCase.expectedUser.Username, user.Username)
			assert.Equal(t, testCase.expectedUser.Role, user.Role)
			assert.Equal(t, testCase.expectedUser.Password, user.Password)
		} else {
			var errDeserializationJSON []byte
			var errMarshalling error
			errDeserializationJSON, errMarshalling = json.Marshal(errDeserialization.GetMessage())
			assert.Nil(t, errMarshalling)
			assert.NotNil(t, errDeserialization)
			assert.Equal(t, testCase.expectedError, string(errDeserializationJSON))
		}
	}
}
