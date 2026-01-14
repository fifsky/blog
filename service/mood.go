package service

import (
	"context"
	"time"

	"app/model"
	apiv1 "app/proto/gen/api/v1"
	"app/store"

	"github.com/samber/lo"
)

var _ apiv1.MoodServiceServer = (*Mood)(nil)

type Mood struct {
	apiv1.UnimplementedMoodServiceServer
	store *store.Store
}

func NewMood(s *store.Store) *Mood {
	return &Mood{store: s}
}

func (m *Mood) List(ctx context.Context, req *apiv1.PageRequest) (*apiv1.MoodListResponse, error) {
	num := 10
	moods, err := m.store.ListMood(ctx, int(req.Page), num)
	if err != nil {
		return nil, err
	}
	uids := lo.Map(moods, func(item model.Mood, index int) int { return item.UserId })
	um, _ := m.store.GetUserByIds(ctx, uids)

	items := make([]*apiv1.MoodItem, 0, len(moods))
	for _, md := range moods {
		item := &apiv1.MoodItem{
			Id:        int32(md.Id),
			Content:   md.Content,
			CreatedAt: md.CreatedAt.Format(time.DateTime),
		}
		if u, ok := um[md.UserId]; ok {
			item.User = &apiv1.UserSummary{Id: int32(u.Id), Name: u.Name, NickName: u.NickName}
		}
		items = append(items, item)
	}
	total, err := m.store.CountMoodTotal(ctx)
	if err != nil {
		return nil, err
	}
	return &apiv1.MoodListResponse{
		List:  items,
		Total: int32(total),
	}, nil
}
