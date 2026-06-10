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
	footprints, err := f.store.ListFootprint(ctx, int(req.GetPage()), num)
	if err != nil {
		return nil, err
	}

	items := make([]*adminv1.FootprintItem, 0, len(footprints))
	for _, v := range footprints {
		item := adminv1.FootprintItem_builder{Id: int32(v.Id),
			Name:        v.Name,
			Description: v.Description,
			Longitude:   v.Longitude,
			Latitude:    v.Latitude,
			Date:        v.Date,
			MarkerColor: v.MarkerColor,
			Url:         v.Url,
			UrlLabel:    v.UrlLabel,
			CreatedAt:   v.CreatedAt.Format(time.DateTime)}.Build()

		item.SetCategories(append(item.GetCategories(), v.Categories...))

		for _, p := range v.Photos {
			item.SetPhotos(append(item.GetPhotos(), adminv1.FootprintPhotoItem_builder{Src: p.Src,
				Thumbnail: p.Thumbnail}.Build(),
			))
		}

		items = append(items, item)
	}

	total, err := f.store.CountFootprintTotal(ctx)
	if err != nil {
		return nil, err
	}

	return adminv1.FootprintListResponse_builder{List: items,
			Total: int32(total)}.Build(),
		nil
}

func (f *Footprint) Create(ctx context.Context, req *adminv1.FootprintCreateRequest) (*types.IDResponse, error) {
	photos := model.PhotosFromURLs(req.GetPhotoUrls())

	categories := make([]string, 0, len(req.GetCategories()))
	categories = append(categories, req.GetCategories()...)

	fp := &model.Footprint{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Longitude:   req.GetLongitude(),
		Latitude:    req.GetLatitude(),
		Date:        req.GetDate(),
		MarkerColor: req.GetMarkerColor(),
		Categories:  categories,
		Url:         req.GetUrl(),
		UrlLabel:    req.GetUrlLabel(),
		Photos:      photos,
	}

	id, err := f.store.CreateFootprint(ctx, fp)
	if err != nil {
		return nil, err
	}

	return types.IDResponse_builder{Id: int32(id)}.Build(), nil
}

func (f *Footprint) Update(ctx context.Context, req *adminv1.FootprintUpdateRequest) (*types.IDResponse, error) {
	u := &model.UpdateFootprint{Id: int(req.GetId())}

	if req.GetName() != "" {
		v := req.GetName()
		u.Name = &v
	}
	if req.GetDescription() != "" {
		v := req.GetDescription()
		u.Description = &v
	}
	if req.GetLongitude() != "" {
		v := req.GetLongitude()
		u.Longitude = &v
	}
	if req.GetLatitude() != "" {
		v := req.GetLatitude()
		u.Latitude = &v
	}
	if req.GetDate() != "" {
		v := req.GetDate()
		u.Date = &v
	}
	if req.GetMarkerColor() != "" {
		v := req.GetMarkerColor()
		u.MarkerColor = &v
	}
	if req.GetCategories() != nil {
		u.Categories = req.GetCategories()
	}
	if req.GetUrl() != "" {
		v := req.GetUrl()
		u.Url = &v
	}
	if req.GetUrlLabel() != "" {
		v := req.GetUrlLabel()
		u.UrlLabel = &v
	}
	if req.GetPhotoUrls() != nil {
		u.PhotoUrls = req.GetPhotoUrls()
	}

	if err := f.store.UpdateFootprint(ctx, u); err != nil {
		return nil, err
	}

	return types.IDResponse_builder{Id: req.GetId()}.Build(), nil
}

func (f *Footprint) Delete(ctx context.Context, req *adminv1.FootprintDeleteRequest) (*emptypb.Empty, error) {
	if err := f.store.DeleteFootprint(ctx, int(req.GetId())); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
