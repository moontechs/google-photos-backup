package scanner

import (
	"go.etcd.io/bbolt"
)

const (
	rescanRequestsBucketName = "rescan_requests"
	rescanRequestKey         = "rescan"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Repository
type Repository interface {
	UpdateRescanRequest(rescanType, email string, value []byte) error
	GetRescanRequest(rescanType, email string) ([]byte, error)
	DeleteRescanRequest(rescanType, email string) error
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
		bucket, err := tx.CreateBucketIfNotExists([]byte(rescanRequestsBucketName + ":" + rescanType))
		if err != nil {
			return err
		}

		return bucket.Put([]byte(email), value)
	})
}

func (r *repo) GetRescanRequest(rescanType, email string) ([]byte, error) {
	value := make([]byte, 0)

	err := r.DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(rescanRequestsBucketName + ":" + rescanType))
		if bucket == nil {
			return nil
		}

		value = bucket.Get([]byte(email))

		return nil
	})

	return value, err
}

func (r *repo) DeleteRescanRequest(rescanType, email string) error {
	return r.DB.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(rescanRequestsBucketName + ":" + rescanType))
		if bucket == nil {
			return nil
		}

		return bucket.Delete([]byte(email))
	})
}
