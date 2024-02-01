package google_client

import (
	"go.etcd.io/bbolt"
)

const (
	clientBucketName           = "clients"
	assignedAccountsBucketName = "assigned_accounts"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Repository
type Repository interface {
	Find(key string) ([]byte, error)
	FindAll() (map[string][]byte, error)
	FindAssignedAccounts(clientId string) ([]byte, error)
	FindAllAssignedAccounts() (map[string][]byte, error)
	Save(key string, value []byte) error
	SaveAssignedAccounts(clientId string, value []byte) error
	Delete(key string) error
}

type repository struct {
	DB *bbolt.DB
}

func NewRepository(db *bbolt.DB) repository {
	return repository{DB: db}
}

func (r repository) Find(key string) ([]byte, error) {
	var value []byte
	value = nil

	err := r.DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(clientBucketName))
		if bucket == nil {
			return nil
		}

		value = bucket.Get([]byte(key))

		return nil
	})

	return value, err
}

func (r repository) FindAll() (map[string][]byte, error) {
	values := make(map[string][]byte)

	err := r.DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(clientBucketName))
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			values[string(k)] = v

			return nil
		})
	})

	return values, err
}

func (r repository) FindAssignedAccounts(clientId string) ([]byte, error) {
	var value []byte
	value = nil

	err := r.DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(assignedAccountsBucketName))
		if bucket == nil {
			return nil
		}

		value = bucket.Get([]byte(clientId))

		return nil
	})

	return value, err
}

func (r repository) FindAllAssignedAccounts() (map[string][]byte, error) {
	values := make(map[string][]byte)

	err := r.DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(assignedAccountsBucketName))
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			values[string(k)] = v

			return nil
		})
	})

	return values, err
}

func (r repository) Save(key string, value []byte) error {
	return r.DB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(clientBucketName))
		if err != nil {
			return err
		}

		return bucket.Put([]byte(key), value)
	})
}

func (r repository) SaveAssignedAccounts(clientId string, value []byte) error {
	return r.DB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(assignedAccountsBucketName))
		if err != nil {
			return err
		}

		return bucket.Put([]byte(clientId), value)
	})
}

func (r repository) Delete(key string) error {
	return r.DB.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(clientBucketName))
		if bucket == nil {
			return nil
		}

		return bucket.Delete([]byte(key))
	})
}
