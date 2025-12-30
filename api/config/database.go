package config

import (
	"database/sql"

	"app/connect/db"
)

type DBConf struct {
	Blog db.Config `yaml:"blog"`
}

func NewBlogDB(conf *Config) (*sql.DB, func()) {
	return db.Connect(conf.DB.Blog)
}
