package repo

import (
	"app/provider/model"
	"github.com/ilibs/gosql/v2"
)

type Article struct {
	Base
	commRepo *Comment
}

func NewArticle(db *gosql.DB, commRepo *Comment) *Article {
	return &Article{
		Base:     Base{db: db},
		commRepo: commRepo,
	}
}

type UserPosts struct {
	model.Posts
	Cate       *model.Cates `json:"cate" db:"-" relation:"cate_id,id"`
	User       *model.Users `json:"user" db:"-" relation:"user_id,id"`
	CommentNum int          `json:"comment_num" db:"-"`
}

func (a *Article) GetUserPost(id int, url string) (*UserPosts, error) {
	post := &UserPosts{}
	post.Id = id
	post.Url = url

	err := a.db.Model(post).Where("status = 1").Get()

	if err != nil {
		return nil, err
	}
	return post, nil
}

func (a *Article) PostPrev(id int) (*model.Posts, error) {
	m := &model.Posts{
		Type: 1,
	}
	err := gosql.Model(m).Where("id < ? and status = 1", id).OrderBy("id desc").Limit(1).Get()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (a *Article) PostNext(id int) (*model.Posts, error) {
	m := &model.Posts{
		Type: 1,
	}
	err := gosql.Model(m).Where("id > ? and status = 1", id).OrderBy("id asc").Limit(1).Get()

	if err != nil {
		return nil, err
	}
	return m, nil
}

func (a *Article) PostArchive() ([]map[string]string, error) {
	m := make([]map[string]string, 0)
	result, err := gosql.Queryx("select ym,count(ym) total from (select DATE_FORMAT(created_at,'%Y/%m') as ym from posts where type = 1 and status = 1) s group by ym order by ym desc")

	if err != nil {
		return nil, err
	}

	for result.Next() {
		var ym, total string
		_ = result.Scan(&ym, &total)
		m = append(m, map[string]string{
			"ym":    ym,
			"total": total,
		})
	}

	return m, err
}

func (a *Article) PostGetList(p *model.Posts, start int, num int, artdate, keyword string) ([]*UserPosts, error) {
	var posts = make([]*UserPosts, 0)
	start = (start - 1) * num

	args := make([]interface{}, 0)
	where := "status = 1"

	if p.CateId > 0 {
		where += " and cate_id = ?"
		args = append(args, p.CateId)
	}

	if p.Type > 0 {
		where += " and type = ?"
		args = append(args, p.Type)
	}

	if artdate != "" {
		where += " and DATE_FORMAT(created_at,'%Y-%m') = ?"
		args = append(args, artdate)
	}

	if keyword != "" {
		where += " and title like ?"
		args = append(args, "%"+keyword+"%")
	}

	err := gosql.Model(&posts).Where(where, args...).Limit(num).Offset(start).OrderBy("id desc").All()

	if err != nil {
		return nil, err
	}

	postIds := make([]int, 0)

	for _, v := range posts {
		postIds = append(postIds, v.Id)
	}

	cm, err := a.commRepo.PostCommentNum(postIds)

	if err != nil {
		return nil, err
	}

	for _, v := range posts {
		if c, ok := cm[v.Id]; ok {
			v.CommentNum = c
		}
	}

	return posts, err
}
