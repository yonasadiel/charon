package auth

import (
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/yonasadiel/helios"
)

var src = rand.NewSource(time.Now().UnixNano())

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

// generateUserToken generates token of length userTokenLength
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func generateUserToken() string {
	n := userTokenLength
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), userTokenIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), userTokenIdxMax
		}
		if idx := int(cache & userTokenIdxMask); idx < len(userTokenBytes) {
			sb.WriteByte(userTokenBytes[idx])
			i--
		}
		cache >>= userTokenIdxBits
		remain--
	}

	return sb.String()
}

// Login will try to authenticate user and store the session
// if it fails, it will give helios.Error, If it success, it will
// return a new session
func Login(username string, password string, ip string) (*Session, helios.Error) {
	var user User
	var session Session
	var token string

	helios.DB.Where(&User{Username: username}).First(&user)

	if user.ID == 0 {
		return nil, errWrongUsernamePassword
	}

	if !checkPasswordHash(password, user.Password) {
		return nil, errWrongUsernamePassword
	}

	token = generateUserToken()
	session = Session{
		UserID:    user.ID,
		Token:     token,
		User:      &user,
		IPAddress: ip,
	}
	helios.DB.Create(&session)

	return &session, nil
}

// Logout invalidates the session token
func Logout(user User) {
	helios.DB.Where("user_id = ?", user.ID).Delete(&Session{})
}
