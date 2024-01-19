package scanner

import (
	"encoding/json"
	"fmt"
)

type Scheduler interface {
	ScheduleRescan(email string) error
}

type scheduler struct {
	repository Repository
}

type RescanRequest struct {
	NextPageToken string `json:"next_page_token"`
}

func NewScheduler(repository Repository) *scheduler {
	return &scheduler{repository: repository}
}

func (s *scheduler) ScheduleRescan(email string) error {
	rescanRequest := RescanRequest{
		NextPageToken: "",
	}

	rescanRequestJson, err := json.Marshal(rescanRequest)
	if err != nil {
		return fmt.Errorf("marshal scan data: %w", err)
	}

	return s.repository.UpdateRescanRequest(email, rescanRequestJson)
}
