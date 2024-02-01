package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"google-backup/internal/account/accountfakes"
	"google-backup/internal/google_client/google_clientfakes"
	"google-backup/internal/handlers"
	"google-backup/internal/settings/settingsfakes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type response struct {
	Data []struct {
		ID               string `json:"id"`
		Secret           string `json:"secret"`
		RedirectURL      string `json:"redirectUrl"`
		AssignedAccounts []struct {
			Email      string `json:"email"`
			GivenName  string `json:"givenName"`
			FamilyName string `json:"familyName"`
			Picture    string `json:"picture"`
		} `json:"assignedAccounts"`
	} `json:"data"`
}

func TestClientsHandle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("get list of clients without assigned accounts", func(t *testing.T) {
		fakeAccountRepository := new(accountfakes.FakeRepository)
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		handler := handlers.NewClientsApiHandler(fakeAccountRepository, fakeGoogleClientRepository, fakeSettingsRepository)

		fakeGoogleClientRepository.FindAllReturns(map[string][]byte{
			"id1": []byte(`{"id":"id1","secret":"secret1","redirectUrl":"http://localhost:8080/redirect_url/id1"}`),
			"id2": []byte(`{"id":"id2","secret":"secret2","redirectUrl":"http://localhost:8080/redirect_url/id2"}`),
		}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/clients", nil)

		handler.Handle(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"data":[{"id":"id1","secret":"secret1","redirectUrl":"http://localhost:8080/redirect_url/id1","assignedAccounts":[]},{"id":"id2","secret":"secret2","redirectUrl":"http://localhost:8080/redirect_url/id2","assignedAccounts":[]}]}`, w.Body.String())
	})

	t.Run("get list of clients with assigned accounts", func(t *testing.T) {
		fakeAccountRepository := new(accountfakes.FakeRepository)
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		handler := handlers.NewClientsApiHandler(fakeAccountRepository, fakeGoogleClientRepository, fakeSettingsRepository)

		fakeGoogleClientRepository.FindAllReturns(map[string][]byte{
			"id1": []byte(`{"id":"id1","secret":"secret1","redirectUrl":"http://localhost:8080/redirect_url/id1"}`),
			"id2": []byte(`{"id":"id2","secret":"secret2","redirectUrl":"http://localhost:8080/redirect_url/id2"}`),
		}, nil)

		fakeGoogleClientRepository.FindAssignedAccountsReturns([]byte(`["email1@test.com","email2@test.com"]`), nil)

		fakeAccountRepository.FindAccountReturns([]byte(`{"email":"email","picture":"picture","givenName":"Bob","familyName":"Alice"}`), nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/clients", nil)

		handler.Handle(c)

		var actualData response
		_ = json.Unmarshal(w.Body.Bytes(), &actualData)

		var expectedData response
		_ = json.Unmarshal([]byte(`{"data":[{"id":"id1","secret":"secret1","redirectUrl":"http://localhost:8080/redirect_url/id1","assignedAccounts":[{"email":"email","givenName":"Bob","familyName":"Alice","picture":"picture"},{"email":"email","givenName":"Bob","familyName":"Alice","picture":"picture"}]},{"id":"id2","secret":"secret2","redirectUrl":"http://localhost:8080/redirect_url/id2","assignedAccounts":[{"email":"email","givenName":"Bob","familyName":"Alice","picture":"picture"},{"email":"email","givenName":"Bob","familyName":"Alice","picture":"picture"}]}]}`), &expectedData)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.ElementsMatch(t, expectedData.Data, actualData.Data)
	})

	t.Run("get one client", func(t *testing.T) {
		fakeAccountRepository := new(accountfakes.FakeRepository)
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		handler := handlers.NewClientsApiHandler(fakeAccountRepository, fakeGoogleClientRepository, fakeSettingsRepository)

		fakeGoogleClientRepository.FindReturns(
			[]byte(`{"id":"id1","secret":"secret1","redirectUrl":"http://localhost:8080/redirect_url/id1"}`),
			nil,
		)

		fakeGoogleClientRepository.FindAssignedAccountsReturns([]byte(`["email1@test.com","email2@test.com"]`), nil)

		fakeAccountRepository.FindAccountReturns([]byte(`{"email":"email","picture":"picture","givenName":"Bob","familyName":"Alice"}`), nil)

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
		assert.Equal(t, `{"data":{"id":"id1","secret":"secret1","redirectUrl":"http://localhost:8080/redirect_url/id1","assignedAccounts":[{"email":"email","givenName":"Bob","familyName":"Alice","picture":"picture"},{"email":"email","givenName":"Bob","familyName":"Alice","picture":"picture"}]}}`, w.Body.String())
	})

	t.Run("list of clients not found", func(t *testing.T) {
		fakeAccountRepository := new(accountfakes.FakeRepository)
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		handler := handlers.NewClientsApiHandler(fakeAccountRepository, fakeGoogleClientRepository, fakeSettingsRepository)

		fakeGoogleClientRepository.FindAllReturns(nil, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/clients", nil)

		handler.Handle(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"data":[]}`, w.Body.String())
	})

	t.Run("one client not found", func(t *testing.T) {
		fakeAccountRepository := new(accountfakes.FakeRepository)
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		handler := handlers.NewClientsApiHandler(fakeAccountRepository, fakeGoogleClientRepository, fakeSettingsRepository)

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
		fakeAccountRepository := new(accountfakes.FakeRepository)
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		handler := handlers.NewClientsApiHandler(fakeAccountRepository, fakeGoogleClientRepository, fakeSettingsRepository)

		fakeGoogleClientRepository.FindAllReturns(nil, errors.New("error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/clients", nil)

		handler.Handle(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, `{"message":"error"}`, w.Body.String())
	})

	t.Run("one client not found", func(t *testing.T) {
		fakeAccountRepository := new(accountfakes.FakeRepository)
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		handler := handlers.NewClientsApiHandler(fakeAccountRepository, fakeGoogleClientRepository, fakeSettingsRepository)

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
		assert.Equal(t, `{"message":"error"}`, w.Body.String())
	})

	t.Run("create a client", func(t *testing.T) {
		fakeAccountRepository := new(accountfakes.FakeRepository)
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		handler := handlers.NewClientsApiHandler(fakeAccountRepository, fakeGoogleClientRepository, fakeSettingsRepository)

		fakeSettingsRepository.FindReturns(
			[]byte(`{"host":"http://domain"}`),
			nil,
		)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/clients", bytes.NewBuffer(
			[]byte(`{"id":"id1","secret":"secret1"}`),
		))

		handler.Handle(c)

		clientId, clientData := fakeGoogleClientRepository.SaveArgsForCall(0)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, `{"data":{"id":"id1","secret":"secret1","redirectUrl":"http://domain/auth/google/callback/id1"}}`, w.Body.String())
		assert.Equal(t, 1, fakeGoogleClientRepository.SaveCallCount())
		assert.Equal(t, "id1", clientId)
		assert.Equal(t, `{"id":"id1","secret":"secret1","redirectUrl":"http://domain/auth/google/callback/id1"}`, string(clientData))
	})

	t.Run("create a client validation", func(t *testing.T) {
		fakeAccountRepository := new(accountfakes.FakeRepository)
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		handler := handlers.NewClientsApiHandler(fakeAccountRepository, fakeGoogleClientRepository, fakeSettingsRepository)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/clients", bytes.NewBuffer(
			[]byte("{}"),
		))

		handler.Handle(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, `{"message":"Key: 'updateClientRequest.ID' Error:Field validation for 'ID' failed on the 'required' tag\nKey: 'updateClientRequest.Secret' Error:Field validation for 'Secret' failed on the 'required' tag"}`, w.Body.String())
		assert.Equal(t, 0, fakeGoogleClientRepository.SaveCallCount())
	})

	t.Run("delete a client", func(t *testing.T) {
		fakeAccountRepository := new(accountfakes.FakeRepository)
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		handler := handlers.NewClientsApiHandler(fakeAccountRepository, fakeGoogleClientRepository, fakeSettingsRepository)

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
		fakeAccountRepository := new(accountfakes.FakeRepository)
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		handler := handlers.NewClientsApiHandler(fakeAccountRepository, fakeGoogleClientRepository, fakeSettingsRepository)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPut, "/clients", nil)

		handler.Handle(c)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}
