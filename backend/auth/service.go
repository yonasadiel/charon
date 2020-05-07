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

	helios.DB.Where("username = ?", username).First(&user)

	if user.ID == 0 {
		return nil, errWrongUsernamePassword
	}

	if !checkPasswordHash(password, user.Password) {
		return nil, errWrongUsernamePassword
	}

	if user.SessionLocked {
		return nil, errSessionLocked
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
	user.SessionLocked = false
	helios.DB.Save(&user)
	helios.DB.Where("user_id = ?", user.ID).Delete(&Session{})
}

// GetAllUser returns all users with lower role.
func GetAllUser(user User) []User {
	var users []User
	helios.DB.Where("role < ?", user.Role).Find(&users)
	return users
}

// UpsertUser creates or updates a user. It creates if
// ID = 0, or updates otherwise. The role should be
// less then invoker's role
// If it is create, then user.ID will be changed.
func UpsertUser(user User, newUser *User) helios.Error {
	if newUser.Role >= user.Role {
		return errUserRoleTooHigh
	}

	newUser.Password = hashPassword(newUser.Password)
	if newUser.ID == 0 {
		helios.DB.Create(newUser)
	} else {
		helios.DB.Save(newUser)
	}
	return nil
}

// LockUserSession locks user so they can't login from new device
func LockUserSession(user User, username string) helios.Error {
	var targetUser User
	helios.DB.Where("username = ?", username).Where("role < ?", user.Role).First(&targetUser)
	if targetUser.ID == 0 {
		return errUserNotFound
	}
	targetUser.SessionLocked = true
	helios.DB.Save(&targetUser)
	return nil
}

// UnlockUserSession reverse the lock effect, so the user can login from new device
func UnlockUserSession(user User, username string) helios.Error {
	var targetUser User
	helios.DB.Where("username = ?", username).Where("role < ?", user.Role).First(&targetUser)
	if targetUser.ID == 0 {
		return errUserNotFound
	}
	targetUser.SessionLocked = false
	helios.DB.Save(&targetUser)
	return nil
}
