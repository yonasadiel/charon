package auth

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/yonasadiel/charon/app"
)

func hashPassword(password string) string {
	// we ignore error because the failure
	// usually because of cost error (which we set 14)
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Login will try to authenticate user and store the session
// if it fails, it will give APIError
func Login(r LoginRequest) (*User, *app.APIError) {
	var user User
	app.Charon.DB.Where(&User{Email: r.Email}).First(&user)

	if user.ID == 0 {
		return nil, &errWrongUsernamePassword
	}

	if !checkPasswordHash(r.Password, user.Password) {
		return nil, &errWrongUsernamePassword
	}

	return &user, nil
}

var errWrongUsernamePassword = app.APIError{
	StatusCode: http.StatusBadRequest,
	Code:       "login_wrong_email_or_password",
	Message:    "Wrong email / password",
}
