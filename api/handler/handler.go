package handler

import "github.com/google/wire"

type Handler struct {
	Article *Article
	Cate    *Cate
	Comment *Comment
	Common  *Common
	Link    *Link
	User    *User
	Mood    *Mood
	Remind  *Remind
	Setting *Setting
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
	wire.Struct(new(Handler), "*"))
