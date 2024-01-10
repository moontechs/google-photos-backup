package config

import (
	"fmt"

	"go.etcd.io/bbolt"
)

const (
	appBucketName = "app"
	configKey     = "config"
)

type Repository interface {
	Get() ([]byte, error)
	Set(data []byte) error
}

type repository struct {
	DB *bbolt.DB
}

func NewRepository(db *bbolt.DB) repository {
	return repository{DB: db}
}

func (r repository) Get() ([]byte, error) {
	var value []byte

	err := r.DB.View(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(appBucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %w", err)
		}

		value = bucket.Get([]byte(configKey))

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("get config: %w", err)
	}

	return value, nil
}

func (r repository) Set(data []byte) error {
	return r.DB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(appBucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %w", err)
		}

		err = bucket.Put([]byte(configKey), data)
		if err != nil {
			return fmt.Errorf("put config: %w", err)
		}

		return nil
	})
}
