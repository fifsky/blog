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

var _ adminv1.FootprintServiceServer = (*Footprint)(nil)

type Footprint struct {
	adminv1.UnimplementedFootprintServiceServer
	store *store.Store
}

func NewFootprint(s *store.Store) *Footprint {
	return &Footprint{store: s}
}

func (f *Footprint) List(ctx context.Context, req *adminv1.FootprintListRequest) (*adminv1.FootprintListResponse, error) {
	num := 10
	footprints, err := f.store.ListFootprint(ctx, int(req.Page), num)
	if err != nil {
		return nil, err
	}

	items := make([]*adminv1.FootprintItem, 0, len(footprints))
	for _, v := range footprints {
		item := &adminv1.FootprintItem{
			Id:          int32(v.Id),
			Name:        v.Name,
			Description: v.Description,
			Longitude:   v.Longitude,
			Latitude:    v.Latitude,
			Date:        v.Date,
			MarkerColor: v.MarkerColor,
			Url:         v.Url,
			UrlLabel:    v.UrlLabel,
			CreatedAt:   v.CreatedAt.Format(time.DateTime),
		}

		for _, c := range v.Categories {
			item.Categories = append(item.Categories, c)
		}

		for _, p := range v.Photos {
			item.Photos = append(item.Photos, &adminv1.FootprintPhotoItem{
				Src:       p.Src,
				Thumbnail: p.Thumbnail,
			})
		}

		items = append(items, item)
	}

	total, err := f.store.CountFootprintTotal(ctx)
	if err != nil {
		return nil, err
	}

	return &adminv1.FootprintListResponse{
		List:  items,
		Total: int32(total),
	}, nil
}

func (f *Footprint) Create(ctx context.Context, req *adminv1.FootprintCreateRequest) (*types.IDResponse, error) {
	photos := model.PhotosFromURLs(req.PhotoUrls)

	categories := make([]string, 0, len(req.Categories))
	categories = append(categories, req.Categories...)

	fp := &model.Footprint{
		Name:        req.Name,
		Description: req.Description,
		Longitude:   req.Longitude,
		Latitude:    req.Latitude,
		Date:        req.Date,
		MarkerColor: req.MarkerColor,
		Categories:  categories,
		Url:         req.Url,
		UrlLabel:    req.UrlLabel,
		Photos:      photos,
	}

	id, err := f.store.CreateFootprint(ctx, fp)
	if err != nil {
		return nil, err
	}

	return &types.IDResponse{Id: int32(id)}, nil
}

func (f *Footprint) Update(ctx context.Context, req *adminv1.FootprintUpdateRequest) (*types.IDResponse, error) {
	u := &model.UpdateFootprint{Id: int(req.Id)}

	if req.Name != "" {
		u.Name = &req.Name
	}
	if req.Description != "" {
		u.Description = &req.Description
	}
	if req.Longitude != "" {
		u.Longitude = &req.Longitude
	}
	if req.Latitude != "" {
		u.Latitude = &req.Latitude
	}
	if req.Date != "" {
		u.Date = &req.Date
	}
	if req.MarkerColor != "" {
		u.MarkerColor = &req.MarkerColor
	}
	if req.Categories != nil {
		u.Categories = req.Categories
	}
	if req.Url != "" {
		u.Url = &req.Url
	}
	if req.UrlLabel != "" {
		u.UrlLabel = &req.UrlLabel
	}
	if req.PhotoUrls != nil {
		u.PhotoUrls = req.PhotoUrls
	}

	if err := f.store.UpdateFootprint(ctx, u); err != nil {
		return nil, err
	}

	return &types.IDResponse{Id: req.Id}, nil
}

func (f *Footprint) Delete(ctx context.Context, req *adminv1.FootprintDeleteRequest) (*emptypb.Empty, error) {
	if err := f.store.DeleteFootprint(ctx, int(req.Id)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
