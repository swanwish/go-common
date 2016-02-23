package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/swanwish/go-common/logs"
)

func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func GetCurrentDirectory() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logs.Errorf("Failed to get current directory")
		return "", err
	}
	return strings.Replace(dir, "\\", "/", -1), nil
}
