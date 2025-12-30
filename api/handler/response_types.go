package handler

import "time"

type UserItem struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	NickName  string    `json:"nick_name"`
	Email     string    `json:"email"`
	Status    int       `json:"status"`
	Type      int       `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type IDResponse struct {
	Id int `json:"id"`
}

type UserSummary struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	NickName string `json:"nick_name"`
}

type CateSummary struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

type CateMenuItem struct {
	Url     string `json:"url"`
	Content string `json:"content"`
}

type CateListItem struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Desc      string    `json:"desc"`
	Domain    string    `json:"domain"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Num       int       `json:"num"`
}

type CateListResponse struct {
	List      []CateListItem `json:"list"`
	PageTotal int            `json:"page_total"`
}

type DateArchiveItem struct {
	Url     string `json:"url"`
	Content string `json:"content"`
}

type ArticleItem struct {
	Id        int          `json:"id"`
	CateId    int          `json:"cate_id"`
	Type      int          `json:"type"`
	UserId    int          `json:"user_id"`
	Title     string       `json:"title"`
	Url       string       `json:"url"`
	Content   string       `json:"content"`
	Status    int          `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	User      *UserSummary `json:"user,omitempty"`
	Cate      *CateSummary `json:"cate,omitempty"`
}

type ArticleListResponse struct {
	List      []ArticleItem `json:"list"`
	PageTotal int           `json:"page_total"`
}

type PrevNextItem struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type PrevNextResponse struct {
	Prev PrevNextItem `json:"prev"`
	Next PrevNextItem `json:"next"`
}

type LinkItem struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	Desc      string    `json:"desc"`
	CreatedAt time.Time `json:"created_at"`
}

type LinkMenuItem struct {
	Url     string `json:"url"`
	Content string `json:"content"`
}

type LinkListResponse struct {
	List      []LinkItem `json:"list"`
	PageTotal int        `json:"page_total"`
}

type MoodItem struct {
	Id        int          `json:"id"`
	Content   string       `json:"content"`
	User      *UserSummary `json:"user,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
}

type MoodListResponse struct {
	List      []MoodItem `json:"list"`
	PageTotal int        `json:"page_total"`
}

type RemindItem struct {
	Id        int       `json:"id"`
	Type      int       `json:"type"`
	Content   string    `json:"content"`
	Month     int       `json:"month"`
	Week      int       `json:"week"`
	Day       int       `json:"day"`
	Hour      int       `json:"hour"`
	Minute    int       `json:"minute"`
	Status    int       `json:"status"`
	NextTime  time.Time `json:"next_time"`
	CreatedAt time.Time `json:"created_at"`
}

type RemindListResponse struct {
	List      []RemindItem `json:"list"`
	PageTotal int          `json:"page_total"`
}

type OptionsResponse map[string]string

type UserListReponse struct {
	List      []UserItem `json:"list"`
	PageTotal int        `json:"page_total"`
}
