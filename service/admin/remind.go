package admin

import (
	"context"
	"time"

	adminv1 "app/proto/gen/admin/v1"
	"app/proto/gen/types"
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
			Type:      int32(v.Type),
			Content:   v.Content,
			Month:     int32(v.Month),
			Week:      int32(v.Week),
			Day:       int32(v.Day),
			Hour:      int32(v.Hour),
			Minute:    int32(v.Minute),
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
		Type:      int(req.Type),
		Content:   req.Content,
		Month:     int(req.Month),
		Week:      int(req.Week),
		Day:       int(req.Day),
		Hour:      int(req.Hour),
		Minute:    int(req.Minute),
		Status:    1,
		CreatedAt: time.Now(),
	}
	lastId, err := r.store.CreateRemind(ctx, c)
	if err != nil {
		return nil, err
	}
	return &types.IDResponse{Id: int32(lastId)}, nil
}

func (r *Remind) Update(ctx context.Context, req *adminv1.RemindUpdateRequest) (*types.IDResponse, error) {
	u := &model.UpdateRemind{Id: int(req.Id)}
	if req.Type > 0 {
		v := int(req.Type)
		u.Type = &v
	}
	if req.Content != "" {
		v := req.Content
		u.Content = &v
	}
	if req.Month > 0 {
		v := int(req.Month)
		u.Month = &v
	}
	if req.Week > 0 {
		v := int(req.Week)
		u.Week = &v
	}
	if req.Day > 0 {
		v := int(req.Day)
		u.Day = &v
	}
	if req.Hour > 0 {
		v := int(req.Hour)
		u.Hour = &v
	}
	if req.Minute > 0 {
		v := int(req.Minute)
		u.Minute = &v
	}
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
