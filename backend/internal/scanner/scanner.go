package scanner

import (
	"context"
	"encoding/json"
	"fmt"

	"google-backup/internal/account"
	"google-backup/internal/downloader"
	"google-backup/internal/media"
	"google-backup/internal/media_reader"

	"golang.org/x/sync/errgroup"
)

const (
	RescanTypePhotos = "photos"
	RescanTypeDrive  = "drive"
)

type UpdatesScanner interface {
	ScanAll(ctx context.Context) error
}

type RescanRequests struct {
	Photos RescanRequest `json:"photos"`
	Drive  RescanRequest `json:"drive"`
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
	// rescanRequests, err := u.getRescanRequests(email)
	// if err != nil {
	// 	return fmt.Errorf("get rescan requests: %w", err)
	// }

	// mediaItems, err := mediaReader.GetMediaItems(string(email), rescanRequest.NextPageToken)
	// if err != nil {
	// 	if errors.As(err, &media.TooManyRequestsError{}) {
	// 		u.accountLimiter.SetLimitReached(string(email), account.ApiRequestLimitType, true)
	// 	}

	// 	return fmt.Errorf("get media items: %w", err)
	// }

	// u.accountLimiter.SetLimitReached(string(email), account.ApiRequestLimitType, false)

	// if len(mediaItems.Items) == 0 {
	// 	return u.repository.DeleteRescanRequest(email)
	// }

	// for _, item := range mediaItems.Items {
	// 	err := u.downloadScheduler.ScheduleDownload(email, item.ID)
	// 	if err != nil {
	// 		return fmt.Errorf("schedule download: %w", err)
	// 	}
	// }

	// rescanRequest.NextPageToken = mediaItems.NextPageToken
	// if rescanRequest.NextPageToken == "" {
	// 	return u.repository.DeleteRescanRequest(email)

	// }

	// rescanRequestJson, err := json.Marshal(rescanRequest)
	// if err != nil {
	// 	return fmt.Errorf("marshal next page token: %w", err)
	// }

	// return u.repository.UpdateRescanRequest(email, rescanRequestJson)

	return nil
}

func (u updatesScanner) getRescanRequests(email string) (RescanRequests, error) {
	rescanRequestsMap, err := u.repository.GetRescanRequests(string(email))
	if err != nil {
		return RescanRequests{}, fmt.Errorf("get rescan request: %w", err)
	}

	var rescanRequests RescanRequests

	for rescanType, rescanRequestJson := range rescanRequestsMap {
		if rescanRequestJson == nil {
			continue
		}

		var rescanRequest RescanRequest
		err = json.Unmarshal(rescanRequestJson, &rescanRequest)
		if err != nil {
			return RescanRequests{}, fmt.Errorf("unmarshal next batch: %w", err)
		}

		switch rescanType {
		case RescanTypePhotos:
			rescanRequests.Photos = rescanRequest
		case RescanTypeDrive:
			rescanRequests.Drive = rescanRequest
		}
	}

	return rescanRequests, err
}
