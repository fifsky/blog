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

func TestRemind_RemindGetList(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("reminds"))
		s := New(db)
		ret, err := s.ListRemind(context.Background(), 1, 1)
		assert.NoError(t, err)
		assert.NotNil(t, ret)
	})
}

func TestRemind_ListRemind(t *testing.T) {
	tests := []struct {
		name    string
		start   int
		num     int
		wantLen int
	}{
		{name: "全部第一页", start: 1, num: 10, wantLen: 2},
		{name: "每页1条第1页", start: 1, num: 1, wantLen: 1},
		{name: "每页1条第3页无数据", start: 3, num: 1, wantLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("reminds"))
				s := New(db)
				ret, err := s.ListRemind(context.Background(), tt.start, tt.num)
				require.NoError(t, err)
				assert.Len(t, ret, tt.wantLen)
			})
		})
	}
}

func TestRemind_GetRemind(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		wantErr  bool
		wantName string
	}{
		{name: "存在提醒", id: 8, wantErr: false, wantName: "提醒！！！"},
		{name: "存在提醒2", id: 9, wantErr: false, wantName: "生日快乐"},
		{name: "不存在提醒", id: 999, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("reminds"))
				s := New(db)
				ret, err := s.GetRemind(context.Background(), tt.id)
				if tt.wantErr {
					require.Error(t, err)
					assert.Nil(t, ret)
					return
				}
				require.NoError(t, err)
				assert.NotNil(t, ret)
				assert.Equal(t, tt.wantName, ret.Content)
			})
		})
	}
}

func TestRemind_RemindAll(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("reminds"))
		s := New(db)
		ret, err := s.RemindAll(context.Background())
		require.NoError(t, err)
		// fixture 中2条都是ACTIVE状态，应全部返回
		assert.Len(t, ret, 2)
		for _, r := range ret {
			assert.Contains(t, []model.RemindStatus{model.RemindStatusActive, model.RemindStatusPending}, r.Status)
		}
	})

	t.Run("只返回ACTIVE和PENDING状态", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), nil)
			s := New(db)

			now := time.Now()
			// 创建 ACTIVE 状态
			_, err := s.CreateRemind(context.Background(), &model.Remind{
				Cron: "* * * * *", Content: "active", Status: model.RemindStatusActive,
				NextTime: now, CreatedAt: now, UpdatedAt: now,
			})
			require.NoError(t, err)

			// 创建 PENDING 状态
			_, err = s.CreateRemind(context.Background(), &model.Remind{
				Cron: "* * * * *", Content: "pending", Status: model.RemindStatusPending,
				NextTime: now, CreatedAt: now, UpdatedAt: now,
			})
			require.NoError(t, err)

			// 创建 DONE 状态
			_, err = s.CreateRemind(context.Background(), &model.Remind{
				Cron: "* * * * *", Content: "done", Status: model.RemindStatusDone,
				NextTime: now, CreatedAt: now, UpdatedAt: now,
			})
			require.NoError(t, err)

			ret, err := s.RemindAll(context.Background())
			require.NoError(t, err)
			assert.Len(t, ret, 2) // 只有 ACTIVE 和 PENDING
		})
	})
}

func TestRemind_CountRemindTotal(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("reminds"))
		s := New(db)
		total, err := s.CountRemindTotal(context.Background())
		require.NoError(t, err)
		assert.Equal(t, 2, total)
	})

	t.Run("空表", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), nil)
			s := New(db)
			total, err := s.CountRemindTotal(context.Background())
			require.NoError(t, err)
			assert.Equal(t, 0, total)
		})
	})
}

func TestRemind_UpdateRemindStatus(t *testing.T) {
	tests := []struct {
		name   string
		id     int
		status model.RemindStatus
	}{
		{name: "更新为DONE", id: 8, status: model.RemindStatusDone},
		{name: "更新为PENDING", id: 8, status: model.RemindStatusPending},
		{name: "更新为ACTIVE", id: 8, status: model.RemindStatusActive},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("reminds"))
				s := New(db)
				err := s.UpdateRemindStatus(context.Background(), tt.id, tt.status)
				require.NoError(t, err)

				got, err := s.GetRemind(context.Background(), tt.id)
				require.NoError(t, err)
				assert.Equal(t, tt.status, got.Status)
			})
		})
	}
}

