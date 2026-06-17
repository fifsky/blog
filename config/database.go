package config

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
)

// Config is database connection configuration
type Database struct {
	Driver       string `yaml:"driver"`
	Dsn          string `yaml:"dsn"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxLifetime  int    `yaml:"max_lifetime"`
}

func (d *Database) Connect() *sql.DB {
	dsn, err := mysql.ParseDSN(d.Dsn)
	if err != nil {
		log.Fatalf("[db] failed parse dsn: %s\n", err)
	}

	log.Printf("[db] connect: %s %s\n", dsn.DBName, dsn.Addr)
	var sess *sql.DB
	sess, err = sql.Open(d.Driver, d.Dsn)
	if err != nil {
		log.Fatalf("[db] failed to connect: %s %s\n", dsn.DBName, err)
	}

	if err = sess.Ping(); err != nil {
		_ = sess.Close()
		log.Fatalf("[db] failed to connect: %s %s\n", dsn.DBName, err)
	}

	sess.SetMaxOpenConns(d.MaxOpenConns)
	sess.SetMaxIdleConns(d.MaxIdleConns)
	if d.MaxLifetime > 0 {
		sess.SetConnMaxLifetime(time.Duration(d.MaxLifetime) * time.Second)
	}
	return sess
}
