package utils

import (
	"crypto/sha1"
	"encoding/hex"
)

func GenerateSHA1(name string) string {
	hasher := sha1.New()
	hasher.Write([]byte(name))
	sha := hasher.Sum(nil)
	shaHex := hex.EncodeToString(sha)
	return shaHex
}
