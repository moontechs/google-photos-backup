package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"google-backup/internal/account/accountfakes"
	"google-backup/internal/auth"
	"google-backup/internal/auth/authfakes"
	"google-backup/internal/google_client/google_clientfakes"
	"google-backup/internal/settings/settingsfakes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGoogleCallbackHandle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("receive callback", func(t *testing.T) {
		fakeAccountRepository := new(accountfakes.FakeRepository)
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		fakeGoogleAuth := new(authfakes.FakeAuth)
		handler := NewGoogleCallbackHandler(fakeAccountRepository, fakeGoogleClientRepository, fakeSettingsRepository, fakeGoogleAuth)

		fakeGoogleAuth.GetUserInfoReturns(auth.UserInfo{
			Picture: "picture",
			Email:   "email",
		}, nil)

		fakeSettingsRepository.FindReturns([]byte(`{"host": "http://localhost:8080"}`), nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{
			{
				Key:   "clientId",
				Value: "id1",
			},
		}
		c.Request, _ = http.NewRequest(http.MethodGet, "http://localhost:8080/auth/google/callback/:clientId", nil)

		handler.Handle(c)

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Equal(t, "http://localhost:8080", w.Header().Get("Location"))
		assert.Equal(t, 1, fakeAccountRepository.SaveTokenCallCount())
		assert.Equal(t, 1, fakeGoogleClientRepository.SaveAssignedAccountsCallCount())

		clientId, assignedAccounts := fakeGoogleClientRepository.SaveAssignedAccountsArgsForCall(0)
		assert.Equal(t, "id1", clientId)
		assert.Equal(t, []byte(`["email"]`), assignedAccounts)
	})

	t.Run("receive callback account assigned to another client", func(t *testing.T) {
		fakeAccountRepository := new(accountfakes.FakeRepository)
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeSettingsRepository := new(settingsfakes.FakeRepository)
		fakeGoogleAuth := new(authfakes.FakeAuth)
		handler := NewGoogleCallbackHandler(fakeAccountRepository, fakeGoogleClientRepository, fakeSettingsRepository, fakeGoogleAuth)

		fakeGoogleAuth.GetUserInfoReturns(auth.UserInfo{
			Picture: "picture",
			Email:   "email",
		}, nil)

		fakeSettingsRepository.FindReturns([]byte(`{"host": "http://localhost:8080"}`), nil)

		fakeGoogleClientRepository.FindAllAssignedAccountsReturns(map[string][]byte{
			"id2": []byte(`["email", "email2"]`),
		}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{
			{
				Key:   "clientId",
				Value: "id1",
			},
		}
		c.Request, _ = http.NewRequest(http.MethodGet, "http://localhost:8080/auth/google/callback/:clientId", nil)

		handler.Handle(c)

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Equal(t, "http://localhost:8080", w.Header().Get("Location"))
		assert.Equal(t, 1, fakeAccountRepository.SaveTokenCallCount())
		assert.Equal(t, 2, fakeGoogleClientRepository.SaveAssignedAccountsCallCount())

		// first call to unassign account from a previous client
		clientId, assignedAccounts := fakeGoogleClientRepository.SaveAssignedAccountsArgsForCall(0)
		assert.Equal(t, "id2", clientId)
		assert.Equal(t, []byte(`["email2"]`), assignedAccounts)

		// second call to assign account to a new client
		clientId, assignedAccounts = fakeGoogleClientRepository.SaveAssignedAccountsArgsForCall(1)
		assert.Equal(t, "id1", clientId)
		assert.Equal(t, []byte(`["email"]`), assignedAccounts)
	})
}
