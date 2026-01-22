package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Config is database connection configuration
type Config struct {
	Driver       string `yaml:"driver"`
	Dsn          string `yaml:"dsn"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxLifetime  int    `yaml:"max_lifetime"`
}

func (d *Config) connect() (*sql.DB, error) {
	sess, err := sql.Open("pgx", d.Dsn)

	if err != nil {
		return nil, err
	}

	if err := sess.Ping(); err != nil {
		_ = sess.Close()
		return nil, err
	}

	sess.SetMaxOpenConns(d.MaxOpenConns)
	sess.SetMaxIdleConns(d.MaxIdleConns)
	if d.MaxLifetime > 0 {
		sess.SetConnMaxLifetime(time.Duration(d.MaxLifetime) * time.Second)
	}
	return sess, nil
}

func Connect(conf Config) (*sql.DB, func()) {
	log.Printf("[db] connect: %s\n", conf.Dsn)
	db, err := conf.connect()
	if err != nil {
		log.Fatalf("[db] failed to connect: %s\n", err)
	}
	return db, func() {
		if err := db.Close(); err != nil {
			log.Printf("[db] database close error: %s", err)
		}
	}
}
