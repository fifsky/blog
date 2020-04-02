package model

import (
	"strings"
	"time"

	"github.com/goapt/golib/convert"

	"github.com/ilibs/gosql/v2"
)

type Comments struct {
	Id        int       `form:"id" json:"id" db:"id"`
	PostId    int       `form:"post_id" json:"post_id" db:"post_id"`
	Pid       int       `form:"pid" json:"pid" db:"pid"`
	Name      string    `form:"name" json:"name" db:"name"`
	Content   string    `form:"content" json:"content" db:"content"`
	IP        string    `form:"-" json:"ip" db:"ip"`
	CreatedAt time.Time `form:"-" json:"created_at" time_format:"2006-01-02 15:04:05" db:"created_at"`
}

func (c *Comments) TableName() string {
	return "comments"
}

func (c *Comments) PK() string {
	return "id"
}

func (c *Comments) AfterChange() {
	Cache.Delete("new-comments")
}

func PostComments(postId, start, num int) ([]*Comments, error) {
	var m = make([]*Comments, 0)
	start = (start - 1) * num
	err := gosql.Model(&m).Where("post_id = ?", postId).OrderBy("id asc").Limit(num).Offset(start).All()
	if err != nil {
		return nil, err
	}
	return m, nil
}

type NewComment struct {
	Comments
	Type         int    `json:"type" db:"type"`
	ArticleTitle string `json:"article_title" db:"title"`
	Url          string `json:"url" db:"url"`
}

func NewComments() ([]*NewComment, error) {
	if v, ok := Cache.Get("new-comments"); ok {
		return v.([]*NewComment), nil
	}

	var m = make([]*NewComment, 0)
	err := gosql.Select(&m, "select p.type,p.title,c.* from comments c left join posts p on c.post_id = p.id order by c.id desc limit 10")
	if err != nil {
		return nil, err
	}

	Cache.Set("new-comments", m, 1*time.Hour)

	return m, nil
}

func CommentList(start, num int) ([]*NewComment, error) {
	var m = make([]*NewComment, 0)
	start = (start - 1) * num
	err := gosql.Select(&m, "select p.type,p.title,c.* from comments c left join posts p on c.post_id = p.id order by c.id desc limit ?,?", start, num)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func PostCommentNum(postId []int) (map[int]int, error) {
	m := make(map[int]int)

	postIds := make([]string, 0)
	for _, v := range postId {
		postIds = append(postIds, convert.ToStr(v))
	}

	if len(postIds) == 0 {
		return m, nil
	}

	rows, err := gosql.Queryx("select count(*) comment_num,post_id from comments where post_id in(" + strings.Join(postIds, ",") + ") group by post_id")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var commentNum, postId int
		err := rows.Scan(&commentNum, &postId)
		if err != nil {
			return nil, err
		}

		m[postId] = commentNum
	}

	return m, nil
}
