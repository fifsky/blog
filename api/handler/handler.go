package handler

import "github.com/google/wire"

type Handler struct {
	Article  *Article
	Cate     *Cate
	Comment  *Comment
	Common   *Common
	Link     *Link
	User     *User
	Mood     *Mood
	Remind   *Remind
	Setting  *Setting
	DingTalk *DingTalk
}

var ProviderSet = wire.NewSet(
	NewArticle,
	NewCate,
	NewComment,
	NewCommon,
	NewLink,
	NewUser,
	NewMood,
	NewRemind,
	NewSetting,
	NewDingTalk,
	wire.Struct(new(Handler), "*"))
