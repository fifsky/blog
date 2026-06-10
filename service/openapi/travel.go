package openapi

import (
	"context"

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
	footprints, err := t.store.ListAllFootprints(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]*apiv1.Footprint, 0, len(footprints))
	for _, v := range footprints {
		item := apiv1.Footprint_builder{Id: int32(v.Id),
			Name:        v.Name,
			Description: v.Description,
			Longitude:   v.Longitude,
			Latitude:    v.Latitude,
			Date:        v.Date,
			MarkerColor: v.MarkerColor,
			Url:         v.Url,
			UrlLabel:    v.UrlLabel}.Build()

		item.SetCategories(append(item.GetCategories(), v.Categories...))

		for _, p := range v.Photos {
			item.SetPhotos(append(item.GetPhotos(), apiv1.FootprintPhoto_builder{Src: p.Src,
				Thumbnail: p.Thumbnail}.Build(),
			))
		}

		items = append(items, item)
	}

	return apiv1.GetFootprintsResponse_builder{Footprints: items}.Build(),
		nil
}
