package model

// Option 配置项模型
type Option struct {
	Id          int    // PK
	OptionKey   string // 配置项唯一key
	OptionValue string // 配置内容
}
