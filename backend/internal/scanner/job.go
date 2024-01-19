package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google-backup/internal/settings"

	log "github.com/sirupsen/logrus"
)

type scannerJob struct {
	updatesScanner     UpdatesScanner
	settingsRepository settings.Repository
}

func NewScannerJob(updatesScanner UpdatesScanner, settingsRepository settings.Repository) scannerJob {
	return scannerJob{updatesScanner: updatesScanner, settingsRepository: settingsRepository}
}

func (s scannerJob) Run(ctx context.Context) error {
	err := s.updatesScanner.ScanAll(ctx)
	if err != nil {
		log.Error(fmt.Errorf("scan all: %w", err))
	}

	return err
}

func (s scannerJob) GetDelay() time.Duration {
	settingsJson, err := s.settingsRepository.Find()
	if err != nil {
		return 1 * time.Minute
	}

	var settingsData settings.SettingsData
	err = json.Unmarshal(settingsJson, &settingsData)
	if err != nil {
		return 1 * time.Minute
	}

	return settingsData.ScannerJobDelay
}

func (s scannerJob) GetName() string {
	return "scanner"
}
