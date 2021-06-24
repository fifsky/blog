package model

type Options struct {
	Id          int    `form:"id" json:"id" db:"id"`
	OptionKey   string `form:"option_key" json:"option_key" db:"option_key"`
	OptionValue string `form:"option_value" json:"option_value" db:"option_value"`
}

func (o *Options) TableName() string {
	return "options"
}

func (o *Options) PK() string {
	return "id"
}
