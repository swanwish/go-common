package web

import (
	"encoding/gob"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/swanwish/go-common/logs"
	"github.com/swanwish/go-common/utils"
)

const (
	SessionKeyLoginUser = "session_login_user"
	defaultMaxMemory    = 32 << 20 // 32 MB
	SessionKeyLength    = 16
)

func init() {
	gob.Register(&LoginUser{})
}

type LoginUser struct {
	UserId      string
	LoginName   string
	DisplayName string
	Email       string
}

var (
	sessionName string
	store       *sessions.CookieStore
	storeLock   = &sync.Mutex{}
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

func InitCookieStore(keyPairs ...[]byte) {
	if keyPairs == nil || len(keyPairs) == 0 || len(keyPairs[0]) == 0 {
		store = sessions.NewCookieStore(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
	} else {
		store = sessions.NewCookieStore(keyPairs...)
	}
}

func getStore() *sessions.CookieStore {
	if store != nil {
		return store
	}
	storeLock.Lock()
	defer storeLock.Unlock()
	if store != nil {
		return store
	}
	if store == nil {
		InitCookieStore(nil)
	}
	return store
}

type HandlerContext struct {
	W          http.ResponseWriter
	R          *http.Request
	bodyValues map[string]string
}

type SessionHandlerContext struct {
	HandlerContext
	LoginUser
}

func CreateHandlerContext(rw http.ResponseWriter, r *http.Request) HandlerContext {
	return HandlerContext{rw, r, nil}
}

func (ctx HandlerContext) Var(key string) string {
	return mux.Vars(ctx.R)[key]
}

func (ctx HandlerContext) VarInt(key string) int64 {
	return stringToInt(ctx.Var(key))
}

func stringToInt(str string) int64 {
	if str != "" {
		intValue, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			logs.Errorf("Failed to parse string value %s to int, the error is %v", str, err)
		} else {
			return intValue
		}
	}
	return 0
}

func (ctx HandlerContext) FormValue(key string) string {
	return ctx.R.FormValue(key)
}

func (ctx HandlerContext) FormValues(key string) []string {
	if ctx.R.Form == nil {
		ctx.R.ParseMultipartForm(defaultMaxMemory)
	}
	return ctx.R.Form[key]
}

func (ctx HandlerContext) FormIntValue(key string) int64 {
	return stringToInt(ctx.FormValue(key))
}

func (ctx HandlerContext) HeaderValue(key string) string {
	return ctx.R.Header.Get(key)
}

func (ctx HandlerContext) HeaderIntValue(key string) int64 {
	return stringToInt(ctx.R.Header.Get(key))
}

func (ctx HandlerContext) FormFloatValue(key string) float64 {
	strValue := ctx.FormValue(key)
	if strValue != "" {
		floatValue, err := strconv.ParseFloat(strValue, 64)
		if err != nil {
			logs.Errorf("Failed to parse string value %s to float, the error is %v", strValue, err)
		} else {
			return floatValue
		}
	}
	return 0.0
}

func (ctx *HandlerContext) ParseBodyValues() error {
	requestContent, err := ctx.GetRequestContent()
	if err != nil {
		logs.Errorf("Failed to get request content, the error is %v", err)
		return err
	}
	unescapedContent, err := url.QueryUnescape(string(requestContent))
	if err != nil {
		logs.Errorf("Failed to unescape request content %s", requestContent)
		return err
	}
	ctx.bodyValues = make(map[string]string, 0)
	lines := strings.Split(unescapedContent, "\n")
	for _, line := range lines {
		logs.Debugf("line %s", line)
		parts := strings.Split(line, "&")
		for _, part := range parts {
			pair := strings.Split(part, "=")
			if len(pair) == 2 {
				if pair[0] != "" {
					ctx.bodyValues[pair[0]] = pair[1]
				}
			}
		}
	}
	return nil
}

func (ctx *HandlerContext) BodyValue(key string) string {
	if ctx.bodyValues == nil {
		if err := ctx.ParseBodyValues(); err != nil {
			logs.Errorf("Failed to parse body value")
			return ""
		}
	}
	return ctx.bodyValues[key]
}

func (ctx *HandlerContext) BodyIntValue(key string) int64 {
	value := ctx.BodyValue(key)
	if value != "" {
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			logs.Errorf("Failed to parse string value %s to int, the error is %v", value, err)
		} else {
			return intValue
		}
	}
	return 0
}

func (ctx HandlerContext) Redirect(url string, code int) {
	http.Redirect(ctx.W, ctx.R, url, code)
}

func (ctx HandlerContext) ServeFile(filePath string) {
	utils.ServeFile(filePath, ctx.W, ctx.R, false)
}

func (ctx HandlerContext) GetRequestContent() ([]byte, error) {
	contents, err := ioutil.ReadAll(ctx.R.Body)
	if err != nil {
		logs.Errorf("Failed to read content from response, the error is %v", err)
	}
	return contents, err
}

func (ctx HandlerContext) GetClientIp() string {
	ip := ctx.HeaderValue("X-Real-IP")
	if ip == "" {
		ip, _, _ = net.SplitHostPort(ctx.R.RemoteAddr)
	}
	return ip
}

func (ctx HandlerContext) GetSession() *sessions.Session {
	session, _ := getStore().Get(ctx.R, getSessionName())
	return session
}

func (ctx HandlerContext) GetSessionValue(key interface{}) interface{} {
	session := ctx.GetSession()
	return session.Values[key]
}

func (ctx HandlerContext) SetSessionValue(key, value interface{}) {
	session := ctx.GetSession()
	session.Values[key] = value
}

func (ctx HandlerContext) SetSessionOptions(options *sessions.Options) {
	session := ctx.GetSession()
	session.Options = options
}

func (ctx HandlerContext) DeleteSessionKey(key interface{}) {
	session := ctx.GetSession()
	delete(session.Values, key)
}

func (ctx HandlerContext) SaveSession() {
	session := ctx.GetSession()
	session.Save(ctx.R, ctx.W)
}

func (ctx HandlerContext) GetLoginUser() *LoginUser {
	loginUser := ctx.GetSessionValue(SessionKeyLoginUser)
	if loginUser != nil {
		return loginUser.(*LoginUser)
	}
	return nil
}
