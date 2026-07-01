package admin

import (
	"context"
	"fmt"
	"time"

	"app/pkg/remindutil"
	adminv1 "app/proto/gen/admin/v1"
	"app/proto/gen/types"
	"app/store"
	"app/store/model"

	"google.golang.org/protobuf/types/known/emptypb"
)

var _ adminv1.RemindServiceHTTPServer = (*Remind)(nil)

type Remind struct {
	store *store.Store
}

func NewRemind(s *store.Store) *Remind {
	return &Remind{store: s}
}

func (r *Remind) List(ctx context.Context, req *adminv1.RemindListRequest) (*adminv1.RemindListResponse, error) {
	num := 10
	reminds, err := r.store.ListRemind(ctx, int(req.GetPage()), num)
	if err != nil {
		return nil, err
	}
	items := make([]*adminv1.RemindItem, 0, len(reminds))
	for _, v := range reminds {
		items = append(items, adminv1.RemindItem_builder{Id: int32(v.Id),
			Cron:      v.Cron,
			Content:   v.Content,
			Status:    int32(v.Status),
			NextTime:  v.NextTime.Format(time.DateTime),
			CreatedAt: v.CreatedAt.Format(time.DateTime)}.Build(),
		)
	}
	total, err := r.store.CountRemindTotal(ctx)
	if err != nil {
		return nil, err
	}
	return adminv1.RemindListResponse_builder{List: items,
			Total: int32(total)}.Build(),
		nil
}

func (r *Remind) Create(ctx context.Context, req *adminv1.RemindCreateRequest) (*types.IDResponse, error) {
	c := &model.Remind{
		Cron:      req.GetCron(),
		Content:   req.GetContent(),
		Status:    1,
		CreatedAt: time.Now(),
	}
	c.NextTime = remindutil.NextTimeFromRule(c.CreatedAt, c)

	// 如果由于解析失败等原因获取到了零值，返回错误，不存入数据库
	if c.NextTime.IsZero() {
		return nil, fmt.Errorf("无效的 Cron 表达式或时间格式: %s", req.GetCron())
	}

	lastId, err := r.store.CreateRemind(ctx, c)
	if err != nil {
		return nil, err
	}
	return types.IDResponse_builder{Id: int32(lastId)}.Build(), nil
}

func (r *Remind) Update(ctx context.Context, req *adminv1.RemindUpdateRequest) (*types.IDResponse, error) {
	u := &model.UpdateRemind{Id: int(req.GetId())}

	if req.GetCron() != "" {
		v := req.GetCron()
		u.Cron = &v
	}
	if req.GetContent() != "" {
		v := req.GetContent()
		u.Content = &v
	}

	nextTime := remindutil.NextTimeFromRule(time.Now(), &model.Remind{
		Cron:      req.GetCron(),
		CreatedAt: time.Now(),
	})

	// 更新时如果包含了 Cron，则校验有效性
	if req.GetCron() != "" && nextTime.IsZero() {
		return nil, fmt.Errorf("无效的 Cron 表达式或时间格式: %s", req.GetCron())
	}

	u.NextTime = &nextTime

	if err := r.store.UpdateRemind(ctx, u); err != nil {
		return nil, err
	}
	return types.IDResponse_builder{Id: int32(req.GetId())}.Build(), nil
}

func (r *Remind) Delete(ctx context.Context, req *adminv1.RemindDeleteRequest) (*emptypb.Empty, error) {
	if err := r.store.DeleteRemind(ctx, int(req.GetId())); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
