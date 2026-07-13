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

func TestMood_MoodGetList(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("moods", "users"))
		s := New(db)
		ret, err := s.ListMood(context.Background(), 1, 1)
		assert.NoError(t, err)
		assert.NotNil(t, ret)
		assert.True(t, len(ret) > 0)
	})
}

func TestMood_ListMood(t *testing.T) {
	tests := []struct {
		name      string
		start     int
		num       int
		wantEmpty bool
	}{
		{name: "全部第一页", start: 1, num: 10, wantEmpty: false},
		{name: "每页1条第1页", start: 1, num: 1, wantEmpty: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("moods"))
				s := New(db)
				ret, err := s.ListMood(context.Background(), tt.start, tt.num)
				require.NoError(t, err)
				if tt.wantEmpty {
					assert.Empty(t, ret)
					return
				}
				assert.NotEmpty(t, ret)
				assert.LessOrEqual(t, len(ret), tt.num)
			})
		})
	}
}

func TestMood_RandomMood(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("moods"))
		s := New(db)
		ret, err := s.RandomMood(context.Background())
		require.NoError(t, err)
		assert.NotNil(t, ret)
		assert.True(t, ret.Id > 0)
		assert.NotEmpty(t, ret.Content)
	})

	t.Run("空表返回错误", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), nil)
			s := New(db)
			ret, err := s.RandomMood(context.Background())
			require.Error(t, err)
			assert.Nil(t, ret)
		})
	})
}

func TestMood_CountMoodTotal(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("moods"))
		s := New(db)
		total, err := s.CountMoodTotal(context.Background())
		require.NoError(t, err)
		assert.Greater(t, total, 0)
	})

	t.Run("空表", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), nil)
			s := New(db)
			total, err := s.CountMoodTotal(context.Background())
			require.NoError(t, err)
			assert.Equal(t, 0, total)
		})
	})
}

func TestMood_CreateMood(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), nil)
		s := New(db)

		beforeTotal, err := s.CountMoodTotal(context.Background())
		require.NoError(t, err)

		now := time.Now()
		md := &model.Mood{
			Content:   "今天心情不错",
			UserId:    1,
			CreatedAt: now,
			UpdatedAt: now,
		}

		id, err := s.CreateMood(context.Background(), md)
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))

		afterTotal, err := s.CountMoodTotal(context.Background())
		require.NoError(t, err)
		assert.Equal(t, beforeTotal+1, afterTotal)

		// 验证创建的数据
		list, err := s.ListMood(context.Background(), 1, 10)
		require.NoError(t, err)
		assert.NotEmpty(t, list)

		var created model.Mood
		for _, m := range list {
			if m.Id == int(id) {
				created = m
			}
		}
		assert.Equal(t, "今天心情不错", created.Content)
		assert.Equal(t, 1, created.UserId)
	})
}

func TestMood_UpdateMood(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("moods"))
		s := New(db)

		err := s.UpdateMood(context.Background(), &model.UpdateMood{
			Id:      1,
			Content: new("更新后的心情"),
		})
		require.NoError(t, err)

		// 验证更新
		list, err := s.ListMood(context.Background(), 1, 10)
		require.NoError(t, err)
		for _, m := range list {
			if m.Id == 1 {
				assert.Equal(t, "更新后的心情", m.Content)
				return
			}
		}
		require.Fail(t, "未找到 id=1 的心情")
	})
}

func TestMood_DeleteMood(t *testing.T) {
	t.Run("删除单个心情", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("moods"))
			s := New(db)

			beforeTotal, err := s.CountMoodTotal(context.Background())
			require.NoError(t, err)

			err = s.DeleteMood(context.Background(), []int{1})
			require.NoError(t, err)

			afterTotal, err := s.CountMoodTotal(context.Background())
			require.NoError(t, err)
			assert.Equal(t, beforeTotal-1, afterTotal)
		})
	})

	t.Run("删除多个心情", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("moods"))
			s := New(db)

			beforeTotal, err := s.CountMoodTotal(context.Background())
			require.NoError(t, err)

			err = s.DeleteMood(context.Background(), []int{1, 2})
			require.NoError(t, err)

			afterTotal, err := s.CountMoodTotal(context.Background())
			require.NoError(t, err)
			assert.Less(t, afterTotal, beforeTotal)
		})
	})

	t.Run("空ID列表不报错", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("moods"))
			s := New(db)
			err := s.DeleteMood(context.Background(), []int{})
			require.NoError(t, err)
		})
	})

	t.Run("删除不存在ID不报错", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("moods"))
			s := New(db)

			beforeTotal, err := s.CountMoodTotal(context.Background())
			require.NoError(t, err)

			err = s.DeleteMood(context.Background(), []int{999})
			require.NoError(t, err)

			// 验证总数不变
			afterTotal, err := s.CountMoodTotal(context.Background())
			require.NoError(t, err)
			assert.Equal(t, beforeTotal, afterTotal)
		})
	})
}
