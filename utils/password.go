package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func EncryptPassword(password, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	md := hash.Sum([]byte(salt))
	return hex.EncodeToString(md)
}
