package admin

import (
	"context"
	"time"

	adminv1 "app/proto/gen/admin/v1"
	"app/proto/gen/types"
	"app/service/remind"
	"app/store"
	"app/store/model"

	"google.golang.org/protobuf/types/known/emptypb"
)

var _ adminv1.RemindServiceServer = (*Remind)(nil)

type Remind struct {
	adminv1.UnimplementedRemindServiceServer
	store *store.Store
}

func NewRemind(s *store.Store) *Remind {
	return &Remind{store: s}
}

func (r *Remind) List(ctx context.Context, req *adminv1.RemindListRequest) (*adminv1.RemindListResponse, error) {
	num := 10
	reminds, err := r.store.ListRemind(ctx, int(req.Page), num)
	if err != nil {
		return nil, err
	}
	items := make([]*adminv1.RemindItem, 0, len(reminds))
	for _, v := range reminds {
		items = append(items, &adminv1.RemindItem{
			Id:        int32(v.Id),
			Cron:      v.Cron,
			Content:   v.Content,
			Status:    int32(v.Status),
			NextTime:  v.NextTime.Format(time.DateTime),
			CreatedAt: v.CreatedAt.Format(time.DateTime),
		})
	}
	total, err := r.store.CountRemindTotal(ctx)
	if err != nil {
		return nil, err
	}
	return &adminv1.RemindListResponse{
		List:  items,
		Total: int32(total),
	}, nil
}

func (r *Remind) Create(ctx context.Context, req *adminv1.RemindCreateRequest) (*types.IDResponse, error) {
	c := &model.Remind{
		Cron:      req.Cron,
		Content:   req.Content,
		Status:    1,
		CreatedAt: time.Now(),
	}
	c.NextTime = remind.NextTimeFromRule(c.CreatedAt, c)
	lastId, err := r.store.CreateRemind(ctx, c)
	if err != nil {
		return nil, err
	}
	return &types.IDResponse{Id: int32(lastId)}, nil
}

func (r *Remind) Update(ctx context.Context, req *adminv1.RemindUpdateRequest) (*types.IDResponse, error) {
	u := &model.UpdateRemind{Id: int(req.Id)}

	if req.Cron != "" {
		v := req.Cron
		u.Cron = &v
	}
	if req.Content != "" {
		v := req.Content
		u.Content = &v
	}

	nextTime := remind.NextTimeFromRule(time.Now(), &model.Remind{
		Cron:      req.Cron,
		CreatedAt: time.Now(),
	})
	u.NextTime = &nextTime

	if err := r.store.UpdateRemind(ctx, u); err != nil {
		return nil, err
	}
	return &types.IDResponse{Id: int32(req.Id)}, nil
}

func (r *Remind) Delete(ctx context.Context, req *adminv1.RemindDeleteRequest) (*emptypb.Empty, error) {
	if err := r.store.DeleteRemind(ctx, int(req.Id)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
