package exam

import (
	"math/big"
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

func generateNRandomBigInt(n int) []big.Int {
	var primelength uint = 256
	var twoPower *big.Int = new(big.Int).Lsh(big.NewInt(1), primelength)
	var randSrc *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	var randoms []big.Int
	for i := 0; i < n; i++ {
		var rand *big.Int = new(big.Int).Rand(randSrc, twoPower)
		randoms = append(randoms, *rand)
	}
	return randoms
}
