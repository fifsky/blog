package openapi

import (
	"database/sql"
	"net/http"
	"strconv"

	"app/pkg/errors"
	"app/server/response"
	"app/store"
)

type Geo struct {
	store *store.Store
}

type NearestRegionResponse struct {
	ProvinceID   int    `json:"province_id"`
	ProvinceName string `json:"province_name"`
	CityID       int    `json:"city_id"`
	CityName     string `json:"city_name"`
}

func NewGeo(s *store.Store) *Geo {
	return &Geo{store: s}
}

func (g *Geo) Nearest(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	lat, err := strconv.ParseFloat(q.Get("latitude"), 64)
	if err != nil {
		response.Fail(w, errors.BadRequest("INVALID_LATITUDE", "latitude 参数错误").WithCause(err))
		return
	}
	lng, err := strconv.ParseFloat(q.Get("longitude"), 64)
	if err != nil {
		response.Fail(w, errors.BadRequest("INVALID_LONGITUDE", "longitude 参数错误").WithCause(err))
		return
	}

	city, province, err := g.store.FindNearestCity(r.Context(), lat, lng)
	if err != nil {
		if err == sql.ErrNoRows {
			response.Fail(w, errors.NotFound("REGION_NOT_FOUND", "未找到匹配的城市"))
			return
		}
		response.Fail(w, errors.ErrSystem.WithCause(err))
		return
	}

	response.Success(w, &NearestRegionResponse{
		ProvinceID:   province.RegionId,
		ProvinceName: province.RegionName,
		CityID:       city.RegionId,
		CityName:     city.RegionName,
	})
}
