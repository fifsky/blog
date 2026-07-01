package model

import (
	"encoding/json"
	"time"
)

// Footprint 旅行足迹模型
type Footprint struct {
	Id          int              // PK
	Name        string           // 地点名称
	Description string           // 描述
	Longitude   string           // 经度
	Latitude    string           // 纬度
	Date        string           // 到访日期
	MarkerColor string           // 标记颜色
	Categories  []string         // 分类标签数组
	Url         string           // 关联链接
	UrlLabel    string           // 链接按钮文案
	Photos      []FootprintPhoto // 照片列表
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// FootprintPhoto 足迹照片
type FootprintPhoto struct {
	Src       string // 原图地址
	Thumbnail string // 缩略图地址
}

// UpdateFootprint 更新足迹参数
type UpdateFootprint struct {
	Id          int
	Name        *string
	Description *string
	Longitude   *string
	Latitude    *string
	Date        *string
	MarkerColor *string
	Categories  []string
	Url         *string
	UrlLabel    *string
	PhotoUrls   []string
}

// FootprintJSON 足迹JSON序列化模型
type FootprintJSON struct {
	Id          int              `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Longitude   string           `json:"longitude"`
	Latitude    string           `json:"latitude"`
	Date        string           `json:"date"`
	MarkerColor string           `json:"marker_color"`
	Categories  []string         `json:"categories"`
	Url         string           `json:"url"`
	UrlLabel    string           `json:"url_label"`
	Photos      []FootprintPhoto `json:"photos"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

func (f *Footprint) ScanCategories(raw json.RawMessage) error {
	if raw == nil {
		f.Categories = nil
		return nil
	}
	return json.Unmarshal(raw, &f.Categories)
}

func (f *Footprint) ScanPhotos(raw json.RawMessage) error {
	if raw == nil {
		f.Photos = nil
		return nil
	}
	return json.Unmarshal(raw, &f.Photos)
}

func PhotosFromURLs(urls []string) []FootprintPhoto {
	photos := make([]FootprintPhoto, 0, len(urls))
	for _, u := range urls {
		photos = append(photos, FootprintPhoto{
			Src:       u,
			Thumbnail: u + "!photothumb",
		})
	}
	return photos
}
