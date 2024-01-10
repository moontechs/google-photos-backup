package config

import (
	"encoding/json"
	"fmt"
	"time"
)

type Config interface {
	Get() ConfigData
	Set(configdata ConfigData) error
	Init() error
}

type ConfigData struct {
	RootPath           string        `json:"root_path"`
	ScannerJobDelay    time.Duration `json:"scanner_job_delay"`
	DownloaderJobDelay time.Duration `json:"downloader_job_delay"`
}

type config struct {
	repository Repository
}

func NewConfig(repository Repository) config {
	return config{
		repository: repository,
	}
}

func (c config) Get() ConfigData {
	defaultConfig := ConfigData{
		RootPath:           "/data",
		ScannerJobDelay:    time.Minute,
		DownloaderJobDelay: time.Minute,
	}

	data, err := c.repository.Get()
	if err != nil {
		return defaultConfig
	}

	if data == nil {
		return defaultConfig
	}

	var configData ConfigData
	err = json.Unmarshal(data, &configData)
	if err != nil {
		return defaultConfig
	}

	return configData
}

func (c config) Set(configData ConfigData) error {
	data, err := json.Marshal(configData)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	err = c.repository.Set(data)
	if err != nil {
		return fmt.Errorf("set config: %w", err)
	}

	return nil
}

// Saves default config if it doesn't exist
func (c config) Init() error {
	configData := c.Get()

	err := c.Set(configData)
	if err != nil {
		return fmt.Errorf("set config data: %w", err)
	}

	return nil
}
