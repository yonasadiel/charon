package exam

import (
	"math/rand"
	"strings"
	"time"
)

const (
	tokenBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	tokenIdxBits = 6                   // 6 bits to represent a token index
	tokenIdxMask = 1<<tokenIdxBits - 1 // All 1-bits, as many as tokenIdxBits
	tokenIdxMax  = 63 / tokenIdxBits   // # of token indices fitting in 63 bits
)

var randomSource = rand.NewSource(time.Now().UnixNano())

// generateRandomToken generates token
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func generateRandomToken(tokenLength int) string {
	sb := strings.Builder{}
	sb.Grow(tokenLength)
	// A randomSource.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := tokenLength-1, randomSource.Int63(), tokenIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randomSource.Int63(), tokenIdxMax
		}
		if idx := int(cache & tokenIdxMask); idx < len(tokenBytes) {
			sb.WriteByte(tokenBytes[idx])
			i--
		}
		cache >>= tokenIdxBits
		remain--
	}

	return sb.String()
}

// isAnswerValidChoice returns true if the answer is exist in one of question choices.
// Returns true if the choices are empty. It means that the question is not multiple choices type of question.
func isAnswerValidChoice(answer string, choices string) bool {
	var exist bool = false
	var choicesArr []string = strings.Split(choices, "|")
	var isEmpty bool = true
	for _, choice := range choicesArr {
		if choice != "" {
			isEmpty = false
			if choice == answer {
				exist = true
			}
		}
	}
	if isEmpty {
		return true
	}
	return exist
}
