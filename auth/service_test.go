package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yonasadiel/helios"
)

func TestLoginSuccess(t *testing.T) {
	helios.App.BeforeTest()

	password := hashPassword("def")
	user := User{Email: "abc", Password: password}
	helios.DB.Create(&User{Email: "def", Password: password})
	helios.DB.Create(&user)
	helios.DB.Create(&User{Email: "ghi", Password: password})
	helios.DB.Create(&User{Email: "jkl", Password: password})

	userSession, errLoggedIn := Login("abc", "def", "1.2.3.4")

	assert.Nil(t, errLoggedIn, "Expected success login, but get error: %s", errLoggedIn)
	assert.NotNil(t, userSession, "Empty session returned")
	assert.Equal(t, userTokenLength, len(userSession.Token), "Different token length")
	assert.Equal(t, user.ID, userSession.UserID, "Different user ID returned")

	var userSessionDB Session
	helios.DB.Where("token = ?", userSession.Token).First(&userSessionDB)
	assert.NotEqual(t, 0, userSessionDB.ID, "User token not found")
	assert.Equal(t, user.ID, userSessionDB.UserID, "Different user logged in")
}

func TestLoginWrongUsername(t *testing.T) {
	helios.App.BeforeTest()

	password := hashPassword("def")
	user := User{Email: "abc", Password: password}
	helios.DB.Create(&User{Email: "def", Password: password})
	helios.DB.Create(&user)
	helios.DB.Create(&User{Email: "ghi", Password: password})
	helios.DB.Create(&User{Email: "jkl", Password: password})

	userLoggedIn, errLoggedIn := Login("mno", "def", "1.2.3.4")

	assert.Equal(t, errWrongUsernamePassword, *errLoggedIn, "Expected wrong username / password, but success logging in")
	assert.Nil(t, userLoggedIn, "Not nil user session")
}

func TestLoginWrongPassword(t *testing.T) {
	helios.App.BeforeTest()

	password := hashPassword("def")
	user := User{Email: "abc", Password: password}
	helios.DB.Create(&User{Email: "def", Password: password})
	helios.DB.Create(&user)
	helios.DB.Create(&User{Email: "ghi", Password: password})
	helios.DB.Create(&User{Email: "jkl", Password: password})

	userLoggedIn, errLoggedIn := Login("abc", "abc", "1.2.3.4")

	assert.Equal(t, errWrongUsernamePassword, *errLoggedIn, "Expected wrong username / password, but success logging in")
	assert.Nil(t, userLoggedIn, "Not nil user session")
}

func TestHashPassword(t *testing.T) {
	helios.App.BeforeTest()

	passwordHashed := hashPassword("charon")
	assert.NotEmpty(t, passwordHashed, "Hashed Password is empty")
}

func TestCheckPasswordHash(t *testing.T) {
	check := checkPasswordHash("charon", "$2a$14$RgL6IqGdMZTkTibAWfuoSeOoc6OpuHezUh3PK4hBLza45pwHx4f7K")
	assert.True(t, check, "Password mismatch")
}
