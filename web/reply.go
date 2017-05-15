package web

import (
	"encoding/json"
	"net/http"

	"github.com/swanwish/go-common/logs"
)

func (ctx HandlerContext) ReplyText(text string) {
	ctx.W.Header().Set("Content-Type", ContentTypePlain)
	if enableCors {
		ctx.EnableCors(enableCors)
	}
	ctx.W.Write([]byte(text))
}

func (ctx HandlerContext) ReplyOK() {
	ctx.ReplyText("OK")
}

func (ctx HandlerContext) ReplyError(errorMessage string, errorCode int) {
	logs.Errorf("Error message: %s and error code is: %d", errorMessage, errorCode)
	data := make(map[string]interface{})
	data[JsonKeyErrorMessage] = errorMessage
	data[JsonKeyErrorNo] = errorCode
	ctx.ReplyJson(data)
}

func (ctx HandlerContext) ReplyHttpError(error string, code int) {
	http.Error(ctx.W, error, code)
}

func (ctx HandlerContext) ReplyInvalidParameterError() {
	ctx.ReplyError(ErrorMessageInvalidParameter, StatusCodeInvalidParameter)
}

func (ctx HandlerContext) ReplyNotExistError() {
	ctx.ReplyError(ErrorMessageNotExist, StatusCodeNotExists)
}

func (ctx HandlerContext) ReplyInvalidParameterErrorWithMesssage(message string) {
	if message == "" {
		message = ErrorMessageInvalidParameter
	}
	ctx.ReplyError(message, StatusCodeInvalidParameter)
}

func (ctx HandlerContext) ReplyInternalError() {
	ctx.ReplyError(ErrorMessageInternalError, StatusCodeInternalError)
}

func (ctx HandlerContext) ReplyForbiddenError() {
	ctx.ReplyError(ErrorMessageForbidden, StatusCodeForbidden)
}

func (ctx HandlerContext) ReplyJson(data map[string]interface{}) {
	if data[JsonKeyErrorNo] == nil {
		data[JsonKeyErrorNo] = StatusCodeOK
	}
	if data[JsonKeyErrorMessage] == nil {
		data[JsonKeyErrorMessage] = JsonSuccessMessage
	}

	ctx.ReplyJsonObject(data)
}

func (ctx HandlerContext) ReplyJsonObject(object interface{}) {
	js, err := json.Marshal(object)
	if err != nil {
		logs.Errorf("Marshal json failed, the error is %v.", err)
		return
	}

	ctx.ReplyData(ContentTypeJson, js)
}

func (ctx HandlerContext) ReplyHtml(html string) {
	ctx.ReplyData(ContentTypeHtml, []byte(html))
}

func (ctx HandlerContext) ReplyData(contentType string, data []byte) {
	ctx.W.Header().Set("Content-Type", contentType)
	if enableCors {
		ctx.EnableCors(enableCors)
	}
	ctx.W.Write(data)
}

func (ctx HandlerContext) ReplyJsonData(data interface{}) {
	result := make(map[string]interface{}, 0)
	if data != nil {
		result[JsonKeyData] = data
	}
	ctx.ReplyJson(result)
}

func (ctx HandlerContext) ReplyJsonList(list interface{}) {
	result := make(map[string]interface{}, 0)
	result[JsonKeyList] = list
	ctx.ReplyJson(result)
}

func (ctx HandlerContext) ReplyInvalidTokenError() {
	ctx.ReplyError(ErrorMessageInvalidToken, StatusCodeTokenInvalid)
}
