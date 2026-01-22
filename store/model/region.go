package model

type Region struct {
	RegionId   int
	ParentId   int
	Level      int
	RegionName string
	Longitude  string
	Latitude   string
	Pinyin     string
	AzNo       string
}
