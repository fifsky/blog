package admin

import (
	"context"
	"time"

	adminv1 "app/proto/gen/admin/v1"
	"app/proto/gen/types"
	"app/store"
	"app/store/model"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ adminv1.FootprintServiceHTTPServer = (*Footprint)(nil)

type Footprint struct {
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

	items := lo.Map(footprints, func(v *model.Footprint, _ int) *adminv1.FootprintItem {
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

		photos := lo.Map(v.Photos, func(p model.FootprintPhoto, _ int) *adminv1.FootprintPhotoItem {
			return adminv1.FootprintPhotoItem_builder{Src: p.Src, Thumbnail: p.Thumbnail}.Build()
		})
		item.SetPhotos(photos)

		return item
	})

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
		u.Name = new(req.GetName())
	}
	if req.GetDescription() != "" {
		u.Description = new(req.GetDescription())
	}
	if req.GetLongitude() != "" {
		u.Longitude = new(req.GetLongitude())
	}
	if req.GetLatitude() != "" {
		u.Latitude = new(req.GetLatitude())
	}
	if req.GetDate() != "" {
		u.Date = new(req.GetDate())
	}
	if req.GetMarkerColor() != "" {
		u.MarkerColor = new(req.GetMarkerColor())
	}
	if req.GetCategories() != nil {
		u.Categories = req.GetCategories()
	}
	if req.GetUrl() != "" {
		u.Url = new(req.GetUrl())
	}
	if req.GetUrlLabel() != "" {
		u.UrlLabel = new(req.GetUrlLabel())
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
