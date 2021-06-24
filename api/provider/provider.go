package provider

import (
	"app/provider/repo"
	"github.com/google/wire"
)

var RepoSet = wire.NewSet(
	repo.NewArticle,
	repo.NewCate,
	repo.NewComment,
	repo.NewLink,
	repo.NewMood,
	repo.NewRemind,
	repo.NewSetting,
	repo.NewUser,
)
