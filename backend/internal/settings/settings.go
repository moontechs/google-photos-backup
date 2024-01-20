package settings

import (
	"encoding/json"
	"fmt"
	"time"
)

type SettingsInitializer interface {
	Init() error
}

type SettingsData struct {
	RootPath                 string        `json:"root_path"`
	PhotosScannerJobDelay    time.Duration `json:"photos_scanner_job_delay"`
	PhotosDownloaderJobDelay time.Duration `json:"photos_downloader_job_delay"`
	Host                     string        `json:"host"`
	PhotosBackupEnabled      bool          `json:"photos_backup_enabled"`
	DriveBackupEnabled       bool          `json:"drive_backup_enabled"`
}

type settings struct {
	repository Repository
}

func NewSettings(repository Repository) settings {
	return settings{
		repository: repository,
	}
}

// Saves default settings if doesn't exist
func (c settings) Init() error {
	defaultSettingsData := SettingsData{
		RootPath:                 "/data",
		PhotosScannerJobDelay:    time.Minute,
		PhotosDownloaderJobDelay: time.Minute,
		Host:                     "http://localhost:8080",
		PhotosBackupEnabled:      true,
		DriveBackupEnabled:       true,
	}

	defaultSettings, err := json.Marshal(defaultSettingsData)
	if err != nil {
		return fmt.Errorf("marshal default settings: %w", err)
	}

	err = c.repository.Save(defaultSettings)
	if err != nil {
		return fmt.Errorf("set config data: %w", err)
	}

	return nil
}
