package errors

var (
	ErrApiNotFound     = NotFound("API_NOT_FOUND", "接口不存在")
	ErrArticleNotFound = NotFound("ARTICLE_NOT_FOUND", "页面不存在")
	ErrSystem          = InternalServer("SYSTEM_ERROR", "系统错误")
	ErrUnauthorized    = Unauthorized("UNAUTHORIZED", "登录过期，请重新登录")
)
