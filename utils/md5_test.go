package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileMd5Hash(t *testing.T) {
	testCases := []struct {
		filePath string
		hasError bool
		expected string
	}{
		{filePath: "", hasError: true, expected: ""},
		{filePath: "/usr/bin/tail", hasError: false, expected: "0476c3a044da707a3c5319be141ba217"},
	}

	for _, testCase := range testCases {
		md5, err := GetFileMd5Hash(testCase.filePath)
		if testCase.hasError {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, testCase.expected, md5)
		}
	}
}
