// Package ipgeo 提供IP地理位置查询功能，使用 api.ip.sb 服务
package ipgeo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GeoInfo IP 地理位置信息
type GeoInfo struct {
	Organization  string  `json:"organization"`
	Region        string  `json:"region"`
	ISP           string  `json:"isp"`
	RegionCode    string  `json:"region_code"`
	City          string  `json:"city"`
	ASN           int     `json:"asn"`
	Offset        int     `json:"offset"`
	Latitude      float64 `json:"latitude"`
	IP            string  `json:"ip"`
	ContinentCode string  `json:"continent_code"`
	Timezone      string  `json:"timezone"`
	Country       string  `json:"country"`
	Longitude     float64 `json:"longitude"`
	CountryCode   string  `json:"country_code"`
}

// Lookup 查询 IP 地址的地理位置信息
func Lookup(ctx context.Context, httpClient *http.Client, ip string) (*GeoInfo, error) {
	url := fmt.Sprintf("https://api.ip.sb/geoip/%s", ip)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求IP地理位置失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IP地理位置查询返回异常状态码: %d", resp.StatusCode)
	}

	var info GeoInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("解析IP地理位置响应失败: %w", err)
	}

	return &info, nil
}
