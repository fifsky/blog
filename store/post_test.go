package store

import (
	"context"
	"testing"
	"time"

	"app/pkg/dbunit"
	"app/store/model"
	"app/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArticle_PostPrev(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts"))
		s := New(db)
		ret, err := s.PrevPost(context.Background(), 7)
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.Equal(t, ret.Id, 8)
	})
}

func TestArticle_PostNext(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts"))
		s := New(db)
		ret, err := s.NextPost(context.Background(), 7)
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.Equal(t, ret.Id, 4)
	})
}

func TestArticle_PostArchive(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts"))
		s := New(db)
		ret, err := s.PostArchive(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.True(t, len(ret) > 0)
	})
}

func TestArticle_PostGetList(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts", "users", "cates"))
		s := New(db)

		p := &model.Post{
			CateId: 1,
		}

		ret, err := s.ListPost(context.Background(), p, 1, 1, "", "", "")
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.True(t, len(ret) > 0)
	})
}

func TestArticle_GetPost(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		url     string
		wantErr bool
		wantId  int
	}{
		{name: "按ID查询", id: 4, url: "", wantErr: false, wantId: 4},
		{name: "按URL查询", id: 0, url: "about", wantErr: false, wantId: 7},
		{name: "不存在的ID", id: 999, url: "", wantErr: true},
		{name: "不存在的URL", id: 0, url: "nonexistent", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts"))
				s := New(db)
				ret, err := s.GetPost(context.Background(), tt.id, tt.url)
				if tt.wantErr {
					require.Error(t, err)
					assert.Nil(t, ret)
					return
				}
				require.NoError(t, err)
				assert.NotNil(t, ret)
				assert.Equal(t, tt.wantId, ret.Id)
			})
		})
	}
}

func TestArticle_IncrementPostViewNum(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts"))
		s := New(db)

		// 获取初始浏览量
		before, err := s.GetPost(context.Background(), 4, "")
		require.NoError(t, err)

		err = s.IncrementPostViewNum(context.Background(), 4)
		require.NoError(t, err)

		after, err := s.GetPost(context.Background(), 4, "")
		require.NoError(t, err)
		assert.Equal(t, before.ViewNum+1, after.ViewNum)
	})
}

func TestArticle_GetPostDaysInMonth(t *testing.T) {
	tests := []struct {
		name     string
		year     int
		month    int
		wantDays []int32
	}{
		{name: "2012年9月有文章", year: 2012, month: 9, wantDays: []int32{10, 28}},
		{name: "2012年10月有文章", year: 2012, month: 10, wantDays: []int32{28}},
		{name: "2020年3月有文章", year: 2020, month: 3, wantDays: nil},
		{name: "无文章的月份", year: 2025, month: 1, wantDays: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts"))
				s := New(db)
				days, err := s.GetPostDaysInMonth(context.Background(), tt.year, tt.month)
				require.NoError(t, err)
				if len(tt.wantDays) == 0 {
					assert.Empty(t, days)
					return
				}
				assert.ElementsMatch(t, tt.wantDays, days)
			})
		})
	}
}

func TestArticle_CountPosts(t *testing.T) {
	tests := []struct {
		name      string
		post      *model.Post
		artdate   string
		keyword   string
		tag       string
		wantCount int
	}{
		{name: "全部文章", post: &model.Post{}, artdate: "", keyword: "", tag: "", wantCount: 3},
		{name: "按分类", post: &model.Post{CateId: 1}, artdate: "", keyword: "", tag: "", wantCount: 3},
		{name: "按类型-文章", post: &model.Post{Type: 1}, artdate: "", keyword: "", tag: "", wantCount: 2},
		{name: "按类型-页面", post: &model.Post{Type: 2}, artdate: "", keyword: "", tag: "", wantCount: 1},
		{name: "按日期", post: &model.Post{}, artdate: "2012-09", keyword: "", tag: "", wantCount: 2},
		{name: "按关键字", post: &model.Post{}, artdate: "", keyword: "fifsky", tag: "", wantCount: 1},
		{name: "无匹配关键字", post: &model.Post{}, artdate: "", keyword: "不存在的关键字", tag: "", wantCount: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts"))
				s := New(db)
				count, err := s.CountPosts(context.Background(), tt.post, tt.artdate, tt.keyword, tt.tag)
				require.NoError(t, err)
				assert.Equal(t, tt.wantCount, count)
			})
		})
	}
}

