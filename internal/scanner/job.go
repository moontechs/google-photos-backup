package scanner

import (
	"context"
	"fmt"
	"time"

	"github.com/moontechs/photos-backup/internal/config"
	log "github.com/sirupsen/logrus"
)

type scannerJob struct {
	updatesScanner UpdatesScanner
	config         config.Config
}

func NewScannerJob(updatesScanner UpdatesScanner, config config.Config) scannerJob {
	return scannerJob{updatesScanner: updatesScanner, config: config}
}

func (s scannerJob) Run(ctx context.Context) error {
	err := s.updatesScanner.ScanAll(ctx)
	if err != nil {
		log.Error(fmt.Errorf("scan all: %w", err))
	}

	return err
}

func (s scannerJob) GetDelay() time.Duration {
	configData := s.config.Get()

	return configData.ScannerJobDelay
}

func (s scannerJob) GetName() string {
	return "scanner"
}
