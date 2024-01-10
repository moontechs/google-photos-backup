package db

import (
	"fmt"
	"os"

	"go.etcd.io/bbolt"
)

type Connection struct {
	DB *bbolt.DB
}

func NewConnection() (*Connection, error) {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "/Users/michael/github.com/moontechs/photos-backup/database.db" // TODO fix
	}

	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	return &Connection{DB: db}, nil
}

func (c *Connection) Close() error {
	if c.DB == nil {
		return nil
	}

	return c.DB.Close()
}
