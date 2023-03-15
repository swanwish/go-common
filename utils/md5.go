package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"

	"github.com/swanwish/go-common/logs"
)

func GetMD5Hash(text string) string {
	return GetByteMD5Hash([]byte(text))
}

func GetByteMD5Hash(content []byte) string {
	hasher := md5.New()
	hasher.Write(content)
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetFileMd5Hash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		logs.Errorf("Failed to open file %s, the error is %#v", filePath, err)
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		logs.Errorf("Failed to calculate file md5, the error is %#v", err)
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
