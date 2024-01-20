package scanner

import (
	"errors"

	"go.etcd.io/bbolt"
)

const (
	rescanRequestKey = "rescan"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Repository
type Repository interface {
	UpdateRescanRequest(rescanType, email string, value []byte) error
	GetRescanRequests(email string) (map[string][]byte, error)
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

func (r *repo) UpdateRescanRequest(rescanType, email string, value []byte) error {
	return r.DB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(email))
		if err != nil {
			return err
		}

		return bucket.Put([]byte(rescanRequestKey+"-"+rescanType), value)
	})
}

func (r *repo) GetRescanRequests(email string) (map[string][]byte, error) {
	var values map[string][]byte

	err := r.DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(email))
		if bucket == nil {
			return errors.New("account not found")
		}

		values[RescanTypePhotos] = bucket.Get([]byte(rescanRequestKey + "-" + RescanTypePhotos))
		values[RescanTypeDrive] = bucket.Get([]byte(rescanRequestKey + "-" + RescanTypeDrive))

		return nil
	})

	return values, err
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
