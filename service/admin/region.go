package admin

import (
	"context"

	adminv1 "app/proto/gen/admin/v1"
	"app/store"
)

var _ adminv1.RegionServiceServer = (*Region)(nil)

type Region struct {
	adminv1.UnimplementedRegionServiceServer
	store *store.Store
}

func NewRegion(s *store.Store) *Region {
	return &Region{store: s}
}

func (r *Region) ListByParent(ctx context.Context, req *adminv1.RegionListRequest) (*adminv1.RegionListResponse, error) {
	regions, err := r.store.ListRegionByParent(ctx, int(req.ParentId))
	if err != nil {
		return nil, err
	}

	items := make([]*adminv1.RegionItem, 0, len(regions))
	for _, v := range regions {
		items = append(items, &adminv1.RegionItem{
			RegionId:   int32(v.RegionId),
			ParentId:   int32(v.ParentId),
			Level:      int32(v.Level),
			RegionName: v.RegionName,
			Longitude:  v.Longitude,
			Latitude:   v.Latitude,
			Pinyin:     v.Pinyin,
			AzNo:       v.AzNo,
		})
	}

	return &adminv1.RegionListResponse{
		List: items,
	}, nil
}
