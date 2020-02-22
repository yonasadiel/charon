package auth

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/yonasadiel/helios"
)

func hashPassword(password string) string {
	// we ignore error because the failure
	// usually because of cost error
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	return string(bytes)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Login will try to authenticate user and store the session
// if it fails, it will give APIError
func Login(r LoginRequest) (*User, *helios.APIError) {
	var user User
	helios.DB.Where(&User{Email: r.Email}).First(&user)

	if user.ID == 0 {
		return nil, &errWrongUsernamePassword
	}

	if !checkPasswordHash(r.Password, user.Password) {
		return nil, &errWrongUsernamePassword
	}

	return &user, nil
}
