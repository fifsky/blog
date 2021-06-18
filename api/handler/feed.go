package handler

import (
	"fmt"
	"time"

	"app/response"
	"github.com/goapt/gee"
	"github.com/goapt/golib/convert"
	"github.com/gorilla/feeds"

	"app/model"
)

var FeedGet gee.HandlerFunc = func(c *gee.Context) gee.Response {
	now := time.Now()
	options, err := model.GetOptions()
	if err != nil {
		response.Fail(c, 202, err)
	}

	feed := &feeds.Feed{
		Title:       options["site_name"],
		Link:        &feeds.Link{Href: "https://fifsky.com"},
		Description: options["site_desc"],
		Author:      &feeds.Author{Name: "fifsky", Email: "fifsky@gmail.com"},
		Created:     now,
	}

	cid := convert.StrTo(c.DefaultQuery("cid", "0")).MustInt()

	post := &model.Posts{}
	if cid > 0 {
		post.CateId = cid
	}

	posts, err := model.PostGetList(post, 1, 10, "", "")

	if err != nil {
		return response.Fail(c, 500, err)
	}

	for _, v := range posts {
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       v.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("https://fifsky.com/article/%d", v.Id)},
			Description: v.Content,
			Author:      &feeds.Author{Name: v.User.NickName, Email: "fifsky@gmail.com"},
			Created:     now,
		})
	}

	err = feed.WriteAtom(c.Writer)
	if err != nil {
		return response.Fail(c, 500, err)
	}
	return nil
}
