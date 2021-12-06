package utils

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/swanwish/go-common/logs"
)

var (
	ErrInvalidStatus = errors.New("Invalid status")
	ErrAlreadyExists = errors.New("Already exists")
	ErrNotFound      = errors.New("not found")
)

func DownloadFromUrl(rawUrl, dest, filePath string) error {
	u, err := url.Parse(rawUrl)
	if err != nil {
		logs.Errorf("The url %s is invalid, the error is %v", rawUrl, err)
		return err
	}
	//fileName := ""
	if filePath == "" {
		logs.Debugf("The path is %s", u.Path)
		filePath = u.Path
	}

	destFilePath := filepath.Join(dest, filePath) // fmt.Sprintf("%s%s", dest, path)
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
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return ErrNotFound
		}
		logs.Errorf("Failed to download file from %s, the status code is %d, the status is %s", rawUrl, response.StatusCode, response.Status)
		return ErrInvalidStatus
	}

	parentDir := path.Dir(destFilePath)
	if !FileExists(parentDir) {
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			logs.Errorf("Failed to create dir %s, the error is %#v", parentDir, err)
			return err
		}
	}

	output, err := os.Create(destFilePath)
	if err != nil {
		logs.Errorf("Error while creating %s, the error is %v", filePath, err)
		return err
	}
	defer func() {
		_ = output.Close()
	}()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		logs.Errorf("Error while downloading from %s, the error is %v", rawUrl, err)
		return err
	}

	logs.Debugf("%d bytes downloaded.", n)
	return nil
}
