package photos

import (
	"context"
	"fmt"

	"google-backup/internal/account"
	"google-backup/internal/auth"
)

type ReaderCreater interface {
	CreateMediaReader(ctx context.Context, account account.AccountData) (Reader, error)
}

type readerCreater struct {
	account    account.Account
	googleAuth auth.Auth
}

func NewReaderCreater(
	account account.Account,
	googleAuth auth.Auth,
) readerCreater {
	return readerCreater{
		account:    account,
		googleAuth: googleAuth,
	}
}

func (r readerCreater) CreateMediaReader(ctx context.Context, account account.AccountData) (Reader, error) {
	authToken, err := r.account.GetTokenByEmail(account.Email)
	if err != nil {
		return nil, fmt.Errorf("get token by email: %w", err)
	}

	clientId, err := r.account.GetAccountOauthClientId(account.Email)
	if err != nil {
		return nil, fmt.Errorf("get account oauth client name: %w", err)
	}

	gClient, err := r.googleAuth.GetHttpClient(ctx, clientId, &authToken)
	if err != nil {
		return nil, fmt.Errorf("get google client: %w", err)
	}

	mediaReader, err := NewReader(clientId, gClient)
	if err != nil {
		return nil, fmt.Errorf("new media reader: %w", err)
	}

	return mediaReader, nil
}
