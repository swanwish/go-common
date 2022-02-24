package utils

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRequest(t *testing.T) {
	requestUrl := "https://4.dbt.io/api/download/list"
	querys := url.Values{}
	querys.Set("v", "4")
	querys.Set("key", "1462b719-42d8-0874-7c50-905063472458")
	querys.Set("page", "2")
	_, content, err := GetRequest(requestUrl, nil, querys, nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, content)
}
