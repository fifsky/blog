package model

import (
	"encoding/json"
	"time"
)

type Footprint struct {
	Id          int
	Name        string
	Description string
	Longitude   string
	Latitude    string
	Date        string
	MarkerColor string
	Categories  []string
	Url         string
	UrlLabel    string
	Photos      []FootprintPhoto
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type FootprintPhoto struct {
	Src       string
	Thumbnail string
}

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
