package scanner

import (
	"errors"

	"go.etcd.io/bbolt"
)

const (
	rescanRequestKey = "rescan"
)

type Repository interface {
	UpdateRescanRequest(email string, value []byte) error
	GetRescanRequest(email string) ([]byte, error)
	DeleteRescanRequest(email string) error
}

type repo struct {
	DB *bbolt.DB
}

type scanQueueQueryResult struct {
	Key   []byte
	Value []byte
}

func NewRepository(db *bbolt.DB) *repo {
	return &repo{DB: db}
}

func (r *repo) UpdateRescanRequest(email string, value []byte) error {
	return r.DB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(email))
		if err != nil {
			return err
		}

		return bucket.Put([]byte(rescanRequestKey), value)
	})
}

func (r *repo) GetRescanRequest(email string) ([]byte, error) {
	var value []byte

	err := r.DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(email))
		if bucket == nil {
			return errors.New("account not found")
		}

		value = bucket.Get([]byte(rescanRequestKey))

		return nil
	})

	return value, err
}

func (r *repo) DeleteRescanRequest(email string) error {
	return r.DB.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(email))
		if bucket == nil {
			return errors.New("account not found")
		}

		return bucket.Delete([]byte(rescanRequestKey))
	})
}
