package downloader

import (
	"context"
	"fmt"
	"time"

	"github.com/moontechs/photos-backup/internal/config"
	log "github.com/sirupsen/logrus"
)

type downloaderJob struct {
	downloader Downloader
	config     config.Config
}

func NewDownloaderJob(downloader Downloader, config config.Config) downloaderJob {
	return downloaderJob{downloader: downloader, config: config}
}

func (s downloaderJob) Run(ctx context.Context) error {
	err := s.downloader.DownloadAll(ctx)
	if err != nil {
		log.Error(fmt.Errorf("download all: %w", err))
	}

	return err
}

func (s downloaderJob) GetDelay() time.Duration {
	configData := s.config.Get()

	return configData.DownloaderJobDelay
}

func (s downloaderJob) GetName() string {
	return "downloader"
}
