package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/swanwish/go-common/logs"
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
		destDir := fmt.Sprintf("%s%s", dest, parentPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			logs.Errorf("Failed to create dir %s, the error is %v", destDir, err)
			return err
		}
		//fileName = path[lastSplashIndex+1:]
		//logs.Debugf("The file name is %s", fileName)
	}
	destFilePath := fmt.Sprintf("%s%s", dest, path)
	if FileExists(destFilePath) {
		logs.Debugf("The file %s already exists", destFilePath)
		return nil
	}

	output, err := os.Create(destFilePath)
	if err != nil {
		logs.Errorf("Error while creating %s, the error is %v", fileName, err)
		return err
	}
	defer output.Close()

	logs.Debugf("Downloading file from %s, save to %s", rawUrl, destFilePath)
	response, err := http.Get(rawUrl)
	if err != nil {
		logs.Errorf("Error while downloading from %s, the error is %v", rawUrl, err)
		return err
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		logs.Errorf("Error while downloading from %s, the error is %v", rawUrl, err)
		return err
	}

	logs.Debugf("%d bytes downloaded.", n)
	return nil
}
