package media

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Reader interface {
	GetMediaItems(email string, nextPageToken string) (MediaItems, error)
	GetMediaItem(mediaItemId string) (MediaItem, error)
}

type TooManyRequestsError struct {
	error
}

type NotOkRequestError struct {
	error
}

type reader struct {
	httpClient *http.Client
}

func NewReader(httpClient *http.Client) (*reader, error) {
	return &reader{
		httpClient: httpClient,
	}, nil
}

func (m *reader) GetMediaItems(email string, nextPageToken string) (MediaItems, error) {
	resp, err := m.httpClient.Get(
		"https://photoslibrary.googleapis.com/v1/mediaItems?pageSize=100&pageToken=" + nextPageToken,
	)
	if err != nil {
		return MediaItems{}, fmt.Errorf("media items search request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return MediaItems{}, TooManyRequestsError{}
	}

	if resp.StatusCode != http.StatusOK {
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return MediaItems{}, fmt.Errorf("read response body: %w", err)
		}

		return MediaItems{}, NotOkRequestError{fmt.Errorf(string(responseBody))}
	}

	var responseBody mediaItemsListResponseBody
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return MediaItems{}, fmt.Errorf("decode response body: %w", err)
	}

	return MediaItems{
		Items:         responseBody.MediaItems,
		NextPageToken: responseBody.NextPageToken,
	}, nil
}

func (m *reader) GetMediaItem(mediaItemId string) (MediaItem, error) {
	resp, err := m.httpClient.Get(
		"https://photoslibrary.googleapis.com/v1/mediaItems/" + mediaItemId,
	)
	if err != nil {
		return MediaItem{}, fmt.Errorf("media item get request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return MediaItem{}, TooManyRequestsError{}
	}

	if resp.StatusCode != http.StatusOK {
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return MediaItem{}, fmt.Errorf("read response body: %w", err)
		}

		return MediaItem{}, NotOkRequestError{fmt.Errorf(string(responseBody))}
	}

	var mediaItem MediaItem
	err = json.NewDecoder(resp.Body).Decode(&mediaItem)
	if err != nil {
		return MediaItem{}, fmt.Errorf("decode response body: %w", err)
	}

	return mediaItem, nil
}

// func (m *reader) getFilters(startDate, endDate time.Time) filters {
// 	return filters{
// 		DateFilter: dateFilter{
// 			Ranges: []ranges{
// 				{
// 					StartDate: date{
// 						Year:  startDate.Year(),
// 						Month: int(startDate.Month()),
// 						Day:   startDate.Day(),
// 					},
// 					EndDate: date{
// 						Year:  endDate.Year(),
// 						Month: int(endDate.Month()),
// 						Day:   endDate.Day(),
// 					},
// 				},
// 			},
// 		},
// 	}
// }
