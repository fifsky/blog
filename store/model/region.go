package model

// Region 中国行政区域模型
type Region struct {
	RegionId   int    // 区域ID
	ParentId   int    // 上级ID
	Level      int    // 层级
	RegionName string // 区域名称
	Longitude  string // 经度
	Latitude   string // 纬度
	Pinyin     string // 拼音
	AzNo       string // 首字母
}
