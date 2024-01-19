package scanner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"google-backup/internal/account"
	"google-backup/internal/downloader"
	"google-backup/internal/media"
	"google-backup/internal/media_reader"

	"golang.org/x/sync/errgroup"
)

type UpdatesScanner interface {
	ScanAll(ctx context.Context) error
}

type updatesScanner struct {
	repository        Repository
	downloadScheduler downloader.Scheduler
	accountLimiter    account.Limiter
	mediaReader       media_reader.Reader
}

func NewUpdatesScanner(
	repository Repository,
	downloadScheduler downloader.Scheduler,
	accountLimiter account.Limiter,
	mediaReader media_reader.Reader,
) updatesScanner {
	return updatesScanner{
		repository:        repository,
		downloadScheduler: downloadScheduler,
		accountLimiter:    accountLimiter,
		mediaReader:       mediaReader,
	}
}

func (u updatesScanner) ScanAll(ctx context.Context) error {
	readers, err := u.mediaReader.CreateMediaReaders(ctx)
	if err != nil {
		return fmt.Errorf("create media readers: %w", err)
	}

	errs, ctx := errgroup.WithContext(ctx)

	for email, reader := range readers {
		r := reader
		e := email

		errs.Go(
			func() error {
				err := u.scan(ctx, r, e)
				if err != nil {
					return fmt.Errorf("scan updates: %w", err)
				}

				return nil
			},
		)
	}

	return errs.Wait()
}

func (u updatesScanner) scan(ctx context.Context, mediaReader media.Reader, email string) error {
	rescanRequest, err := u.getRescanRequest(email)
	if err != nil {
		return fmt.Errorf("get rescan request: %w", err)
	}

	mediaItems, err := mediaReader.GetMediaItems(string(email), rescanRequest.NextPageToken)
	if err != nil {
		if errors.As(err, &media.TooManyRequestsError{}) {
			u.accountLimiter.SetLimitReached(string(email), account.ApiRequestLimitType, true)
		}

		return fmt.Errorf("get media items: %w", err)
	}

	u.accountLimiter.SetLimitReached(string(email), account.ApiRequestLimitType, false)

	if len(mediaItems.Items) == 0 {
		return u.repository.DeleteRescanRequest(email)
	}

	for _, item := range mediaItems.Items {
		err := u.downloadScheduler.ScheduleDownload(email, item.ID)
		if err != nil {
			return fmt.Errorf("schedule download: %w", err)
		}
	}

	rescanRequest.NextPageToken = mediaItems.NextPageToken
	if rescanRequest.NextPageToken == "" {
		return u.repository.DeleteRescanRequest(email)

	}

	rescanRequestJson, err := json.Marshal(rescanRequest)
	if err != nil {
		return fmt.Errorf("marshal next page token: %w", err)
	}

	return u.repository.UpdateRescanRequest(email, rescanRequestJson)
}

func (u updatesScanner) getRescanRequest(email string) (RescanRequest, error) {
	rescanRequestJson, err := u.repository.GetRescanRequest(string(email))
	if err != nil {
		return RescanRequest{}, fmt.Errorf("get rescan request: %w", err)
	}

	if len(rescanRequestJson) == 0 {
		return RescanRequest{}, fmt.Errorf("rescan request not found")
	}

	var rescanRequest RescanRequest
	err = json.Unmarshal(rescanRequestJson, &rescanRequest)
	if err != nil {
		return RescanRequest{}, fmt.Errorf("unmarshal next batch: %w", err)
	}

	return rescanRequest, err
}
