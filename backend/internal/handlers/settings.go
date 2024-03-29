package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"google-backup/internal/settings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type settingsApiHandler struct {
	settingsRepository settings.Repository
}

type settingsUpdateRequest struct {
	RootPath                 string `json:"rootPath" binding:"required,ascii"`
	PhotosScannerJobDelay    int64  `json:"photosScannerJobDelay" binding:"required,numeric"`
	PhotosDownloaderJobDelay int64  `json:"photosDownloaderJobDelay" binding:"required,numeric"`
	Host                     string `json:"host" binding:"required,ascii"`
	PhotosBackupEnabled      bool   `json:"photosBackupEnabled" binding:"required,boolean"`
	DriveBackupEnabled       bool   `json:"driveBackupEnabled" binding:"required,boolean"`
}

func NewSettingsHandler(settingsRepository settings.Repository) *settingsApiHandler {
	return &settingsApiHandler{settingsRepository: settingsRepository}
}

func (h *settingsApiHandler) Handle(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		h.handleGet(c)

		return
	case "POST":
		h.handlePost(c)

		return
	}

	c.JSON(http.StatusMethodNotAllowed, gin.H{})
}

func (h *settingsApiHandler) handleGet(c *gin.Context) {
	settingsJson, err := h.settingsRepository.Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

		return
	}

	var settingsData settings.SettingsData
	err = json.Unmarshal(settingsJson, &settingsData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		log.Error(fmt.Errorf("unmarshall settings json: %w", err))

		return
	}

	settingsData = h.convertDurationToMinutes(settingsData)

	c.JSON(http.StatusOK, gin.H{"data": settingsData})
}

func (h *settingsApiHandler) handlePost(c *gin.Context) {
	var request settingsUpdateRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		return
	}

	settingsData := settings.SettingsData{
		RootPath:                 request.RootPath,
		PhotosScannerJobDelay:    time.Duration(request.PhotosScannerJobDelay * int64(time.Minute)),
		PhotosDownloaderJobDelay: time.Duration(request.PhotosDownloaderJobDelay * int64(time.Minute)),
		Host:                     request.Host,
		PhotosBackupEnabled:      true,
		DriveBackupEnabled:       true,
	}

	settingsJson, err := json.Marshal(settingsData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		log.Error(fmt.Errorf("marshall settings data: %w", err))

		return
	}

	err = h.settingsRepository.Save(settingsJson)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		log.Error(fmt.Errorf("save settings: %w", err))

		return
	}

	settingsData = h.convertDurationToMinutes(settingsData)

	c.JSON(http.StatusOK, gin.H{"data": settingsData})
}

func (h *settingsApiHandler) convertDurationToMinutes(settingsData settings.SettingsData) settings.SettingsData {
	settingsData.PhotosScannerJobDelay = settingsData.PhotosScannerJobDelay / time.Minute
	settingsData.PhotosDownloaderJobDelay = settingsData.PhotosDownloaderJobDelay / time.Minute

	return settingsData
}
