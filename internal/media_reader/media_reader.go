package media_reader

import (
	"context"
	"fmt"

	"github.com/moontechs/photos-backup/internal/account"
	"github.com/moontechs/photos-backup/internal/auth"
	"github.com/moontechs/photos-backup/internal/media"
)

type Reader interface {
	CreateMediaReaders(ctx context.Context) (map[string]media.Reader, error)
}

type reader struct {
	account        account.Account
	googleAuth     auth.Auth
	accountLimiter account.Limiter
}

func NewMediaReader(
	account account.Account,
	googleAuth auth.Auth,
	accountLimiter account.Limiter,
) reader {
	return reader{
		account:        account,
		googleAuth:     googleAuth,
		accountLimiter: accountLimiter,
	}
}

func (r reader) CreateMediaReaders(ctx context.Context) (map[string]media.Reader, error) {
	accounts, err := r.account.GetAccounts()
	if err != nil {
		return nil, fmt.Errorf("get accounts: %w", err)
	}

	readers := make(map[string]media.Reader, len(accounts))

	for _, email := range accounts {
		limitReached, err := r.accountLimiter.LimitReached(string(email), account.ApiRequestLimitType)
		if err != nil {
			return nil, fmt.Errorf("limit reached check: %w", err)
		}

		if limitReached {
			continue
		}

		authToken, err := r.account.GetTokenByEmail(string(email))
		if err != nil {
			return nil, fmt.Errorf("get token by email: %w", err)
		}

		clientName, err := r.account.GetAccountOauthClientName(string(email))
		if err != nil {
			return nil, fmt.Errorf("get account oauth client name: %w", err)
		}

		gClient, err := r.googleAuth.GetClient(ctx, clientName, &authToken)
		if err != nil {
			return nil, fmt.Errorf("get google client: %w", err)
		}

		mediaReader, err := media.NewReader(gClient)
		if err != nil {
			return nil, fmt.Errorf("new media reader: %w", err)
		}

		readers[string(email)] = mediaReader
	}

	return readers, nil
}
