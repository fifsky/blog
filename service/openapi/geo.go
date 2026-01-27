package openapi

import (
	"context"
	"database/sql"

	"app/pkg/errors"
	apiv1 "app/proto/gen/api/v1"
	"app/store"
)

var _ apiv1.GeoServiceServer = (*Geo)(nil)

type Geo struct {
	apiv1.UnimplementedGeoServiceServer
	store *store.Store
}

func NewGeo(s *store.Store) *Geo {
	return &Geo{store: s}
}

func (g *Geo) GetNearestRegion(ctx context.Context, req *apiv1.GetNearestRegionRequest) (*apiv1.GetNearestRegionResponse, error) {
	city, province, err := g.store.FindNearestCity(ctx, req.Latitude, req.Longitude)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFound("REGION_NOT_FOUND", "未找到匹配的城市")
		}
		return nil, errors.ErrSystem.WithCause(err)
	}

	return &apiv1.GetNearestRegionResponse{
		ProvinceId:   int32(province.RegionId),
		ProvinceName: province.RegionName,
		CityId:       int32(city.RegionId),
		CityName:     city.RegionName,
	}, nil
}
