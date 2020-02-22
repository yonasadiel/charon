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

	userLoggedIn, errLoggedIn := Login(LoginRequest{Email: "abc", Password: "def"})

	assert.Nil(t, errLoggedIn, "Expected success login, but get error: %s", errLoggedIn)
	assert.Equal(t, user.ID, userLoggedIn.ID, "Wrong user returned")
}

func TestLoginWrongUsername(t *testing.T) {
	helios.App.BeforeTest()

	password := hashPassword("def")
	user := User{Email: "abc", Password: password}
	helios.DB.Create(&User{Email: "def", Password: password})
	helios.DB.Create(&user)
	helios.DB.Create(&User{Email: "ghi", Password: password})
	helios.DB.Create(&User{Email: "jkl", Password: password})

	userLoggedIn, errLoggedIn := Login(LoginRequest{Email: "mno", Password: "def"})

	assert.Equal(t, errWrongUsernamePassword, *errLoggedIn, "Expected wrong username / password, but success logging in")
	assert.Nil(t, userLoggedIn, "Not nil user")
}

func TestLoginWrongPassword(t *testing.T) {
	helios.App.BeforeTest()

	password := hashPassword("def")
	user := User{Email: "abc", Password: password}
	helios.DB.Create(&User{Email: "def", Password: password})
	helios.DB.Create(&user)
	helios.DB.Create(&User{Email: "ghi", Password: password})
	helios.DB.Create(&User{Email: "jkl", Password: password})

	userLoggedIn, errLoggedIn := Login(LoginRequest{Email: "abc", Password: "abc"})

	assert.Equal(t, errWrongUsernamePassword, *errLoggedIn, "Expected wrong username / password, but success logging in")
	assert.Nil(t, userLoggedIn, "Not nil user")
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
