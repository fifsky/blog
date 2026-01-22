package openapi

import (
	"context"
	"strconv"

	apiv1 "app/proto/gen/api/v1"
	"app/store"
)

var _ apiv1.TravelServiceServer = (*Travel)(nil)

type Travel struct {
	apiv1.UnimplementedTravelServiceServer
	store *store.Store
}

func NewTravel(s *store.Store) *Travel {
	return &Travel{store: s}
}

func (t *Travel) GetFootprints(ctx context.Context, req *apiv1.GetFootprintsRequest) (*apiv1.GetFootprintsResponse, error) {
	provinces, err := t.store.ListProvincesWithPhotos(ctx)
	if err != nil {
		return nil, err
	}

	cities, err := t.store.ListCitiesWithPhotos(ctx)
	if err != nil {
		return nil, err
	}

	provinceItems := make([]*apiv1.FootprintRegion, 0, len(provinces))
	for _, v := range provinces {
		provinceItems = append(provinceItems, &apiv1.FootprintRegion{
			RegionId:  strconv.Itoa(v.RegionId),
			Name:      v.RegionName,
			Longitude: v.Longitude,
			Latitude:  v.Latitude,
		})
	}

	cityItems := make([]*apiv1.FootprintRegion, 0, len(cities))
	for _, v := range cities {
		cityItems = append(cityItems, &apiv1.FootprintRegion{
			RegionId:  strconv.Itoa(v.RegionId),
			Name:      v.RegionName,
			Longitude: v.Longitude,
			Latitude:  v.Latitude,
		})
	}

	return &apiv1.GetFootprintsResponse{
		Provinces: provinceItems,
		Cities:    cityItems,
	}, nil
}

func (t *Travel) ListCityPhotos(ctx context.Context, req *apiv1.ListCityPhotosRequest) (*apiv1.ListCityPhotosResponse, error) {
	photos, err := t.store.ListPhotoByCity(ctx, req.RegionId)
	if err != nil {
		return nil, err
	}

	items := make([]*apiv1.TravelPhoto, 0, len(photos))
	for _, v := range photos {
		items = append(items, &apiv1.TravelPhoto{
			Title:       v.Title,
			Description: v.Description,
			Src:         v.Src,
			Thumbnail:   v.Thumbnail,
		})
	}

	return &apiv1.ListCityPhotosResponse{
		Photos: items,
	}, nil
}
