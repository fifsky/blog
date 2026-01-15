package errors

var (
	ErrApiNotFound     = NotFound("API_NOT_FOUND", "接口不存在")
	ErrArticleNotFound = BadRequest("ARTICLE_NOT_FOUND", "文章不存在")
	ErrSystem          = InternalServer("SYSTEM_ERROR", "系统错误")
	ErrUnauthorized    = Unauthorized("UNAUTHORIZED", "登录过期，请重新登录")
)
