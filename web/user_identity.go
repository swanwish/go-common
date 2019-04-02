package web

import (
	"net/http"

	"github.com/swanwish/go-common/logs"
)

var (
	EnableUserIdentityCheck bool
	UserIdentityList        []UserIdentity
)

type UserIdentity struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Role  int64  `json:"role"`
}

func MakeHeaderCheckHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		LogRequest(r)
		logs.Debugf("The header is: %v", r.Header)
		if EnableUserIdentityCheck {
			ctx := CreateHandlerContext(rw, r)
			valid := false
			for _, userIdentity := range UserIdentityList {
				if ctx.HeaderValue(userIdentity.Key) == userIdentity.Value {
					valid = true
					break
				}
			}
			if !valid {
				ctx.ReplyForbiddenError()
				return
			}
		}
		fn(rw, r)
	}
}

type HeaderCheckHandler struct {
	next http.Handler
}

func NewHeaderCheckHandler(handler http.Handler) *HeaderCheckHandler {
	return &HeaderCheckHandler{next: handler}
}

func (handler *HeaderCheckHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	MakeHeaderCheckHandler(handler.next.ServeHTTP)(rw, r)
}
