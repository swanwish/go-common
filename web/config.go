package web

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/gorilla/sessions"
	"github.com/swanwish/go-common/config"
	"github.com/swanwish/go-common/logs"
	"github.com/swanwish/go-common/utils"
)

const (
	KeyEnableUserIdentity        = "enable_user_identity"
	KeyValidUserIdentityListJson = "valid_user_identity_list_json"
	KeyAllowedRefers             = "allowed_refers"
	KeyEnableCors                = "enable_cors"
)

var (
	sessionName       string
	store             *sessions.CookieStore
	storeLock         = &sync.Mutex{}
	enableCors        = false
	AllowedRefers     []string
	ReferCheckEnabled bool
)

func SetSessionName(name string) {
	sessionName = name
	logs.Debugf("The session name is %s", sessionName)
}

func getSessionName() string {
	if sessionName == "" {
		sessionName = utils.GenerateRandomStringEx(utils.RandomTypeCapitalString|utils.RandomTypeLowercaseChar|utils.RandomTypeDigital, SessionKeyLength)
		logs.Debugf("The session name is %s", sessionName)
	}
	return sessionName
}

func SetEnableCors(enable bool) {
	enableCors = enable
}

func LoadSettings() {
	if enableUserIdentityCheck, err := config.GetInt(KeyEnableUserIdentity); err == nil {
		EnableUserIdentityCheck = enableUserIdentityCheck == 1
		userIdentityListJson, err := config.Get(KeyValidUserIdentityListJson)
		if err != nil || userIdentityListJson == "" {
			logs.Infof("Identity user list json is not configured")
			EnableUserIdentityCheck = false
		}
		var userIdentityList []UserIdentity
		err = json.Unmarshal([]byte(userIdentityListJson), &userIdentityList)
		if err != nil {
			logs.Errorf("Failed to unmarshal user identity list json %s, the error is %v", userIdentityListJson, err)
			EnableUserIdentityCheck = false
		} else if len(userIdentityList) == 0 {
			EnableUserIdentityCheck = false
		} else {
			UserIdentityList = userIdentityList
			for _, userIdentity := range UserIdentityList {
				logs.Debugf("user identity key: %s, value: %s", userIdentity.Key, userIdentity.Value)
			}
		}
		if EnableUserIdentityCheck {
			logs.Infof("Identity check enabled")
		} else {
			logs.Info("Identity check is disabled")
		}
	}
	AllowedRefers = []string{}
	ReferCheckEnabled = false
	if allowedRefers, err := config.Get(KeyAllowedRefers); err == nil && allowedRefers != "" {
		refers := strings.Split(allowedRefers, ",")
		for _, refer := range refers {
			refer := strings.TrimSpace(refer)
			if refer != "" {
				AllowedRefers = append(AllowedRefers, refer)
				ReferCheckEnabled = true
			}
		}
	}
	if ReferCheckEnabled {
		logs.Infof("The allowed refers are %+v", AllowedRefers)
	} else {
		logs.Infof("The allowed refer is not enabled")
	}

	if enableCors, err := config.GetInt(KeyEnableCors); err == nil {
		if enableCors == 1 {
			logs.Infof("The cors is enabled")
			SetEnableCors(true)
		}
	}
}
