package store

import (
	"context"
	"testing"

	"app/pkg/dbunit"
	"app/store/model"
	"app/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_ListComments(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("comments", "posts"))
		s := New(db)

		// comments.yml 中 post_id=7 的评论共 5 条
		list, err := s.ListComments(context.Background(), 7)
		require.NoError(t, err)
		assert.Len(t, list, 5)

		// 验证按 created_at 正序：最早的是 id=4（2018-08-11）
		assert.Equal(t, 4, list[0].Id)

		// 验证回复的回复数据：id=12 pid=4 reply_name=站长
		var reply *model.Comment
		for i := range list {
			if list[i].Id == 12 {
				reply = &list[i]
			}
		}
		require.NotNil(t, reply)
		assert.Equal(t, 4, reply.Pid)
		assert.Equal(t, "站长", reply.ReplyName)

		// 验证 email/website 字段正确读取
		var webComment *model.Comment
		for i := range list {
			if list[i].Id == 11 {
				webComment = &list[i]
			}
		}
		require.NotNil(t, webComment)
		assert.Equal(t, "admin@example.com", webComment.Email)
		assert.Equal(t, "https://caixudong.com", webComment.Website)
	})

	t.Run("不存在的文章返回空列表", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("comments"))
			s := New(db)

			list, err := s.ListComments(context.Background(), 999)
			require.NoError(t, err)
			assert.Empty(t, list)
		})
	})
}

func TestStore_CreateComment(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("comments"))
		s := New(db)

		beforeList, err := s.ListComments(context.Background(), 7)
		require.NoError(t, err)
		beforeCount := len(beforeList)

		c := &model.Comment{
			PostId:    7,
			Pid:       0,
			Name:      "测试用户",
			Email:     "test@example.com",
			Website:   "https://example.com",
			Content:   "这是一条测试评论",
			ReplyName: "",
			IP:        "127.0.0.1",
		}

		id, err := s.CreateComment(context.Background(), c)
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))

		afterList, err := s.ListComments(context.Background(), 7)
		require.NoError(t, err)
		assert.Equal(t, beforeCount+1, len(afterList))

		// 通过 ID 找到刚创建的评论（不依赖排序）
		var created model.Comment
		for i := range afterList {
			if afterList[i].Id == int(id) {
				created = afterList[i]
			}
		}
		assert.Equal(t, "测试用户", created.Name)
		assert.Equal(t, "test@example.com", created.Email)
		assert.Equal(t, "https://example.com", created.Website)
		assert.Equal(t, "这是一条测试评论", created.Content)
	})
}

func TestStore_ListAllComments(t *testing.T) {
	tests := []struct {
		name      string
		keyword   string
		wantCount int
	}{
		{name: "全部", keyword: "", wantCount: 5},
		{name: "搜索昵称-站长", keyword: "站长", wantCount: 2},     // id=11 name + id=12 content
		{name: "搜索内容-时光飞逝", keyword: "时光飞逝", wantCount: 1}, // 仅 id=11 内容
		{name: "无匹配", keyword: "不存在的内容", wantCount: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("comments", "posts"))
				s := New(db)

				list, err := s.ListAllComments(context.Background(), tt.keyword, 1, 10)
				require.NoError(t, err)
				assert.Len(t, list, tt.wantCount)

				// 验证关联了文章标题
				if len(list) > 0 {
					assert.NotEmpty(t, list[0].PostTitle)
				}
			})
		})
	}
}

func TestStore_CountComments(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("comments"))
		s := New(db)

		total, err := s.CountComments(context.Background(), "")
		require.NoError(t, err)
		assert.Equal(t, 5, total)

		total, err = s.CountComments(context.Background(), "时光飞逝")
		require.NoError(t, err)
		assert.Equal(t, 1, total)
	})
}

func TestStore_ListNewComments(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("comments", "posts"))
		s := New(db)

		list, err := s.ListNewComments(context.Background(), 10)
		require.NoError(t, err)
		assert.Len(t, list, 5)

		// 验证按 created_at 倒序：最新的是 id=10（2018-09-18T20:05:11）
		assert.Equal(t, 10, list[0].Id)
		// 验证关联文章信息
		assert.Equal(t, "关于", list[0].PostTitle)
		assert.Equal(t, "about", list[0].PostUrl)
	})
}

func TestStore_DeleteComment(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("comments"))
		s := New(db)

		beforeTotal, err := s.CountComments(context.Background(), "")
		require.NoError(t, err)

		err = s.DeleteComment(context.Background(), []int{4, 9})
		require.NoError(t, err)

		afterTotal, err := s.CountComments(context.Background(), "")
		require.NoError(t, err)
		assert.Equal(t, beforeTotal-2, afterTotal)
	})

	t.Run("空ID不报错", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("comments"))
			s := New(db)

			err := s.DeleteComment(context.Background(), []int{})
			require.NoError(t, err)
		})
	})
}
