package downloader

import (
	"encoding/json"
	"fmt"
)

type Scheduler interface {
	ScheduleDownload(email string, mediaItemId string) error
}

type scheduler struct {
	repository Repository
}

func NewScheduler(repository Repository) scheduler {
	return scheduler{repository: repository}
}

func (s scheduler) ScheduleDownload(email string, mediaItemId string) error {
	downloadRequest := DownloadRequest{
		MediaItemId: mediaItemId,
	}

	downloadRequestJson, err := json.Marshal(downloadRequest)
	if err != nil {
		return fmt.Errorf("marshal download data: %w", err)
	}

	return s.repository.UpdateDownloadRequest(email, mediaItemId, downloadRequestJson)
}
