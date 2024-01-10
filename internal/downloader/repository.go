package downloader

import (
	"fmt"

	"go.etcd.io/bbolt"
)

const downloadRequestBucketName = "download_request"

type Repository interface {
	UpdateDownloadRequest(email string, mediaItemId string, value []byte) error
	GetDownloadRequest(email string) ([]byte, error)
	DeleteDownloadRequest(email string, mediaItemId string) error
}

type DownloadRequest struct {
	MediaItemId string `json:"media_item_id"`
}

type repo struct {
	DB *bbolt.DB
}

func NewRepository(db *bbolt.DB) repo {
	return repo{DB: db}
}

func (r repo) UpdateDownloadRequest(email string, mediaItemId string, value []byte) error {
	return r.DB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(email))
		if err != nil {
			return fmt.Errorf("create bucket: %w", err)
		}
		downloadRequestBucket, err := bucket.CreateBucketIfNotExists([]byte(downloadRequestBucketName))
		if err != nil {
			return fmt.Errorf("create download request bucket: %w", err)
		}

		return downloadRequestBucket.Put([]byte(mediaItemId), value)
	})
}

func (r repo) GetDownloadRequest(email string) ([]byte, error) {
	var value []byte

	err := r.DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(email))
		if bucket == nil {
			return fmt.Errorf("account not found")
		}
		downloadRequestBucket := bucket.Bucket([]byte(downloadRequestBucketName))
		if downloadRequestBucket == nil {
			return fmt.Errorf("download request bucket not found")
		}

		c := downloadRequestBucket.Cursor()

		_, v := c.First()
		value = v

		return nil
	})

	return value, err
}

func (r repo) DeleteDownloadRequest(email string, mediaItemId string) error {
	return r.DB.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(email))
		if bucket == nil {
			return fmt.Errorf("account not found")
		}
		downloadRequestBucket := bucket.Bucket([]byte(downloadRequestBucketName))
		if downloadRequestBucket == nil {
			return fmt.Errorf("download request bucket not found")
		}

		return downloadRequestBucket.Delete([]byte(mediaItemId))
	})
}
