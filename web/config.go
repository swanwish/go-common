package web

import (
	"sync"

	"github.com/gorilla/sessions"
	"github.com/swanwish/go-common/logs"
	"github.com/swanwish/go-common/utils"
)

var (
	sessionName string
	store       *sessions.CookieStore
	storeLock   = &sync.Mutex{}
	enableCors  = false
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
