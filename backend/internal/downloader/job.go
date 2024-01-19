package downloader

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google-backup/internal/settings"

	log "github.com/sirupsen/logrus"
)

type downloaderJob struct {
	downloader         Downloader
	settingsRepository settings.Repository
}

func NewDownloaderJob(downloader Downloader, settingsRepository settings.Repository) downloaderJob {
	return downloaderJob{downloader: downloader, settingsRepository: settingsRepository}
}

func (s downloaderJob) Run(ctx context.Context) error {
	err := s.downloader.DownloadAll(ctx)
	if err != nil {
		log.Error(fmt.Errorf("download all: %w", err))
	}

	return err
}

func (s downloaderJob) GetDelay() time.Duration {
	settingsJson, err := s.settingsRepository.Find()
	if err != nil {
		return 1 * time.Minute
	}

	var settingsData settings.SettingsData
	err = json.Unmarshal(settingsJson, &settingsData)
	if err != nil {
		return 1 * time.Minute
	}

	return settingsData.PhotosDownloaderJobDelay
}

func (s downloaderJob) GetName() string {
	return "downloader"
}
