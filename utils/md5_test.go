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
		{filePath: "./md5.go", hasError: false, expected: "916ce96e225213ceaf8a4c2433917091"},
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
