package handler

// 登录请求参数
type LoginRequest struct {
	UserName string `json:"user_name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// 通用分页请求
type PageRequest struct {
	Page int `json:"page" validate:"min=1"`
}

// 通用ID请求
type IDRequest struct {
	Id int `json:"id" validate:"required,gte=1"`
}

// 文章列表请求参数
type ArticleListRequest struct {
	Year    string `json:"year"`
	Month   string `json:"month"`
	Domain  string `json:"domain"`
	Keyword string `json:"keyword"`
	Page    int    `json:"page" validate:"min=1"`
	Type    int    `json:"type"`
}

// 上下篇请求
type PrevNextRequest struct {
	Id int `json:"id" validate:"required,gte=1"`
}

// 文章详情请求
type ArticleDetailRequest struct {
	Id  int    `json:"id" validate:"required_without=Url,gte=1"`
	Url string `json:"url" validate:"required_without=Id"`
}

// 分类提交请求
type CateRequest struct {
	Id     int    `json:"id" validate:"omitempty,gte=1"`
	Name   string `json:"name" validate:"required"`
	Desc   string `json:"desc"`
	Domain string `json:"domain" validate:"required"`
}

// 文章提交请求
type ArticlePostRequest struct {
	Id      int    `json:"id" validate:"omitempty,gte=1"`
	CateId  int    `json:"cate_id" validate:"required,gte=1"`
	Type    int    `json:"type" validate:"required,gte=1"`
	Title   string `json:"title" validate:"required"`
	Url     string `json:"url"`
	Content string `json:"content" validate:"required"`
	Status  int    `json:"status"`
}

// 友链提交请求
type LinkRequest struct {
	Id   int    `json:"id" validate:"omitempty,gte=1"`
	Name string `json:"name" validate:"required"`
	Url  string `json:"url" validate:"required"`
	Desc string `json:"desc"`
}

// 心情提交请求
type MoodRequest struct {
	Id      int    `json:"id" validate:"omitempty,gte=1"`
	Content string `json:"content" validate:"required"`
}

// 提醒提交请求
type RemindRequest struct {
	Id      int    `json:"id" validate:"omitempty,gte=1"`
	Type    int    `json:"type"`
	Content string `json:"content" validate:"required"`
	Month   int    `json:"month"`
	Week    int    `json:"week"`
	Day     int    `json:"day"`
	Hour    int    `json:"hour"`
	Minute  int    `json:"minute"`
}

// 用户提交请求
type UserRequest struct {
	Id       int    `json:"id" validate:"omitempty,gte=1"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
	NickName string `json:"nick_name" validate:"required"`
	Email    string `json:"email"`
	Type     int    `json:"type" validate:"required,gte=1"`
}

// 分类创建
type CateCreateRequest struct {
	Name   string `json:"name" validate:"required"`
	Desc   string `json:"desc"`
	Domain string `json:"domain" validate:"required"`
}

// 分类更新
type CateUpdateRequest struct {
	Id     int    `json:"id" validate:"required,gte=1"`
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Domain string `json:"domain"`
}

// 友链创建
type LinkCreateRequest struct {
	Name string `json:"name" validate:"required"`
	Url  string `json:"url" validate:"required"`
	Desc string `json:"desc"`
}

// 友链更新
type LinkUpdateRequest struct {
	Id   int    `json:"id" validate:"required,gte=1"`
	Name string `json:"name"`
	Url  string `json:"url"`
	Desc string `json:"desc"`
}

// 心情创建
type MoodCreateRequest struct {
	Content string `json:"content" validate:"required"`
}

// 心情更新
type MoodUpdateRequest struct {
	Id      int    `json:"id" validate:"required,gte=1"`
	Content string `json:"content"`
}

// 文章创建（允许 title 为空以便返回自定义错误文案）
type ArticleCreateRequest struct {
	CateId  int    `json:"cate_id" validate:"required,gte=1"`
	Type    int    `json:"type" validate:"required,gte=1"`
	Title   string `json:"title" validate:"required"`
	Url     string `json:"url"`
	Content string `json:"content" validate:"required"`
}

// 文章更新
type ArticleUpdateRequest struct {
	Id      int    `json:"id" validate:"required,gte=1"`
	CateId  int    `json:"cate_id"`
	Type    int    `json:"type"`
	Title   string `json:"title"`
	Url     string `json:"url"`
	Content string `json:"content"`
	Status  int    `json:"status"`
}

// 提醒创建
type RemindCreateRequest struct {
	Type    int    `json:"type"`
	Content string `json:"content" validate:"required"`
	Month   int    `json:"month"`
	Week    int    `json:"week"`
	Day     int    `json:"day"`
	Hour    int    `json:"hour"`
	Minute  int    `json:"minute"`
}

// 提醒更新
type RemindUpdateRequest struct {
	Id      int    `json:"id" validate:"required,gte=1"`
	Type    int    `json:"type"`
	Content string `json:"content"`
	Month   int    `json:"month"`
	Week    int    `json:"week"`
	Day     int    `json:"day"`
	Hour    int    `json:"hour"`
	Minute  int    `json:"minute"`
}

// 用户创建（允许 password 为空以便返回自定义错误文案）
type UserCreateRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password"`
	NickName string `json:"nick_name" validate:"required"`
	Email    string `json:"email"`
	Type     int    `json:"type" validate:"required,gte=1"`
}

// 用户更新
type UserUpdateRequest struct {
	Id       int    `json:"id" validate:"required,gte=1"`
	Name     string `json:"name"`
	Password string `json:"password"`
	NickName string `json:"nick_name"`
	Email    string `json:"email"`
	Type     int    `json:"type"`
}
