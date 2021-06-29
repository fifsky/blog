package model

type Options struct {
	Id          int    `json:"id" db:"id"`
	OptionKey   string `json:"option_key" db:"option_key"`
	OptionValue string `json:"option_value" db:"option_value"`
}

func (o *Options) TableName() string {
	return "options"
}

func (o *Options) PK() string {
	return "id"
}
