package utils

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/swanwish/go-common/logs"
)

var (
	ErrInvalidStatus = errors.New("Invalid status")
	ErrAlreadyExists = errors.New("Already exists")
)

func DownloadFromUrl(rawUrl, dest string) error {
	u, err := url.Parse(rawUrl)
	if err != nil {
		logs.Errorf("The url %s is invalid, the error is %v", rawUrl, err)
		return err
	}
	fileName := ""
	logs.Debugf("The path is %s", u.Path)
	path := u.Path
	lastSplashIndex := strings.LastIndex(path, "/")
	if lastSplashIndex != -1 {
		parentPath := path[:lastSplashIndex]
		logs.Debugf("The parent path is %s", parentPath)
		destDir := filepath.Join(dest, parentPath) // fmt.Sprintf("%s%s", dest, parentPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			logs.Errorf("Failed to create dir %s, the error is %v", destDir, err)
			return err
		}
		//fileName = path[lastSplashIndex+1:]
		//logs.Debugf("The file name is %s", fileName)
	}
	destFilePath := filepath.Join(dest, path) // fmt.Sprintf("%s%s", dest, path)
	if FileExists(destFilePath) {
		logs.Debugf("The file %s already exists", destFilePath)
		return ErrAlreadyExists
	}

	logs.Debugf("Downloading file from %s, save to %s", rawUrl, destFilePath)
	response, err := http.Get(rawUrl)
	if err != nil {
		logs.Errorf("Error while downloading from %s, the error is %v", rawUrl, err)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		logs.Errorf("Failed to download file from %s, the status is %s", rawUrl, response.Status)
		return ErrInvalidStatus
	}

	output, err := os.Create(destFilePath)
	if err != nil {
		logs.Errorf("Error while creating %s, the error is %v", fileName, err)
		return err
	}
	defer output.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		logs.Errorf("Error while downloading from %s, the error is %v", rawUrl, err)
		return err
	}

	logs.Debugf("%d bytes downloaded.", n)
	return nil
}
