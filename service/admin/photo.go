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

var _ adminv1.PhotoServiceServer = (*Photo)(nil)

type Photo struct {
	adminv1.UnimplementedPhotoServiceServer
	store *store.Store
}

func NewPhoto(s *store.Store) *Photo {
	return &Photo{store: s}
}

func (p *Photo) List(ctx context.Context, req *adminv1.PhotoListRequest) (*adminv1.PhotoListResponse, error) {
	num := 10
	photos, err := p.store.ListPhoto(ctx, int(req.Page), num)
	if err != nil {
		return nil, err
	}

	regionIDs := make([]int, 0, len(photos)*2)
	regionSeen := make(map[int]struct{}, len(photos)*2)
	for _, v := range photos {
		if v.Province > 0 {
			if _, ok := regionSeen[v.Province]; !ok {
				regionSeen[v.Province] = struct{}{}
				regionIDs = append(regionIDs, v.Province)
			}
		}
		if v.City > 0 {
			if _, ok := regionSeen[v.City]; !ok {
				regionSeen[v.City] = struct{}{}
				regionIDs = append(regionIDs, v.City)
			}
		}
	}
	regionMap, err := p.store.GetRegionByIds(ctx, regionIDs)
	if err != nil {
		return nil, err
	}

	items := make([]*adminv1.PhotoItem, 0, len(photos))
	for _, v := range photos {
		item := &adminv1.PhotoItem{
			Id:          int32(v.Id),
			Title:       v.Title,
			Description: v.Description,
			Src:         v.Src,
			Thumbnail:   v.Thumbnail,
			Province:    int32(v.Province),
			City:        int32(v.City),
			CreatedAt:   v.CreatedAt.Format(time.DateTime),
		}

		if v.Province > 0 && regionMap != nil {
			if province, ok := regionMap[v.Province]; ok {
				item.ProvinceName = province.RegionName
			}
		}
		if v.City > 0 && regionMap != nil {
			if city, ok := regionMap[v.City]; ok {
				item.CityName = city.RegionName
			}
		}

		items = append(items, item)
	}

	total, err := p.store.CountPhotoTotal(ctx)
	if err != nil {
		return nil, err
	}

	return &adminv1.PhotoListResponse{
		List:  items,
		Total: int32(total),
	}, nil
}

func (p *Photo) Create(ctx context.Context, req *adminv1.PhotoCreateRequest) (*types.IDResponse, error) {
	var lastId int64
	var err error

	// Create a record for each src
	for _, src := range req.Srcs {
		// Build thumbnail URL by appending !photothumb suffix
		thumbnail := src + "!photothumb"

		photo := &model.Photo{
			Title:       req.Title,
			Description: req.Description,
			Src:         src,
			Thumbnail:   thumbnail,
			Province:    int(req.Province),
			City:        int(req.City),
		}

		lastId, err = p.store.CreatePhoto(ctx, photo)
		if err != nil {
			return nil, err
		}
	}

	return &types.IDResponse{Id: int32(lastId)}, nil
}

func (p *Photo) Update(ctx context.Context, req *adminv1.PhotoUpdateRequest) (*types.IDResponse, error) {
	u := &model.UpdatePhoto{Id: int(req.Id)}

	if req.Title != "" {
		u.Title = &req.Title
	}
	if req.Description != "" {
		u.Description = &req.Description
	}
	if req.Province > 0 {
		v := int(req.Province)
		u.Province = &v
	}
	if req.City > 0 {
		v := int(req.City)
		u.City = &v
	}

	if err := p.store.UpdatePhoto(ctx, u); err != nil {
		return nil, err
	}

	return &types.IDResponse{Id: req.Id}, nil
}

func (p *Photo) Delete(ctx context.Context, req *adminv1.PhotoDeleteRequest) (*emptypb.Empty, error) {
	if err := p.store.DeletePhoto(ctx, int(req.Id)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
