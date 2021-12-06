package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloadFromUrl(t *testing.T) {
	testCases := []struct {
		rawUrl        string
		dest          string
		filePath      string
		checkFilePath string
		hasError      bool
	}{
		{rawUrl: "http://notexists", dest: "/tmp/downloader", filePath: "", checkFilePath: "", hasError: true},
		{rawUrl: "https://github.com/swanwish/tools/raw/master/linux/restart_service.sh", dest: "/tmp/downloader", filePath: "", checkFilePath: "/tmp/downloader/swanwish/tools/raw/master/linux/restart_service.sh", hasError: false},
		{rawUrl: "https://github.com/nouser/nopath/nofile.txt", dest: "/tmp/downloader", filePath: "", checkFilePath: "/tmp/downloader/nouser/nopath/nofile.txt", hasError: true},
	}

	for _, testCase := range testCases {
		if testCase.checkFilePath != "" {
			_ = DeleteFile(testCase.checkFilePath)
		}
		err := DownloadFromUrl(testCase.rawUrl, testCase.dest, testCase.filePath)
		if testCase.hasError {
			assert.NotNil(t, err)
			if testCase.checkFilePath != "" {
				assert.False(t, FileExists(testCase.checkFilePath))
			}
		} else {
			assert.Nil(t, err)
			assert.True(t, FileExists(testCase.checkFilePath))
		}
	}
}
