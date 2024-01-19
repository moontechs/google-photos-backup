package auth

import "go.etcd.io/bbolt"

const (
	oauthClientsBucket = "oauth_clients"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Repository
type Repository interface {
	SaveOauthClientData(clientName string, data []byte) error
	GetOauthClientData(clientName string) ([]byte, error)
}

type repository struct {
	DB *bbolt.DB
}

func NewRepository(db *bbolt.DB) repository {
	return repository{DB: db}
}

func (r repository) SaveOauthClientData(clientName string, data []byte) error {
	return r.DB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(oauthClientsBucket))
		if err != nil {
			return err
		}

		return bucket.Put([]byte(clientName), data)
	})
}

func (r repository) GetOauthClientData(clientName string) ([]byte, error) {
	var data []byte

	err := r.DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(oauthClientsBucket))
		if bucket == nil {
			return nil
		}

		data = bucket.Get([]byte(clientName))

		return nil
	})

	return data, err
}
