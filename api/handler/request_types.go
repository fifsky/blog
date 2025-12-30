package handler

// 登录请求参数
type LoginRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

// 通用分页请求
type PageRequest struct {
	Page int `json:"page"`
}

// 通用ID请求
type IDRequest struct {
	Id int `json:"id"`
}

// 文章列表请求参数
type ArticleListRequest struct {
	Year    string `json:"year"`
	Month   string `json:"month"`
	Domain  string `json:"domain"`
	Keyword string `json:"keyword"`
	Page    int    `json:"page"`
	Type    int    `json:"type"`
}

// 上下篇请求
type PrevNextRequest struct {
	Id int `json:"id"`
}

// 文章详情请求
type ArticleDetailRequest struct {
	Id  int    `json:"id"`
	Url string `json:"url"`
}

// 分类提交请求
type CateRequest struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Domain string `json:"domain"`
}

// 文章提交请求
type ArticlePostRequest struct {
	Id      int    `json:"id"`
	CateId  int    `json:"cate_id"`
	Type    int    `json:"type"`
	Title   string `json:"title"`
	Url     string `json:"url"`
	Content string `json:"content"`
	Status  int    `json:"status"`
}

// 评论提交请求
type CommentRequest struct {
	PostId  int    `json:"post_id"`
	Pid     int    `json:"pid"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

// 友链提交请求
type LinkRequest struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
	Desc string `json:"desc"`
}

// 心情提交请求
type MoodRequest struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
}

// 提醒提交请求
type RemindRequest struct {
	Id       int    `json:"id"`
	Type     int    `json:"type"`
	Content  string `json:"content"`
	Month    int    `json:"month"`
	Week     int    `json:"week"`
	Day      int    `json:"day"`
	Hour     int    `json:"hour"`
	Minute   int    `json:"minute"`
	NextTime string `json:"next_time"`
}

// 用户提交请求
type UserRequest struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	NickName string `json:"nick_name"`
	Email    string `json:"email"`
	Status   int    `json:"status"`
	Type     int    `json:"type"`
}
