package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"google-backup/internal/auth"
	"google-backup/internal/auth/authfakes"
	"google-backup/internal/google_client/google_clientfakes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRedirectUrlHandle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("get redirect url", func(t *testing.T) {
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeAuthRepository := new(authfakes.FakeRepository)
		googleAuth := auth.NewGoogleAuth(
			fakeAuthRepository,
			fakeGoogleClientRepository,
		)
		handler := NewGoogleRedirectUrlHandler(googleAuth)

		fakeGoogleClientRepository.FindReturns([]byte(`{"id":"id1","secret":"secret1","redirectUrl":"http://localhost:8080/redirect_url/id1"}`), nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{
			{
				Key:   "clientId",
				Value: "id1",
			},
		}
		c.Request, _ = http.NewRequest(http.MethodGet, "clients/:clientId/redirect-url", nil)

		handler.Handle(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"data": "https://accounts.google.com/o/oauth2/auth?access_type=offline&client_id=id1&prompt=consent&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fredirect_url%2Fid1&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fphotoslibrary.readonly+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fdrive.readonly+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.profile+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email&state=state"}`, w.Body.String())
	})

	t.Run("get redirect url client not found", func(t *testing.T) {
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeAuthRepository := new(authfakes.FakeRepository)
		googleAuth := auth.NewGoogleAuth(
			fakeAuthRepository,
			fakeGoogleClientRepository,
		)
		handler := NewGoogleRedirectUrlHandler(googleAuth)

		fakeGoogleClientRepository.FindReturns(nil, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{
			{
				Key:   "clientId",
				Value: "id1",
			},
		}
		c.Request, _ = http.NewRequest(http.MethodGet, "clients/:clientId/redirect-url", nil)

		handler.Handle(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("get redirect url post method not allowed", func(t *testing.T) {
		fakeGoogleClientRepository := new(google_clientfakes.FakeRepository)
		fakeAuthRepository := new(authfakes.FakeRepository)
		googleAuth := auth.NewGoogleAuth(
			fakeAuthRepository,
			fakeGoogleClientRepository,
		)
		handler := NewGoogleRedirectUrlHandler(googleAuth)

		fakeGoogleClientRepository.FindReturns(nil, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{
			{
				Key:   "clientId",
				Value: "id1",
			},
		}
		c.Request, _ = http.NewRequest(http.MethodPost, "clients/:clientId/redirect-url", nil)

		handler.Handle(c)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}
