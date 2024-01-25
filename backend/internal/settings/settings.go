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
	RootPath                 string        `json:"rootPath"`
	PhotosScannerJobDelay    time.Duration `json:"photosScannerJobDelay"`
	PhotosDownloaderJobDelay time.Duration `json:"photosDownloaderJobDelay"`
	Host                     string        `json:"host"`
	PhotosBackupEnabled      bool          `json:"photosBackupEnabled"`
	DriveBackupEnabled       bool          `json:"driveBackupEnabled"`
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
