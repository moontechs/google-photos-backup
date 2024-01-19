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
	RootPath           string        `json:"root_path"`
	ScannerJobDelay    time.Duration `json:"scanner_job_delay"`
	DownloaderJobDelay time.Duration `json:"downloader_job_delay"`
	Domain             string        `json:"domain"`
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
		RootPath:           "/data",
		ScannerJobDelay:    time.Minute,
		DownloaderJobDelay: time.Minute,
		Domain:             "http://localhost:8080",
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
