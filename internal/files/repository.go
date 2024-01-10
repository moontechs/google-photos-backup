package files

import (
	"fmt"

	"go.etcd.io/bbolt"
)

const (
	downloadErrorsBucketName = "download_errors"
	filesMetaDataBucketName  = "files_meta_data"
)

type Repository interface {
	SaveDownloadError(email string, mediaItemId, message string) error
	SaveFileMeta(email string, key, data []byte) error
	GetFileMeta(email string, key []byte) ([]byte, error)
}

type repository struct {
	db *bbolt.DB
}

func NewRepository(db *bbolt.DB) repository {
	return repository{db: db}
}

func (r repository) SaveDownloadError(email string, mediaItemId string, message string) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(email))
		if err != nil {
			return fmt.Errorf("create bucket: %w", err)
		}

		downloadErrorsBucket, err := bucket.CreateBucketIfNotExists([]byte(downloadErrorsBucketName))
		if err != nil {
			return fmt.Errorf("create download errors bucket: %w", err)
		}

		return downloadErrorsBucket.Put([]byte(mediaItemId), []byte(message))
	})
}

func (r repository) SaveFileMeta(email string, key, data []byte) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(email))
		if err != nil {
			return fmt.Errorf("create bucket: %w", err)
		}

		filesMetaDataBucket, err := bucket.CreateBucketIfNotExists([]byte(filesMetaDataBucketName))
		if err != nil {
			return fmt.Errorf("create files meta data bucket: %w", err)
		}

		return filesMetaDataBucket.Put(key, data)
	})
}

func (r repository) GetFileMeta(email string, key []byte) ([]byte, error) {
	var data []byte
	data = nil

	err := r.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(email))
		if bucket == nil {
			return nil
		}

		filesMetaDataBucket := bucket.Bucket([]byte(filesMetaDataBucketName))
		if filesMetaDataBucket == nil {
			return nil
		}

		data = filesMetaDataBucket.Get(key)

		return nil
	})

	return data, err
}