func TestRemind_UpdateRemindNextTime(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("reminds"))
		s := New(db)

		newTime := time.Date(2025, 12, 31, 23, 59, 59, 0, time.Local)
		err := s.UpdateRemindNextTime(context.Background(), 8, newTime)
		require.NoError(t, err)

		got, err := s.GetRemind(context.Background(), 8)
		require.NoError(t, err)
		assert.Equal(t, newTime, got.NextTime)
	})
}

func TestRemind_CreateRemind(t *testing.T) {
	dbunit.New(t, func(d *dbunit.DBUnit) {
		db := d.NewDatabase(testutil.Schema(), nil)
		s := New(db)

		beforeTotal, err := s.CountRemindTotal(context.Background())
		require.NoError(t, err)

		now := time.Now()
		nextTime := now.Add(24 * time.Hour)
		md := &model.Remind{
			Cron:      "* * * * *",
			Content:   "新提醒",
			Status:    model.RemindStatusActive,
			NextTime:  nextTime,
			CreatedAt: now,
			UpdatedAt: now,
		}

		id, err := s.CreateRemind(context.Background(), md)
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))

		afterTotal, err := s.CountRemindTotal(context.Background())
		require.NoError(t, err)
		assert.Equal(t, beforeTotal+1, afterTotal)

		// 验证创建的数据
		got, err := s.GetRemind(context.Background(), int(id))
		require.NoError(t, err)
		assert.Equal(t, "新提醒", got.Content)
		assert.Equal(t, "* * * * *", got.Cron)
		assert.Equal(t, model.RemindStatusActive, got.Status)
	})
}

func TestRemind_UpdateRemind(t *testing.T) {
	tests := []struct {
		name   string
		update *model.UpdateRemind
		check  func(t *testing.T, r *model.Remind)
	}{
		{
			name: "更新内容",
			update: &model.UpdateRemind{
				Id:      8,
				Content: new("新内容"),
			},
			check: func(t *testing.T, r *model.Remind) {
				assert.Equal(t, "新内容", r.Content)
			},
		},
		{
			name: "更新cron",
			update: &model.UpdateRemind{
				Id:   8,
				Cron: new("0 0 * * *"),
			},
			check: func(t *testing.T, r *model.Remind) {
				assert.Equal(t, "0 0 * * *", r.Cron)
			},
		},
		{
			name: "更新状态",
			update: &model.UpdateRemind{
				Id:     8,
				Status: new(model.RemindStatusDone),
			},
			check: func(t *testing.T, r *model.Remind) {
				assert.Equal(t, model.RemindStatusDone, r.Status)
			},
		},
		{
			name: "更新多个字段",
			update: &model.UpdateRemind{
				Id:      8,
				Content: new("多字段"),
				Cron:    new("0 0 * * *"),
				Status:  new(model.RemindStatusPending),
			},
			check: func(t *testing.T, r *model.Remind) {
				assert.Equal(t, "多字段", r.Content)
				assert.Equal(t, "0 0 * * *", r.Cron)
				assert.Equal(t, model.RemindStatusPending, r.Status)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbunit.New(t, func(d *dbunit.DBUnit) {
				db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("reminds"))
				s := New(db)
				err := s.UpdateRemind(context.Background(), tt.update)
				require.NoError(t, err)

				got, err := s.GetRemind(context.Background(), tt.update.Id)
				require.NoError(t, err)
				tt.check(t, got)
			})
		})
	}
}

func TestRemind_DeleteRemind(t *testing.T) {
	t.Run("删除存在提醒", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("reminds"))
			s := New(db)
			err := s.DeleteRemind(context.Background(), 8)
			require.NoError(t, err)

			_, err = s.GetRemind(context.Background(), 8)
			require.Error(t, err)
		})
	})

	t.Run("删除不存在提醒不报错", func(t *testing.T) {
		dbunit.New(t, func(d *dbunit.DBUnit) {
			db := d.NewDatabase(testutil.Schema(), testutil.Fixtures("reminds"))
			s := New(db)
			err := s.DeleteRemind(context.Background(), 999)
			require.NoError(t, err)
		})
	})
}
