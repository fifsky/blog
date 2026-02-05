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

var _ adminv1.MoodServiceServer = (*Mood)(nil)

type Mood struct {
	adminv1.UnimplementedMoodServiceServer
	store *store.Store
}

func NewMood(s *store.Store) *Mood {
	return &Mood{store: s}
}

func (m *Mood) Create(ctx context.Context, req *adminv1.MoodCreateRequest) (*types.IDResponse, error) {
	loginUser := GetLoginUser(ctx)
	c := &model.Mood{
		Content:   req.Content,
		UserId:    loginUser.Id,
		CreatedAt: time.Now(),
	}
	lastId, err := m.store.CreateMood(ctx, c)
	if err != nil {
		return nil, err
	}
	return &types.IDResponse{Id: int32(lastId)}, nil
}

func (m *Mood) Update(ctx context.Context, req *adminv1.MoodUpdateRequest) (*types.IDResponse, error) {
	u := &model.UpdateMood{Id: int(req.Id)}
	if req.Content != "" {
		v := req.Content
		u.Content = &v
	}
	if err := m.store.UpdateMood(ctx, u); err != nil {
		return nil, err
	}
	return &types.IDResponse{Id: int32(req.Id)}, nil
}

func (m *Mood) Delete(ctx context.Context, req *adminv1.MoodDeleteRequest) (*emptypb.Empty, error) {
	ids := make([]int, len(req.Ids))
	for i, id := range req.Ids {
		ids[i] = int(id)
	}
	if err := m.store.DeleteMood(ctx, ids); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
