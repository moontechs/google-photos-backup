package account

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

const (
	ApiRequestLimitType = "request"
	DownloadLimitType   = "download"
)

type Limiter interface {
	LimitReached(email, limitType string) (bool, error)
	SetLimitReached(email, limitType string, limitReached bool) error
}

type Limits struct {
	Scan     Limit `json:"scan"`
	Download Limit `json:"download"`
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

func (l limiter) LimitReached(email, limitType string) (bool, error) {
	limitsJson, err := l.repository.GetLimits(email)
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
	case ApiRequestLimitType:
		if limits.Scan.Timestamp == 0 {
			return false, nil
		}
		if limits.Scan.Timestamp >= time.Now().Unix() {
			return false, nil
		}
	case DownloadLimitType:
		if limits.Download.Timestamp == 0 {
			return false, nil
		}
		if limits.Download.Timestamp >= time.Now().Unix() {
			return false, nil
		}
	default:
		return false, fmt.Errorf("unknown limit type: %s", limitType)
	}

	return true, nil
}

func (l limiter) SetLimitReached(email, limitType string, limitReached bool) error {
	limitsJson, err := l.repository.GetLimits(email)
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
	case ApiRequestLimitType:
		if limitReached {
			limits.Scan = l.caclulateLimits(limits.Scan)
		} else {
			limits.Scan.Count = 0
			limits.Scan.Timestamp = 0
		}
	case DownloadLimitType:
		if limitReached {
			limits.Download = l.caclulateLimits(limits.Download)
		} else {
			limits.Download.Count = 0
			limits.Download.Timestamp = 0
		}
	default:
		return fmt.Errorf("unknown limit type: %s", limitType)
	}

	limitsJson, err = json.Marshal(limits)
	if err != nil {
		return fmt.Errorf("marshal limits: %w", err)
	}

	l.repository.CreateUpdateLimits(email, limitsJson)

	return nil
}

func (l limiter) caclulateLimits(limit Limit) Limit {
	limit.Count++
	limit.Timestamp = time.Now().Add(time.Duration(math.Pow(3, float64(limit.Count))) * time.Minute).Unix()

	return limit
}
