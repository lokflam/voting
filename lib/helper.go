package lib

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// Hexdigest256 hashes a string to sha256 representation
func Hexdigest256(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	hashBytes := hash.Sum(nil)
	return strings.ToLower(hex.EncodeToString(hashBytes))
}

// Hexdigest512 hashes a string to sha512 representation
func Hexdigest512(str string) string {
	hash := sha512.New()
	hash.Write([]byte(str))
	hashBytes := hash.Sum(nil)
	return strings.ToLower(hex.EncodeToString(hashBytes))
}

// GenerateUUID returns a new UUID
func GenerateUUID(name string) string {
	return uuid.NewV5(uuid.Must(uuid.NewV1()), name).String()
}
