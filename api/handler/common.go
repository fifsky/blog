package handler

import (
	"app/provider/model"
	"github.com/goapt/gee"
	"github.com/ilibs/identicon"
)

const (
	mooodTag = "#心情#"
)

func getLoginUser(c *gee.Context) *model.Users {
	if u, ok := c.Get("userInfo"); ok {
		return u.(*model.Users)
	}
	return nil
}

type Common struct {
}

func NewCommon() *Common {
	return &Common{}
}

func (m *Common) Avatar(c *gee.Context) gee.Response {
	name := c.DefaultQuery("name", "default")

	// New Generator: Rehuse
	ig, err := identicon.New(
		"fifsky", // Namespace
		5,        // Number of blocks (Size)
		5,        // Density
	)

	if err != nil {
		panic(err) // Invalid Size or Density
	}

	ii, err := ig.Draw(name) // Generate an IdentIcon

	if err != nil {
		return nil
	}
	// Takes the size in pixels and any io.Writer
	_ = ii.Png(300, c.Writer) // 300px * 300px
	return nil
}
