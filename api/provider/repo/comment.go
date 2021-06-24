package repo

import (
	"strings"

	"app/provider/model"
	"github.com/goapt/golib/convert"
	"github.com/ilibs/gosql/v2"
)

type Comment struct {
	Base
}

func NewComment(db *gosql.DB) *Comment {
	return &Comment{Base: Base{db: db}}
}

func (o *Comment) PostComments(postId, start, num int) ([]*model.Comments, error) {
	var m = make([]*model.Comments, 0)
	start = (start - 1) * num
	err := o.db.Model(&m).Where("post_id = ?", postId).OrderBy("id asc").Limit(num).Offset(start).All()
	if err != nil {
		return nil, err
	}
	return m, nil
}

type NewComments struct {
	model.Comments
	Type         int    `json:"type" db:"type"`
	ArticleTitle string `json:"article_title" db:"title"`
	Url          string `json:"url" db:"url"`
}

func (o *Comment) NewComments() ([]*NewComments, error) {
	var m = make([]*NewComments, 0)
	err := o.db.Select(&m, "select p.type,p.title,c.* from comments c left join posts p on c.post_id = p.id order by c.id desc limit 10")
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (o *Comment) CommentList(start, num int) ([]*NewComments, error) {
	var m = make([]*NewComments, 0)
	start = (start - 1) * num
	err := o.db.Select(&m, "select p.type,p.title,c.* from comments c left join posts p on c.post_id = p.id order by c.id desc limit ?,?", start, num)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (o *Comment) PostCommentNum(postId []int) (map[int]int, error) {
	m := make(map[int]int)

	postIds := make([]string, 0)
	for _, v := range postId {
		postIds = append(postIds, convert.ToStr(v))
	}

	if len(postIds) == 0 {
		return m, nil
	}

	rows, err := o.db.Queryx("select count(*) comment_num,post_id from comments where post_id in(" + strings.Join(postIds, ",") + ") group by post_id")

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
