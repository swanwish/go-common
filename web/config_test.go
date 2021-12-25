package web

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadViperSettings(t *testing.T) {
	v := viper.GetViper()
	v.SetDefault(fmt.Sprintf("web.%s", KeyEnableUserIdentity), 1)
	v.SetDefault(fmt.Sprintf("web.%s", KeyValidUserIdentityListJson), `[{"key":"k", "value":"v"}]`)
	LoadViperSettings(v)
	assert.Equal(t, true, EnableUserIdentityCheck)
	v.SetDefault(fmt.Sprintf("web.%s", KeyEnableUserIdentity), 0)
	LoadViperSettings(v)
	assert.Equal(t, false, EnableUserIdentityCheck)
}
