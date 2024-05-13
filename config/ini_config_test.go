package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	StringValue string `ini:"string_value"`
	Int64Value  int64  `ini:"int64_value"`
	IntValue    int    `ini:"int_value"`
	BoolValue   bool   `ini:"bool_value"`
}

func TestUnmarshal(t *testing.T) {
	iniConfig := IniConfiguration{}
	iniConfig.LoadContent(`string_value=string
int64_value=64
int_value=32
bool_value=true`)
	testConfig := TestConfig{}
	err := iniConfig.Unmarshal(&testConfig)
	assert.Nil(t, err)
	assert.Equal(t, "string", testConfig.StringValue)
	assert.Equal(t, int64(64), testConfig.Int64Value)
	assert.Equal(t, 32, testConfig.IntValue)
	assert.Equal(t, true, testConfig.BoolValue)
}
