package account

import (
	"errors"
	"fmt"
	"strconv"

	"go.etcd.io/bbolt"
)

const (
	accountLimitsKey   = "limits"
	oauthClientNameKey = "oauth_client_name"
	tokensBucketName   = "tokens"
	accountsBucketName = "accounts"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Repository
type Repository interface {
	SaveToken(email string, token []byte) error
	SaveAccount(email string, userInfo []byte) error
	GetAccounts() ([][]byte, error)
	FindAccount(email string) ([]byte, error)
	AccountExist(email string) (bool, error)
	FindTokenByEmail(email string) ([]byte, error)
	CreateUpdateLimits(email string, limits []byte) error
	GetLimits(email string) ([]byte, error)
	GetAccountOauthClientName(email string) ([]byte, error)
}

type repo struct {
	DB *bbolt.DB
}

func NewRepository(db *bbolt.DB) *repo {
	return &repo{DB: db}
}

func (r *repo) SaveToken(email string, token []byte) error {
	return r.DB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(tokensBucketName))
		if err != nil {
			return err
		}

		return bucket.Put([]byte(email), token)
	})
}

func (r *repo) SaveAccount(email string, userInfo []byte) error {
	return r.DB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(accountsBucketName))
		if err != nil {
			return err
		}

		return bucket.Put([]byte(email), userInfo)
	})
}

func (r *repo) SaveAccountOauthClientName(email string, clientName []byte) error {
	return r.DB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(email))
		if err != nil {
			return err
		}

		return bucket.Put([]byte(oauthClientNameKey), clientName)
	})
}

func (r *repo) AccountExist(email string) (bool, error) {
	var exist bool
	err := r.DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(accountsBucketName))
		if bucket == nil {
			exist = false

			return nil
		}

		exist = bucket.Get([]byte(email)) != nil

		return nil
	})

	return exist, err
}

func (r *repo) FindAccount(email string) ([]byte, error) {
	var account []byte
	account = nil

	err := r.DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(accountsBucketName))
		if bucket == nil {
			return nil
		}

		account = bucket.Get([]byte(email))

		return nil
	})

	return account, err
}

func (r *repo) FindTokenByEmail(email string) ([]byte, error) {
	var token []byte
	token = nil

	err := r.DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(tokensBucketName))
		if bucket == nil {
			return nil
		}

		token = bucket.Get([]byte(email))

		return nil
	})

	return token, err
}

func (r *repo) SetLimitReached(email string, limitReached bool) error {
	return r.DB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(email))
		if err != nil {
			return err
		}

		return bucket.Put([]byte("limit_reached"), []byte(strconv.FormatBool(limitReached)))
	})
}

func (r *repo) GetLimitReached(email string) (bool, error) {
	limitReached := false

	err := r.DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(email))
		if bucket == nil {
			return errors.New("account not found")
		}
		limitReachedFromDB := bucket.Get([]byte("limit_reached"))

		if limitReachedFromDB == nil {
			limitReached = false
		}

		if string(limitReachedFromDB) == "true" {
			limitReached = true
		}

		return nil
	})

	return limitReached, err
}

func (r *repo) GetAccounts() ([][]byte, error) {
	var accounts [][]byte

	err := r.DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(accountsBucketName))
		if bucket == nil {
			return nil
		}

		bucket.ForEach(func(k, v []byte) error {
			accounts = append(accounts, v)

			return nil
		})

		return nil
	})

	return accounts, err
}

func (r *repo) GetLimits(email string) ([]byte, error) {
	var limits []byte

	r.DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(email))
		if bucket == nil {
			return nil
		}

		limits = bucket.Get([]byte(accountLimitsKey))

		return nil
	})

	return limits, nil
}

func (r *repo) CreateUpdateLimits(email string, limits []byte) error {
	return r.DB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(email))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		return bucket.Put([]byte(accountLimitsKey), limits)
	})
}

func (r *repo) GetAccountOauthClientName(email string) ([]byte, error) {
	var oauthClientName []byte

	err := r.DB.View(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(email))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		oauthClientName = bucket.Get([]byte(oauthClientNameKey))

		return nil
	})

	return oauthClientName, err
}
