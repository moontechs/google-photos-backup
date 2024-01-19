package media

type mediaItemsListResponseBody struct {
	MediaItems    []MediaItem `json:"mediaItems"`
	NextPageToken string      `json:"nextPageToken"`
}

type MediaItem struct {
	ID              string          `json:"id"`
	Description     string          `json:"description"`
	ProductUrl      string          `json:"productUrl"`
	BaseUrl         string          `json:"baseUrl"`
	MimeType        string          `json:"mimeType"`
	Filename        string          `json:"filename"`
	MediaMetadata   MediaMetadata   `json:"mediaMetadata"`
	ContributorInfo ContributorInfo `json:"contributorInfo"`
}

type MediaMetadata struct {
	CreationTime string `json:"creationTime"`
	Width        string `json:"width"`
	Height       string `json:"height"`
	Photo        Photo  `json:"photo"`
	Video        Video  `json:"video"`
}

type Photo struct {
	CameraMake      string  `json:"cameraMake"`
	CameraModel     string  `json:"cameraModel"`
	FocalLength     float64 `json:"focalLength"`
	ApertureFNumber float64 `json:"apertureFNumber"`
	IsoEquivalent   int     `json:"isoEquivalent"`
	ExposureTime    string  `json:"exposureTime"`
}

type Video struct {
	CameraMake  string  `json:"cameraMake"`
	CameraModel string  `json:"cameraModel"`
	Fps         float64 `json:"fps"`
	Status      string  `json:"status"`
}

type ContributorInfo struct {
	ProfilePictureBaseUrl string `json:"profilePictureBaseUrl"`
	DisplayName           string `json:"displayName"`
}

type MediaItems struct {
	Items         []MediaItem
	NextPageToken string
}

type requestBody struct {
	Filters   filters `json:"filters"`
	PageToken string  `json:"pageToken"`
	OrderBy   string  `json:"orderBy"`
}

type filters struct {
	DateFilter dateFilter `json:"dateFilter"`
}

type dateFilter struct {
	Ranges []ranges `json:"ranges"`
}

type ranges struct {
	StartDate date `json:"startDate"`
	EndDate   date `json:"endDate"`
}

type date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}
