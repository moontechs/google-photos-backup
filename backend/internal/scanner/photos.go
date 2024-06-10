package scanner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"google-backup/internal/account"
	"google-backup/internal/downloader"
	"google-backup/internal/google_client"
	"google-backup/internal/photos"
	"google-backup/internal/settings"
)

type photosScanner struct {
	repository        Repository
	readerCreater     photos.ReaderCreater
	limiter           google_client.Limiter
	downloadScheduler downloader.Scheduler
}

func NewPhotosScanner(
	repository Repository,
	readerCreater photos.ReaderCreater,
	limiter google_client.Limiter,
	downloadScheduler downloader.Scheduler,
) photosScanner {
	return photosScanner{
		repository:        repository,
		readerCreater:     readerCreater,
		limiter:           limiter,
		downloadScheduler: downloadScheduler,
	}
}

func (p photosScanner) Scan(
	ctx context.Context,
	settingsData settings.SettingsData,
	account account.AccountData,
) error {
	if !settingsData.PhotosBackupEnabled {
		return nil
	}

	mediaReader, err := p.readerCreater.CreateMediaReader(ctx, account)
	if err != nil {
		return fmt.Errorf("create media reader: %w", err)
	}

	rescanRequest, err := p.getRescanRequest(account.Email)
	if err != nil {
		return fmt.Errorf("get rescan request: %w", err)
	}

	limitReached, err := p.limiter.LimitReached(mediaReader.GetClientId(), google_client.PhotosApiRequestLimitType)
	if err != nil {
		return fmt.Errorf("limit reached: %w", err)
	}

	if limitReached {
		return nil
	}

	mediaItems, err := mediaReader.GetMediaItems(string(account.Email), rescanRequest.NextPageToken)
	if err != nil {
		if errors.As(err, &photos.TooManyRequestsError{}) {
			p.limiter.SetLimitReached(mediaReader.GetClientId(), google_client.PhotosApiRequestLimitType, true)
		}

		return fmt.Errorf("get media items: %w", err)
	}

	p.limiter.SetLimitReached(mediaReader.GetClientId(), google_client.PhotosApiRequestLimitType, false)

	if len(mediaItems.Items) == 0 {
		return p.repository.DeleteRescanRequest("photos", account.Email)
	}

	for _, item := range mediaItems.Items {
		err := p.downloadScheduler.ScheduleDownload(account.Email, item.ID)
		if err != nil {
			return fmt.Errorf("schedule download: %w", err)
		}
	}

	rescanRequest.NextPageToken = mediaItems.NextPageToken
	if rescanRequest.NextPageToken == "" {
		return p.repository.DeleteRescanRequest("photos", account.Email)
	}

	rescanRequestJson, err := json.Marshal(rescanRequest)
	if err != nil {
		return fmt.Errorf("marshal next page token: %w", err)
	}

	return p.repository.UpdateRescanRequest("photos", account.Email, rescanRequestJson)
}

func (p photosScanner) getRescanRequest(email string) (RescanRequest, error) {
	rescanRequestJson, err := p.repository.GetRescanRequest("photos", string(email))
	if err != nil {
		return RescanRequest{}, fmt.Errorf("get rescan request: %w", err)
	}

	if rescanRequestJson == nil {
		return RescanRequest{}, nil
	}

	var rescanRequest RescanRequest
	err = json.Unmarshal(rescanRequestJson, &rescanRequest)
	if err != nil {
		return RescanRequest{}, fmt.Errorf("unmarshal next batch: %w", err)
	}

	return rescanRequest, err
}
