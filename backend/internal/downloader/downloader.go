package downloader

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"google-backup/internal/account"
	"google-backup/internal/files"
	"google-backup/internal/settings"

	"golang.org/x/sync/errgroup"
)

const downloadsBatchLimit = 50

type Downloader interface {
	DownloadAll(ctx context.Context) error
}

type AccountDownloader interface {
	Download(ctx context.Context, settingsData settings.SettingsData, account account.AccountData) error
}

type TooManyRequestsError struct {
	error
}

type NotOkRequestError struct {
	error
}

type ManuallyRepeatableError struct {
	error
}

type downloader struct {
	repository         Repository
	accountRepository  account.Repository
	httpClient         *http.Client
	settingsRepository settings.Repository
	filesManager       files.FilesManager
	accountDownloader  AccountDownloader
}

func NewDownloader(
	repository Repository,
	accountRepository account.Repository,
	httpClient *http.Client,
	settingsRepository settings.Repository,
	filesManager files.FilesManager,
	accountDownloader AccountDownloader,
) downloader {
	return downloader{
		repository:         repository,
		accountRepository:  accountRepository,
		httpClient:         httpClient,
		settingsRepository: settingsRepository,
		filesManager:       filesManager,
		accountDownloader:  accountDownloader,
	}
}

func (d downloader) DownloadAll(ctx context.Context) error {
	settingsJson, err := d.settingsRepository.Find()
	if err != nil {
		return fmt.Errorf("find settings: %w", err)
	}

	var settingsData settings.SettingsData
	err = json.Unmarshal(settingsJson, &settingsData)
	if err != nil {
		return fmt.Errorf("unmarshal settings: %w", err)
	}

	accountsJson, err := d.accountRepository.FindAccounts()
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
					err := d.accountDownloader.Download(ctx, s, a)
					if err != nil {
						return fmt.Errorf("account download: %w", err)
					}

					return nil
				},
			)
		}(settingsData, acc)
	}

	return errs.Wait()
}

func (d downloader) download(ctx context.Context, mediaReader media.Reader, email string) error {
	counter := 0

	for counter < downloadsBatchLimit {
		fileMeta, err := d.downloadFromBaseUrl(ctx, mediaReader, email)
		if err != nil {
			if errors.As(err, &media.TooManyRequestsError{}) {
				d.accountLimiter.SetLimitReached(string(email), account.ApiRequestLimitType, true)

				return fmt.Errorf("download from base url: %w", err)
			}

			if errors.As(err, &TooManyRequestsError{}) {
				d.accountLimiter.SetLimitReached(string(email), account.DownloadLimitType, true)

				return fmt.Errorf("download from base url: %w", err)
			}

			if fileMeta.MediaItem.ID != "" {
				d.repository.DeleteDownloadRequest(email, fileMeta.MediaItem.ID)
			}

			// TODO add to errors reporter and continue loop. Files package has a function
		}

		if fileMeta.MediaItem.ID != "" {
			err = d.repository.DeleteDownloadRequest(email, fileMeta.MediaItem.ID)
			if err != nil {
				return fmt.Errorf("delete download request: %w", err)
			}
		}

		err = d.filesManager.SaveFileMeta(email, fileMeta)
		if err != nil {
			return fmt.Errorf("save file meta: %w", err)
		}

		counter++
	}

	return nil
}

func (d downloader) downloadFromBaseUrl(ctx context.Context, mediaReader media.Reader, email string) (files.FileMeta, error) {
	fileMeta := files.FileMeta{}

	downloadRequestJson, err := d.repository.GetDownloadRequest(email)
	if err != nil {
		return fileMeta, fmt.Errorf("get download request: %w", err)
	}

	if downloadRequestJson == nil {
		return fileMeta, nil
	}

	var downloadRequest DownloadRequest
	err = json.Unmarshal(downloadRequestJson, &downloadRequest)
	if err != nil {
		return fileMeta, fmt.Errorf("unmarshal download request: %w", err)
	}

	if downloadRequest.MediaItemId == "" {
		return fileMeta, nil
	}

	limitReached, err := d.accountLimiter.LimitReached(email, account.ApiRequestLimitType)
	if err != nil {
		return fileMeta, fmt.Errorf("limit reached: %w", err)
	}

	if limitReached {
		return fileMeta, nil
	}

	mediaItem, err := mediaReader.GetMediaItem(downloadRequest.MediaItemId)
	if err != nil {
		return fileMeta, fmt.Errorf("get media item: %w", err)
	}

	fileMeta.MediaItem = mediaItem

	filePathName, err := d.filesManager.GenerateFilePathName(email, mediaItem)
	if err != nil {
		return fileMeta, fmt.Errorf("create file path name: %w", err)
	}

	fileExists, err := d.filesManager.FileExists(email, mediaItem)
	if err != nil {
		return fileMeta, fmt.Errorf("file exists: %w", err)
	}

	err = d.filesManager.CreateFolderIfDoesNotExist(filePathName)
	if err != nil {
		return fileMeta, fmt.Errorf("create folder: %w", err)
	}

	err = d.downloadFile(
		filePathName,
		mediaItem.BaseUrl,
		func(reader io.Reader) (bool, error) {
			if fileExists {
				return d.filesManager.EqualHash(filePathName, reader)
			} else {
				return true, nil
			}
		},
	)
	if err != nil {
		return fileMeta, fmt.Errorf("download file: %w", err)
	}

	err = d.filesManager.UpdateCreationTime(filePathName, mediaItem.MediaMetadata.CreationTime)
	if err != nil {
		return fileMeta, fmt.Errorf("update creation time: %w", err)
	}

	fileMeta.FilePathName = filePathName

	return fileMeta, nil
}

func (d downloader) downloadFile(
	filePathName string,
	url string,
	shouldDownload func(reader io.Reader) (bool, error),
) error {
	resp, err := d.httpClient.Get(url + "=d")
	if err != nil {
		return fmt.Errorf("download file: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return TooManyRequestsError{}
	}

	if resp.StatusCode != http.StatusOK {
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response body: %w", err)
		}

		return NotOkRequestError{fmt.Errorf(string(responseBody))}
	}

	ok, err := shouldDownload(resp.Body)
	if err != nil {
		return fmt.Errorf("should download: %w", err)
	}

	if !ok {
		return nil
	}

	out, err := os.Create(d.filesManager.AddRootFolderToPath(filePathName))
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("copy file: %w", err)
	}

	return nil
}
