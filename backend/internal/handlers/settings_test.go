package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"google-backup/internal/settings/settingsfakes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSettingsHandle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("get settings", func(t *testing.T) {
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		handler := NewSettingsHandler(fakeSettingsRepository)

		fakeSettingsRepository.FindReturns([]byte(`{"rootPath": "/root/path", "photosScannerJobDelay": 60000000000, "photosDownloaderJobDelay": 120000000000, "host": "http://localhost:8080", "photosBackupEnabled": true, "driveBackupEnabled": true}`), nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/settings", nil)

		handler.Handle(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"data":{"rootPath":"/root/path","photosScannerJobDelay":1,"photosDownloaderJobDelay":2,"host":"http://localhost:8080","photosBackupEnabled":true,"driveBackupEnabled":true}}`, w.Body.String())
	})

	t.Run("update settings", func(t *testing.T) {
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		handler := NewSettingsHandler(fakeSettingsRepository)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/settings", bytes.NewBuffer(
			[]byte(`{"rootPath": "/root/path", "photosScannerJobDelay": 1, "photosDownloaderJobDelay": 5, "host": "http://localhost:8080", "photosBackupEnabled": true, "driveBackupEnabled": true}`),
		))

		handler.Handle(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"data":{"rootPath":"/root/path","photosScannerJobDelay":1,"photosDownloaderJobDelay":5,"host":"http://localhost:8080","photosBackupEnabled":true,"driveBackupEnabled":true}}`, w.Body.String())

		settingsJson := fakeSettingsRepository.SaveArgsForCall(0)
		assert.Equal(t, `{"rootPath":"/root/path","photosScannerJobDelay":60000000000,"photosDownloaderJobDelay":300000000000,"host":"http://localhost:8080","photosBackupEnabled":true,"driveBackupEnabled":true}`, string(settingsJson))
	})

	t.Run("update settings validation", func(t *testing.T) {
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		handler := NewSettingsHandler(fakeSettingsRepository)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/settings", bytes.NewBuffer(
			[]byte(`{"data": []}`),
		))

		handler.Handle(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, `{"message":"Key: 'settingsUpdateRequest.RootPath' Error:Field validation for 'RootPath' failed on the 'required' tag\nKey: 'settingsUpdateRequest.PhotosScannerJobDelay' Error:Field validation for 'PhotosScannerJobDelay' failed on the 'required' tag\nKey: 'settingsUpdateRequest.PhotosDownloaderJobDelay' Error:Field validation for 'PhotosDownloaderJobDelay' failed on the 'required' tag\nKey: 'settingsUpdateRequest.Host' Error:Field validation for 'Host' failed on the 'required' tag\nKey: 'settingsUpdateRequest.PhotosBackupEnabled' Error:Field validation for 'PhotosBackupEnabled' failed on the 'required' tag\nKey: 'settingsUpdateRequest.DriveBackupEnabled' Error:Field validation for 'DriveBackupEnabled' failed on the 'required' tag"}`, w.Body.String())
		assert.Equal(t, 0, fakeSettingsRepository.SaveCallCount())
	})
}
