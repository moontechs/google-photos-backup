package google_client

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

const (
	PhotosApiRequestLimitType = "photos-request"
	PhotosDownloadLimitType   = "photos-download"
	DriveApiRequestLimitType  = "drive-request"
)

type Limiter interface {
	LimitReached(clientId, limitType string) (bool, error)
	SetLimitReached(clientId, limitType string, limitReached bool) error
}

type Limits struct {
	PhotosScan     Limit `json:"photosScan"`
	PhotosDownload Limit `json:"photosDownload"`
	DriveRequests  Limit `json:"driveRequests"`
}

type Limit struct {
	Count     int64 `json:"count"`
	Timestamp int64 `json:"time"`
}

type limiter struct {
	repository Repository
}

func NewLimiter(repository Repository) limiter {
	return limiter{repository: repository}
}

func (l limiter) LimitReached(clientId, limitType string) (bool, error) {
	limitsJson, err := l.repository.GetLimits(clientId)
	if err != nil {
		return false, fmt.Errorf("get limits: %w", err)
	}

	if limitsJson == nil {
		return false, nil
	}

	var limits Limits
	err = json.Unmarshal(limitsJson, &limits)
	if err != nil {
		return false, fmt.Errorf("unmarshal limitsJson: %w", err)
	}

	switch limitType {
	case PhotosApiRequestLimitType:
		if limits.PhotosScan.Timestamp == 0 {
			return false, nil
		}
		if limits.PhotosScan.Timestamp >= time.Now().Unix() {
			return false, nil
		}
	case PhotosDownloadLimitType:
		if limits.PhotosDownload.Timestamp == 0 {
			return false, nil
		}
		if limits.PhotosDownload.Timestamp >= time.Now().Unix() {
			return false, nil
		}
	case DriveApiRequestLimitType:
		if limits.DriveRequests.Timestamp == 0 {
			return false, nil
		}
		if limits.DriveRequests.Timestamp >= time.Now().Unix() {
			return false, nil
		}
	default:
		return false, fmt.Errorf("unknown limit type: %s", limitType)
	}

	return true, nil
}

func (l limiter) SetLimitReached(clientId, limitType string, limitReached bool) error {
	limitsJson, err := l.repository.GetLimits(clientId)
	if err != nil {
		return fmt.Errorf("get limits: %w", err)
	}

	limits := Limits{}

	if limitsJson != nil {
		err = json.Unmarshal(limitsJson, &limits)
		if err != nil {
			return fmt.Errorf("unmarshal limitsJson: %w", err)
		}
	}

	switch limitType {
	case PhotosApiRequestLimitType:
		if limitReached {
			limits.PhotosScan = l.caclulateLimits(limits.PhotosScan)
		} else {
			limits.PhotosScan.Count = 0
			limits.PhotosScan.Timestamp = 0
		}
	case PhotosDownloadLimitType:
		if limitReached {
			limits.PhotosDownload = l.caclulateLimits(limits.PhotosDownload)
		} else {
			limits.PhotosDownload.Count = 0
			limits.PhotosDownload.Timestamp = 0
		}
	case DriveApiRequestLimitType:
		if limitReached {
			limits.DriveRequests = l.caclulateLimits(limits.DriveRequests)
		} else {
			limits.DriveRequests.Count = 0
			limits.DriveRequests.Timestamp = 0
		}
	default:
		return fmt.Errorf("unknown limit type: %s", limitType)
	}

	limitsJson, err = json.Marshal(limits)
	if err != nil {
		return fmt.Errorf("marshal limits: %w", err)
	}

	l.repository.CreateUpdateLimits(clientId, limitsJson)

	return nil
}

func (l limiter) caclulateLimits(limit Limit) Limit {
	limit.Count++
	limit.Timestamp = time.Now().Add(time.Duration(math.Pow(3, float64(limit.Count))) * time.Minute).Unix()

	return limit
}
