package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"google-backup/internal/google_client/google_clientfakes"
	"google-backup/internal/settings/settingsfakes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestClientsHandle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("get list of clients", func(t *testing.T) {
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		handler := NewClientsApiHandler(fakeGoogleClientRepository, fakeSettingsRepository)

		fakeGoogleClientRepository.FindAllReturns(map[string][]byte{
			"id1": []byte("{\"id\":\"id1\",\"secret\":\"secret1\",\"redirect_url\":\"http://localhost:8080/redirect_url/id1\"}"),
			"id2": []byte("{\"id\":\"id2\",\"secret\":\"secret2\",\"redirect_url\":\"http://localhost:8080/redirect_url/id2\"}"),
		}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/clients", nil)

		handler.Handle(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "{\"data\":[{\"id\":\"id1\",\"secret\":\"secret1\",\"redirect_url\":\"http://localhost:8080/redirect_url/id1\"},{\"id\":\"id2\",\"secret\":\"secret2\",\"redirect_url\":\"http://localhost:8080/redirect_url/id2\"}]}", w.Body.String())
	})

	t.Run("get one client", func(t *testing.T) {
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		handler := NewClientsApiHandler(fakeGoogleClientRepository, fakeSettingsRepository)

		fakeGoogleClientRepository.FindReturns(
			[]byte("{\"id\":\"id1\",\"secret\":\"secret1\",\"redirect_url\":\"http://localhost:8080/redirect_url/id1\"}"),
			nil,
		)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{
			{
				Key:   "clientId",
				Value: "id1",
			},
		}
		c.Request, _ = http.NewRequest(http.MethodGet, "/clients/id1", nil)

		handler.Handle(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "{\"data\":{\"id\":\"id1\",\"secret\":\"secret1\",\"redirect_url\":\"http://localhost:8080/redirect_url/id1\"}}", w.Body.String())
	})

	t.Run("list of clients not found", func(t *testing.T) {
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		handler := NewClientsApiHandler(fakeGoogleClientRepository, fakeSettingsRepository)

		fakeGoogleClientRepository.FindAllReturns(nil, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/clients", nil)

		handler.Handle(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "{\"data\":[]}", w.Body.String())
	})

	t.Run("one client not found", func(t *testing.T) {
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		handler := NewClientsApiHandler(fakeGoogleClientRepository, fakeSettingsRepository)

		fakeGoogleClientRepository.FindReturns(nil, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{
			{
				Key:   "clientId",
				Value: "id1",
			},
		}
		c.Request, _ = http.NewRequest(http.MethodGet, "/clients/id1", nil)

		handler.Handle(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("list of clients error", func(t *testing.T) {
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		handler := NewClientsApiHandler(fakeGoogleClientRepository, fakeSettingsRepository)

		fakeGoogleClientRepository.FindAllReturns(nil, errors.New("error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/clients", nil)

		handler.Handle(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "{\"message\":\"error\"}", w.Body.String())
	})

	t.Run("one client not found", func(t *testing.T) {
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		handler := NewClientsApiHandler(fakeGoogleClientRepository, fakeSettingsRepository)

		fakeGoogleClientRepository.FindReturns(nil, errors.New("error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{
			{
				Key:   "clientId",
				Value: "id1",
			},
		}
		c.Request, _ = http.NewRequest(http.MethodGet, "/clients/id1", nil)

		handler.Handle(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "{\"message\":\"error\"}", w.Body.String())
	})

	t.Run("create a client", func(t *testing.T) {
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		handler := NewClientsApiHandler(fakeGoogleClientRepository, fakeSettingsRepository)

		fakeSettingsRepository.FindReturns(
			[]byte("{\"domain\":\"http://domain\"}"),
			nil,
		)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/clients", bytes.NewBuffer(
			[]byte("{\"id\":\"id1\",\"secret\":\"secret1\"}"),
		))

		handler.Handle(c)

		clientId, clientData := fakeGoogleClientRepository.SaveArgsForCall(0)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, "{\"data\":{\"id\":\"id1\",\"secret\":\"secret1\",\"redirect_url\":\"http://domain/auth/google/callback/id1\"}}", w.Body.String())
		assert.Equal(t, 1, fakeGoogleClientRepository.SaveCallCount())
		assert.Equal(t, "id1", clientId)
		assert.Equal(t, "{\"id\":\"id1\",\"secret\":\"secret1\",\"redirect_url\":\"http://domain/auth/google/callback/id1\"}", string(clientData))
	})

	t.Run("create a client validation", func(t *testing.T) {
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		handler := NewClientsApiHandler(fakeGoogleClientRepository, fakeSettingsRepository)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/clients", bytes.NewBuffer(
			[]byte("{}"),
		))

		handler.Handle(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "{\"message\":\"Key: 'updateClientRequest.ID' Error:Field validation for 'ID' failed on the 'required' tag\\nKey: 'updateClientRequest.Secret' Error:Field validation for 'Secret' failed on the 'required' tag\"}", w.Body.String())
		assert.Equal(t, 0, fakeGoogleClientRepository.SaveCallCount())
	})

	t.Run("delete a client", func(t *testing.T) {
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		handler := NewClientsApiHandler(fakeGoogleClientRepository, fakeSettingsRepository)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{
			{
				Key:   "clientId",
				Value: "id1",
			},
		}
		c.Request, _ = http.NewRequest(http.MethodDelete, "/client/id1", nil)

		handler.Handle(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, 1, fakeGoogleClientRepository.DeleteCallCount())
		assert.Equal(t, "id1", fakeGoogleClientRepository.DeleteArgsForCall(0))
	})

	t.Run("unsupported method", func(t *testing.T) {
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		handler := NewClientsApiHandler(fakeGoogleClientRepository, fakeSettingsRepository)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPut, "/clients", nil)

		handler.Handle(c)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}
