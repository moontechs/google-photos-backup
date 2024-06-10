package scanner

import (
	"context"
	"encoding/json"
	"fmt"

	"google-backup/internal/account"
	"google-backup/internal/settings"

	"golang.org/x/sync/errgroup"
)

type UpdatesScanner interface {
	ScanAll(ctx context.Context) error
}

type AccountScanner interface {
	Scan(ctx context.Context, settingsData settings.SettingsData, account account.AccountData) error
}

type updatesScanner struct {
	accountRepository  account.Repository
	settingsRepository settings.Repository
	accountScanner     AccountScanner
}

func NewUpdatesScanner(
	accountRepository account.Repository,
	settingsRepository settings.Repository,
	accountScanner AccountScanner,
) updatesScanner {
	return updatesScanner{
		accountRepository:  accountRepository,
		settingsRepository: settingsRepository,
		accountScanner:     accountScanner,
	}
}

func (u updatesScanner) ScanAll(ctx context.Context) error {
	settingsJson, err := u.settingsRepository.Find()
	if err != nil {
		return fmt.Errorf("find settings: %w", err)
	}

	var settingsData settings.SettingsData
	err = json.Unmarshal(settingsJson, &settingsData)
	if err != nil {
		return fmt.Errorf("unmarshal settings: %w", err)
	}

	accountsJson, err := u.accountRepository.FindAccounts()
	if err != nil {
		return fmt.Errorf("get accounts: %w", err)
	}

	var accounts []account.AccountData

	for _, accountJson := range accountsJson {
		var accountData account.AccountData
		err = json.Unmarshal(accountJson, &accountData)
		if err != nil {
			return fmt.Errorf("unmarshal account: %w", err)
		}
		accounts = append(accounts, accountData)
	}

	errs, ctx := errgroup.WithContext(ctx)

	for _, acc := range accounts {
		func(s settings.SettingsData, a account.AccountData) {
			errs.Go(
				func() error {
					err := u.accountScanner.Scan(ctx, s, a)
					if err != nil {
						return fmt.Errorf("account scan: %w", err)
					}

					return nil
				},
			)
		}(settingsData, acc)
	}

	return errs.Wait()
}
