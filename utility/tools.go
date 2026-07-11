package utility

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256Hash(str string) string {
	hash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hash[:])
}
