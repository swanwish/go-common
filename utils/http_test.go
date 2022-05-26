package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRequest(t *testing.T) {
	requestUrl := "https://jsonplaceholder.typicode.com/todos/1"
	_, content, err := GetRequest(requestUrl, nil, nil, nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, content)
}
