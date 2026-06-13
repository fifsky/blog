package openapi

import (
	"context"
	"time"

	apiv1 "app/proto/gen/api/v1"
	"app/proto/gen/types"
	"app/store"
	"app/store/model"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ apiv1.MoodServiceHTTPServer = (*Mood)(nil)

type Mood struct {
	store *store.Store
}

func NewMood(s *store.Store) *Mood {
	return &Mood{store: s}
}

func (m *Mood) List(ctx context.Context, req *apiv1.MoodListRequest) (*apiv1.MoodListResponse, error) {
	num := 10
	moods, err := m.store.ListMood(ctx, int(req.GetPage()), num)
	if err != nil {
		return nil, err
	}
	uids := lo.Map(moods, func(item model.Mood, index int) int { return item.UserId })
	um, _ := m.store.GetUserByIds(ctx, uids)

	items := make([]*apiv1.MoodItem, 0, len(moods))
	for _, md := range moods {
		item := apiv1.MoodItem_builder{Id: int32(md.Id),
			Content:   md.Content,
			CreatedAt: md.CreatedAt.Format(time.DateTime)}.Build()

		if u, ok := um[md.UserId]; ok {
			item.SetUser(types.UserSummary_builder{Id: int32(u.Id), Name: u.Name, NickName: u.NickName}.Build())
		}
		items = append(items, item)
	}
	total, err := m.store.CountMoodTotal(ctx)
	if err != nil {
		return nil, err
	}
	return apiv1.MoodListResponse_builder{List: items,
			Total: int32(total)}.Build(),
		nil
}

func (m *Mood) Random(ctx context.Context, _ *emptypb.Empty) (*apiv1.MoodItem, error) {
	md, err := m.store.RandomMood(ctx)
	if err != nil {
		return nil, err
	}
	um, _ := m.store.GetUserByIds(ctx, []int{md.UserId})
	item := apiv1.MoodItem_builder{Id: int32(md.Id),
		Content:   md.Content,
		CreatedAt: md.CreatedAt.Format(time.DateTime)}.Build()

	if u, ok := um[md.UserId]; ok {
		item.SetUser(types.UserSummary_builder{Id: int32(u.Id), Name: u.Name, NickName: u.NickName}.Build())
	}
	return item, nil
}
