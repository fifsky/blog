package store

import (
	"database/sql"
	"sync"
)

type Store struct {
	db           *sql.DB
	optionsCache map[string]string
	optionsMu    sync.RWMutex
}

func New(db *sql.DB) *Store {
	return &Store{db: db}
}