func TestArticle_ListPost(t *testing.T) {
	tests := []struct {
		name    string
		post    *model.Post
		artdate string
		keyword string
		tag     string
		start   int
		num     int
		wantLen int
	}{
		{name: "全部文章第一页", post: &model.Post{}, artdate: "", keyword: "", tag: "", start: 1, num: 10, wantLen: 3},
		{name: "按分类", post: &model.Post{CateId: 1}, artdate: "", keyword: "", tag: "", start: 1, num: 10, wantLen: 3},
		{name: "按类型-文章", post: &model.Post{Type: 1}, artdate: "", keyword: "", tag: "", start: 1, num: 10, wantLen: 2},
		{name: "按类型-页面", post: &model.Post{Type: 2}, artdate: "", keyword: "", tag: "", start: 1, num: 10, wantLen: 1},
		{name: "按年月过滤", post: &model.Post{}, artdate: "2012-09", keyword: "", tag: "", start: 1, num: 10, wantLen: 2},
		{name: "按关键字", post: &model.Post{}, artdate: "", keyword: "fifsky", tag: "", start: 1, num: 10, wantLen: 1},
		{name: "无匹配", post: &model.Post{}, artdate: "", keyword: "不存在", tag: "", start: 1, num: 10, wantLen: 0},
		{name: "分页-每页1条第1页", post: &model.Post{}, artdate: "", keyword: "", tag: "", start: 1, num: 1, wantLen: 1},
		{name: "分页-每页1条第4页无数据", post: &model.Post{}, artdate: "", keyword: "", tag: "", start: 4, num: 1, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts"))
				s := New(db)
				ret, err := s.ListPost(context.Background(), tt.post, tt.start, tt.num, tt.artdate, tt.keyword, tt.tag)
				require.NoError(t, err)
				assert.Len(t, ret, tt.wantLen)
			})
		})
	}
}

func TestArticle_GetCateByDomain(t *testing.T) {
	tests := []struct {
		name     string
		domain   string
		wantErr  bool
		wantName string
	}{
		{name: "存在domain", domain: "default", wantErr: false, wantName: "默认分类"},
		{name: "不存在domain", domain: "nonexistent", wantErr: true},
		{name: "空domain", domain: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("cates"))
				s := New(db)
				ret, err := s.GetCateByDomain(context.Background(), tt.domain)
				if tt.wantErr {
					require.Error(t, err)
					assert.Nil(t, ret)
					return
				}
				require.NoError(t, err)
				assert.NotNil(t, ret)
				assert.Equal(t, tt.wantName, ret.Name)
			})
		})
	}
}

func TestArticle_CreatePost(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts"))
		s := New(db)

		now := time.Now()
		tags := model.Tags{"test", "go"}
		p := &model.Post{
			CateId:    1,
			Type:      1,
			UserId:    1,
			Title:     "新文章",
			Url:       "new-post",
			Content:   "内容",
			Tags:      tags,
			Status:    model.PostStatusActive,
			CreatedAt: now,
			UpdatedAt: now,
		}

		id, err := s.CreatePost(context.Background(), p)
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))

		// 验证创建的文章
		got, err := s.GetPost(context.Background(), int(id), "")
		require.NoError(t, err)
		assert.Equal(t, "新文章", got.Title)
		assert.Equal(t, "new-post", got.Url)
		assert.Equal(t, "内容", got.Content)
		assert.Equal(t, model.PostStatusActive, got.Status)
	})
}

func TestArticle_UpdatePost(t *testing.T) {
	tests := []struct {
		name   string
		update *model.UpdatePost
		check  func(t *testing.T, p *model.Post)
	}{
		{
			name: "更新标题",
			update: &model.UpdatePost{
				Id:    4,
				Title: new("新标题"),
			},
			check: func(t *testing.T, p *model.Post) {
				assert.Equal(t, "新标题", p.Title)
			},
		},
		{
			name: "更新内容",
			update: &model.UpdatePost{
				Id:      4,
				Content: new("新内容"),
			},
			check: func(t *testing.T, p *model.Post) {
				assert.Equal(t, "新内容", p.Content)
			},
		},
		{
			name: "更新URL",
			update: &model.UpdatePost{
				Id:  4,
				Url: new("new-url"),
			},
			check: func(t *testing.T, p *model.Post) {
				assert.Equal(t, "new-url", p.Url)
			},
		},
		{
			name: "更新状态为草稿",
			update: &model.UpdatePost{
				Id:     4,
				Status: new(model.PostStatusDraft),
			},
			check: func(t *testing.T, p *model.Post) {
				assert.Equal(t, model.PostStatusDraft, p.Status)
			},
		},
		{
			name:   "空更新不报错",
			update: &model.UpdatePost{Id: 4},
			check: func(t *testing.T, p *model.Post) {
				assert.Equal(t, "fifsky blog for php!", p.Title)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts"))
				s := New(db)
				err := s.UpdatePost(context.Background(), tt.update)
				require.NoError(t, err)

				got, err := s.GetPost(context.Background(), tt.update.Id, "")
				require.NoError(t, err)
				tt.check(t, got)
			})
		})
	}
}

