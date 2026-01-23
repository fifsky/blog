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

	items := make([]*adminv1.PhotoItem, 0, len(photos))
	for _, v := range photos {
		item := &adminv1.PhotoItem{
			Id:          int32(v.Id),
			Title:       v.Title,
			Description: v.Description,
			Src:         v.Src,
			Thumbnail:   v.Thumbnail,
			Province:    v.Province,
			City:        v.City,
			CreatedAt:   v.CreatedAt.Format(time.DateTime),
		}

		// Get province name
		if v.Province != "" {
			if province, err := p.store.GetRegion(ctx, parseInt(v.Province)); err == nil {
				item.ProvinceName = province.RegionName
			}
		}

		// Get city name
		if v.City != "" {
			if city, err := p.store.GetRegion(ctx, parseInt(v.City)); err == nil {
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
			Province:    req.Province,
			City:        req.City,
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
	if req.Province != "" {
		u.Province = &req.Province
	}
	if req.City != "" {
		u.City = &req.City
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

// parseInt converts string to int, returns 0 if failed
func parseInt(s string) int {
	var n int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}
