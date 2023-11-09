package utils

import (
	"io/fs"
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

func SaveFile(filePath string, data []byte, perm ...fs.FileMode) error {
	fileDir := filepath.Dir(filePath)
	if !FileExists(fileDir) {
		err := os.MkdirAll(fileDir, 0755)
		if err != nil {
			logs.Errorf("Failed to create dir %s, the error is %v", fileDir, err)
			return err
		}
	}
	var filePerm fs.FileMode = 0644
	if len(perm) > 0 {
		filePerm = perm[0]
	}
	return os.WriteFile(filePath, data, filePerm)
}

func DeleteFile(filePath string) error {
	if FileExists(filePath) {
		logs.Debugf("Delete file %s", filePath)
		err := os.Remove(filePath)
		if err != nil {
			logs.Errorf("Failed to delete file %s, the error is %#v", filePath, err)
		}
		return err
	}
	logs.Errorf("The file path %s does not exists", filePath)
	return os.ErrNotExist
}