func TestArticle_SoftDeletePost(t *testing.T) {
	t.Run("软删除单个文章", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts"))
			s := New(db)
			err := s.SoftDeletePost(context.Background(), []int{4})
			require.NoError(t, err)

			// 验证状态变为 DELETED（GetPost 不按状态过滤，仍可查到）
			got, err := s.GetPost(context.Background(), 4, "")
			require.NoError(t, err)
			assert.Equal(t, model.PostStatusDeleted, got.Status)

			// 验证 ACTIVE 文章列表中不再包含
			count, err := s.CountPosts(context.Background(), &model.Post{}, "", "", "")
			require.NoError(t, err)
			assert.Equal(t, 2, count)
		})
	})

	t.Run("软删除多个文章", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts"))
			s := New(db)
			err := s.SoftDeletePost(context.Background(), []int{4, 7, 8})
			require.NoError(t, err)

			count, err := s.CountPosts(context.Background(), &model.Post{}, "", "", "")
			require.NoError(t, err)
			assert.Equal(t, 0, count)
		})
	})

	t.Run("空ID列表不报错", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts"))
			s := New(db)
			err := s.SoftDeletePost(context.Background(), []int{})
			require.NoError(t, err)
		})
	})
}

func TestArticle_RestorePost(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts"))
		s := New(db)

		// 先软删除
		err := s.SoftDeletePost(context.Background(), []int{4})
		require.NoError(t, err)

		// 恢复
		err = s.RestorePost(context.Background(), 4)
		require.NoError(t, err)

		// 验证状态变为 DRAFT
		got, err := s.GetPost(context.Background(), 4, "")
		require.NoError(t, err)
		assert.Equal(t, model.PostStatusDraft, got.Status)
	})
}

func TestArticle_ListPostForAdmin(t *testing.T) {
	tests := []struct {
		name    string
		post    *model.Post
		keyword string
		start   int
		num     int
		wantLen int
	}{
		{name: "全部文章", post: &model.Post{}, keyword: "", start: 1, num: 10, wantLen: 3},
		{name: "按分类", post: &model.Post{CateId: 1}, keyword: "", start: 1, num: 10, wantLen: 3},
		{name: "按状态", post: &model.Post{Status: model.PostStatusActive}, keyword: "", start: 1, num: 10, wantLen: 3},
		{name: "按关键字", post: &model.Post{}, keyword: "fifsky", start: 1, num: 10, wantLen: 1},
		{name: "无匹配", post: &model.Post{}, keyword: "不存在", start: 1, num: 10, wantLen: 0},
		{name: "分页", post: &model.Post{}, keyword: "", start: 1, num: 1, wantLen: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts"))
				s := New(db)
				ret, err := s.ListPostForAdmin(context.Background(), tt.post, tt.start, tt.num, tt.keyword)
				require.NoError(t, err)
				assert.Len(t, ret, tt.wantLen)
			})
		})
	}
}

func TestArticle_CountPostsForAdmin(t *testing.T) {
	tests := []struct {
		name      string
		post      *model.Post
		keyword   string
		wantCount int
	}{
		{name: "全部", post: &model.Post{}, keyword: "", wantCount: 3},
		{name: "按分类", post: &model.Post{CateId: 1}, keyword: "", wantCount: 3},
		{name: "按状态", post: &model.Post{Status: model.PostStatusActive}, keyword: "", wantCount: 3},
		{name: "按关键字", post: &model.Post{}, keyword: "fifsky", wantCount: 1},
		{name: "无匹配", post: &model.Post{}, keyword: "不存在", wantCount: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts"))
				s := New(db)
				count, err := s.CountPostsForAdmin(context.Background(), tt.post, tt.keyword)
				require.NoError(t, err)
				assert.Equal(t, tt.wantCount, count)
			})
		})
	}
}

func TestArticle_DestroyPost(t *testing.T) {
	t.Run("彻底删除已软删除的文章", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts"))
			s := New(db)

			// 先软删除
			err := s.SoftDeletePost(context.Background(), []int{4})
			require.NoError(t, err)

			// 彻底删除
			err = s.DestroyPost(context.Background(), []int{4})
			require.NoError(t, err)

			// 验证已彻底删除
			_, err = s.GetPost(context.Background(), 4, "")
			require.Error(t, err)
		})
	})

	t.Run("未软删除的文章不会被彻底删除", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts"))
			s := New(db)

			// 直接尝试彻底删除 ACTIVE 状态的文章
			err := s.DestroyPost(context.Background(), []int{4})
			require.NoError(t, err)

			// 文章仍应存在
			got, err := s.GetPost(context.Background(), 4, "")
			require.NoError(t, err)
			assert.NotNil(t, got)
		})
	})

	t.Run("空ID列表不报错", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("posts"))
			s := New(db)
			err := s.DestroyPost(context.Background(), []int{})
			require.NoError(t, err)
		})
	})
}
