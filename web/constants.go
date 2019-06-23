package web

const (
	JsonSuccessMessage  = "success"
	JsonKeyErrorNo      = "errno"
	JsonKeyErrorMessage = "msg"
	JsonKeyData         = "data"
	JsonKeyList         = "list"
)

const (
	StatusCodeOK               = 200
	StatusCodeInvalidParameter = 400
	StatusCodeTokenInvalid     = 401
	StatusCodeForbidden        = 403
	StatusCodeNotExists        = 404
	StatusCodeInternalError    = 500
)

const (
	ErrorMessageInternalError    = "Internal Error"
	ErrorMessageInvalidParameter = "Invalid Parameter"
	ErrorMessageForbidden        = "Forbidden"
	ErrorMessageNotExist         = "Not exist"
	ErrorMessageInvalidToken     = "Invalid token"
)

const (
	ContentTypePlain = "text/plain; charset=UTF-8"
	ContentTypeJson  = "application/json; charset=UTF-8"
	ContentTypeHtml  = "text/html; charset=utf-8"
)
