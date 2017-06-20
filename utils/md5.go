package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMD5Hash(text string) string {
	return GetByteMD5Hash([]byte(text))
}

func GetByteMD5Hash(content []byte) string {
	hasher := md5.New()
	hasher.Write(content)
	return hex.EncodeToString(hasher.Sum(nil))
}
