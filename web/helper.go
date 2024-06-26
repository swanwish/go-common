package web

import (
	"net"
	"net/http"

	"github.com/swanwish/go-common/logs"
)

func LogRequest(r *http.Request) {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	logs.Debugf("%s %v from ip: %s", r.Method, r.URL, ip)
}

func MakeLogEnableRawHandler(fn func(rw http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		LogRequest(r)
		fn(rw, r)
	}
}

func MakeLogEnabledHandler(fn func(HandlerContext)) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		LogRequest(r)
		ctx := HandlerContext{W: rw, R: r}
		fn(ctx)
	}
}

func MakeWrappedLogEnabledHandler(fn func(HandlerContext)) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := CreateHandlerContext(rw, r)
		LogRequest(r)
		fn(ctx)
	}
}

func MakeSessionCheckHandler(fn func(SessionHandlerContext)) http.HandlerFunc {
	return MakeSessionCheckHandlerWithCallback(fn, nil)
}

func MakeSessionCheckHandlerWithCallback(fn func(SessionHandlerContext), errorCallback func(ctx HandlerContext)) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		LogRequest(r)
		handlerContext := HandlerContext{W: rw, R: r}
		loginUser := handlerContext.GetSessionValue(SessionKeyLoginUser)
		if loginUser == nil {
			if errorCallback != nil {
				errorCallback(handlerContext)
				return
			}
			handlerContext.ReplyError("Forbidden", 401)
			return
		}
		fn(SessionHandlerContext{handlerContext, *loginUser.(*LoginUser)})
	}
}
